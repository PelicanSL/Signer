package db

import (
	"Signer/internal/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// Database connection
var DB *sql.DB

func InitDB() {
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connection successful!")
}

// EmailExists checks if an email is already registered in the database.
func EmailExists(email string) bool {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

// CreateUser creates a new user in the database and returns the created user.
func CreateUser(email, password string) (*models.User, error) {
	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Prepare the SQL statement to insert the user
	stmt, err := DB.Prepare("INSERT INTO users(email, password) VALUES($1, $2) RETURNING email")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement and get the generated email
	var userEmail string
	err = stmt.QueryRow(email, string(hashedPassword)).Scan(&userEmail)
	if err != nil {
		return nil, err
	}

	// Create the User object
	user := &models.User{
		Email: userEmail,
		// Do not include the password or JWT here
	}

	return user, nil
}

// StoreAnswer stores an answer in the database.
func StoreAnswer(tx *sql.Tx, email string, questionID int64, answerText string) error {
	query := `
        INSERT INTO answers (question_id, user_email, text)
        VALUES ($1, $2, $3)
        ON CONFLICT (question_id, user_email)
        DO UPDATE SET text = EXCLUDED.text;
    `

	_, err := tx.Exec(query, questionID, email, answerText)
	return err
}

// StoreSignature stores a signature in the database.
func StoreSignature(tx *sql.Tx, email, signature string) error {
	query := `
        INSERT INTO signatures (user_email, sign)
        VALUES ($1, $2)
        ON CONFLICT (user_email)
        DO UPDATE SET sign = EXCLUDED.sign;
    `

	_, err := tx.Exec(query, email, signature)
	return err
}

// VerifySignature checks if a given signature belongs to a user and retrieves associated answers.
func VerifySignature(userEmail, signature string) (bool, []models.Answer, time.Time, error) {
	var (
		timestamp time.Time
		answers   []models.Answer
	)

	// First, check if the signature corresponds to the user
	err := DB.QueryRow("SELECT timestamp FROM signatures WHERE user_email = $1 AND sign = $2", userEmail, signature).Scan(&timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			// Signature not found
			return false, nil, time.Time{}, nil
		}
		// Database query error
		log.Printf("Error verifying the signature: %v", err)
		return false, nil, time.Time{}, err
	}

	// If the signature exists, retrieve associated answers
	rows, err := DB.Query("SELECT id, question_id, text FROM answers WHERE user_email = $1", userEmail)
	if err != nil {
		log.Printf("Error getting responses: %v", err)
		return false, nil, time.Time{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var ans models.Answer
		if err := rows.Scan(&ans.ID, &ans.QuestionID, &ans.Text); err != nil {
			log.Printf("Error reading response row: %v", err)
			continue
		}
		answers = append(answers, ans)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over responses: %v", err)
		return false, nil, time.Time{}, err
	}

	return true, answers, timestamp, nil
}

// DeleteUserAnswers deletes all answers associated with a user.
func DeleteUserAnswers(tx *sql.Tx, userEmail string) error {
	_, err := tx.Exec("DELETE FROM answers WHERE user_email = $1", userEmail)
	return err
}
