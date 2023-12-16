package handlers

import (
	"Signer/internal/db"
	"Signer/internal/models"
	"Signer/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUser handles user registration.
func RegisterUser(c *gin.Context) {
	var req models.UserRegistrationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error processing the data"})
		return
	}

	// Check if the provided email is valid
	if !util.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// Check if the provided password is valid
	if !util.IsValidPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters, include an uppercase letter, a symbol, and a number, and not exceed 100 characters"})
		return
	}

	// Check if the email is already registered
	if db.EmailExists(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already registered"})
		return
	}

	// Create a new user in the database
	user, err := db.CreateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the user"})
		return
	}

	// Generate a JWT token for the registered user
	token, err := util.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating the token"})
		return
	}

	// Respond with a success message and the generated token
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": token})
}
