package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/user/server-moni/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

// Register creates a new user and returns their API Key
func Register(username, password string) (string, error) {
	// Check if user exists
	var exists int
	err := db.DB.QueryRow("SELECT 1 FROM users WHERE username = ?", username).Scan(&exists)
	if err == nil {
		return "", errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Generate API Key
	apiKey := uuid.New().String()

	// Insert User
	_, err = db.DB.Exec("INSERT INTO users (username, password_hash, api_key) VALUES (?, ?, ?)", username, string(hashedPassword), apiKey)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

// Login verifies credentials and returns the API Key
func Login(username, password string) (string, error) {
	var hashedPassword string
	var apiKey string

	err := db.DB.QueryRow("SELECT password_hash, api_key FROM users WHERE username = ?", username).Scan(&hashedPassword, &apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return apiKey, nil
}

// GetUserByAPIKey returns the user associated with the given API Key
func GetUserByAPIKey(apiKey string) (*User, error) {
	var user User
	err := db.DB.QueryRow("SELECT id, username, api_key, created_at FROM users WHERE api_key = ?", apiKey).Scan(&user.ID, &user.Username, &user.APIKey, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
