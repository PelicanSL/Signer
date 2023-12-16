package handlers

import (
	"Signer/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifySignature handles the verification of a user's signature.
func VerifySignature(c *gin.Context) {
	// Structure to receive request data
	type verifyRequest struct {
		UserEmail string `json:"user_email"`
		Signature string `json:"signature"`
	}

	var req verifyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the signature
	ok, answers, timestamp, err := db.VerifySignature(req.UserEmail, req.Signature)
	if err != nil {
		// Handle the error appropriately
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying the signature"})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature or does not correspond to the user"})
		return
	}

	// Respond if the signature is valid
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"answers":   answers,
		"timestamp": timestamp,
	})
}
