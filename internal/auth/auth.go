package auth

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	AgentTokens = make(map[string]AgentTokenInfo)
	tokensMutex sync.RWMutex
	StaticKey   string
)

type AgentTokenInfo struct {
	CreatedAt time.Time
	Name      string
}

func InitAuth() {
	// No initialization needed for open dashboard
}

func SetStaticKey(key string) {
	StaticKey = key
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
	if StaticKey != "" {
		return token == StaticKey
	}
	tokensMutex.RLock()
	defer tokensMutex.RUnlock()
	_, ok := AgentTokens[token]
	return ok
}

// AuthMiddleware is now a no-op for the dashboard (Public Access)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Public Dashboard: No check required.
		c.Next()
	}
}

// Middleware for Agent Ingestion (Agent Tokens)
func AgentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization format"})
			return
		}

		token := parts[1]
		if !ValidateAgentToken(token) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Agent Token"})
			return
		}

		c.Next()
	}
}
