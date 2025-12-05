package auth

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/server-moni/internal/db"
	"sync"
	"time"
)

var (
	AgentTokens = make(map[string]AgentTokenInfo)
	tokensMutex sync.RWMutex
)

type AgentTokenInfo struct {
	CreatedAt time.Time
	Name      string
}

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
	}

	// Write key to file for user convenience
	_ = os.Mkdir("data", 0755)
	err = os.WriteFile("data/api_key.txt", []byte(apiKey), 0644)
	if err != nil {
		log.Printf("Failed to write api_key.txt: %v", err)
	}
	
	// Load Agent Tokens (TODO: Persist in DB, for now in-memory is fine as per request scope, or we can use DB)
	// For simplicity in this iteration, we start empty. 
	// Ideally we should load from DB.
}

func AddAgentToken(token string, name string) {
	tokensMutex.Lock()
	defer tokensMutex.Unlock()
	AgentTokens[token] = AgentTokenInfo{
		CreatedAt: time.Now(),
		Name:      name,
	}
}

func ValidateAgentToken(token string) bool {
	tokensMutex.RLock()
	defer tokensMutex.RUnlock()
	_, ok := AgentTokens[token]
	return ok
}

// Middleware for Dashboard/Admin access (Master Key)
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

// Middleware for Agent Ingestion (Agent Tokens)
func AgentAuthMiddleware() gin.HandlerFunc {
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
		if !ValidateAgentToken(token) {
			// Also allow Master Key for testing/backward compatibility
			storedKey, _ := db.GetConfig("api_key")
			if token != storedKey {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Agent Token"})
				return
			}
		}

		c.Next()
	}
}
