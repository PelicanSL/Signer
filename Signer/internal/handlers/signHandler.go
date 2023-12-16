package handlers

import (
	"Signer/internal/db"
	"Signer/internal/models"
	"Signer/internal/util"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// SignAnswers handles the signing of answers by a user.
func SignAnswers(c *gin.Context) {
	var req models.SignAnswersRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the email from the JWT
	token := c.GetHeader("Authorization")
	email, err := util.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting the transaction"})
		return
	}

	// Make sure to rollback in case of an error
	defer func() {
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing the responses"})
		} else {
			err = tx.Commit()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finalizing the transaction"})
			}
		}
	}()

	// Delete all previous user answers
	if err = db.DeleteUserAnswers(tx, email); err != nil {
		return
	}

	// Use WaitGroup to handle concurrency
	var wg sync.WaitGroup
	for _, answer := range req.Answers {
		wg.Add(1)
		go func(answer models.AnswerSubmission) {
			defer wg.Done()
			if err = db.StoreAnswer(tx, email, answer.QuestionID, answer.Text); err != nil {
				return
			}
		}(answer)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Generate a signature for the answers
	signature := util.GenerateSignature(req.Answers)
	if err = db.StoreSignature(tx, email, signature); err != nil {
		return
	}

	// If it reaches here, everything has gone well
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"signature": signature})
	}
}
