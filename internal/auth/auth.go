package auth

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/server-moni/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func Register(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.CreateUser(email, string(hashedPassword))
}

func Login(email, password string) (string, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return "", err // User not found
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err // Invalid password
	}

	// Generate Session Token
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // 1 day session

	if err := db.CreateSession(token, user.ID, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}

func Logout(token string) error {
	return db.DeleteSession(token)
}

// Middleware to protect routes and inject UserID
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		token := parts[1]
		session, err := db.GetSession(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if time.Now().After(session.ExpiresAt) {
			db.DeleteSession(token)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		// Set UserID in context
		c.Set("userID", session.UserID)
		c.Next()
	}
}

// InitAuth is no longer needed for Master Key, but we keep it empty or remove it.
// We'll remove it to clean up.
