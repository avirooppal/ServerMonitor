package auth

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/server-moni/internal/db"
)

func InitAuth() {
	var apiKey string
	val, err := db.GetConfig("api_key")
	if err != nil {
		// Key doesn't exist, generate one
		newKey := uuid.New().String()
		err = db.SetConfig("api_key", newKey)
		if err != nil {
			log.Fatalf("Failed to save API key: %v", err)
		}
		apiKey = newKey
		log.Printf("==================================================")
		log.Printf("NEW API KEY GENERATED: %s", newKey)
		log.Printf("COPY THIS KEY TO SETUP THE FRONTEND")
		log.Printf("==================================================")
	} else {
		apiKey = val
		// log.Println("API Key loaded from database.")
	}

	// Write key to file for user convenience
	// We assume "data" directory exists or we create it
	_ = os.Mkdir("data", 0755)
	err = os.WriteFile("data/api_key.txt", []byte(apiKey), 0644)
	if err != nil {
		log.Printf("Failed to write api_key.txt: %v", err)
	}
}

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

		apiKey := parts[1]
		storedKey, err := db.GetConfig("api_key")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if apiKey != storedKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			return
		}

		c.Next()
	}
}
