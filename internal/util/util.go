package util

import (
	"Signer/internal/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateJWT generates a JWT for a user.
func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) // Token valid for 1 hour

	// Create a new claim with user information and expiration time
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   fmt.Sprintf("%s", user.Email),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates and extracts the subject (email) from a JWT.
func ValidateJWT(tokenString string) (string, error) {

	// Extract the token from the Bearer scheme
	splitToken := strings.Split(tokenString, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token format")
	}
	tokenString = splitToken[1]

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil
}

// GenerateSignature creates a unique signature for a set of answers.
func GenerateSignature(answers []models.AnswerSubmission) string {
	// Sort the answers by question ID to ensure consistent order
	sort.Slice(answers, func(i, j int) bool {
		return answers[i].QuestionID < answers[j].QuestionID
	})

	// Concatenate answers with their question ID to form a string
	var concatenatedAnswers []string
	for _, answer := range answers {
		concatenatedAnswers = append(concatenatedAnswers, fmt.Sprintf("%d:%s", answer.QuestionID, answer.Text))
	}
	answerString := strings.Join(concatenatedAnswers, "|")

	// Add a timestamp for uniqueness
	timestamp := time.Now().Unix()
	signatureBase := fmt.Sprintf("%s:%d", answerString, timestamp)

	// Create an SHA-256 hash of the string
	hasher := sha256.New()
	hasher.Write([]byte(signatureBase))
	return hex.EncodeToString(hasher.Sum(nil))
}

// IsValidEmail checks if the provided string is a valid email address.
func IsValidEmail(email string) bool {
	if len(email) > 100 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if the provided string is a valid password.
func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 100 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
