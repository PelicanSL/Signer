package models

// User represents a user of the system.
type User struct {
	Email    string // Unique identifier of the user (Email)
	Password string // User's password
}

// UserRegistrationRequest represents a user registration request with JSON tags.
type UserRegistrationRequest struct {
	Email    string `json:"email"`    // Email of the user
	Password string `json:"password"` // User's password
}
