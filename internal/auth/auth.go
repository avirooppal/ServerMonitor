package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/user/server-moni/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	AgentTokens = make(map[string]AgentTokenInfo)
	tokensMutex sync.RWMutex
	SaasMode    bool
	JWTSecret   []byte
)

type AgentTokenInfo struct {
	CreatedAt time.Time
	Name      string
}

func InitAuth() {
	mode := os.Getenv("SAAS_MODE")
	SaasMode = mode == "true"

	if SaasMode {
		log.Println("Running in SAAS MODE")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Fatal("JWT_SECRET is required in SAAS_MODE")
		}
		JWTSecret = []byte(secret)
	} else {
		log.Println("Running in SELF-HOSTED MODE")
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
	}
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

// --- SaaS Auth Functions ---

func Register(email, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return db.CreateUser(email, string(hash), "local", "")
}

func Login(email, password string) (string, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return GenerateJWT(user.ID)
}

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})
	return token.SignedString(JWTSecret)
}

func GitHubLogin(code string) (string, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	// Exchange code for token
	reqBody := fmt.Sprintf(`{"client_id": "%s", "client_secret": "%s", "code": "%s"}`, clientID, clientSecret, code)
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	// Get User Info
	req, _ = http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var userResp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return "", err
	}

	email := userResp.Email
	if email == "" {
		// Fallback if email is private (simplified)
		email = fmt.Sprintf("%s@github.com", userResp.Login)
	}

	// Find or Create User
	user, err := db.GetUserByEmail(email)
	if err != nil {
		// Create
		id, err := db.CreateUser(email, "", "github", fmt.Sprintf("%d", userResp.ID))
		if err != nil {
			return "", err
		}
		return GenerateJWT(int(id))
	}

	return GenerateJWT(user.ID)
}

// --- Middleware ---

// Middleware for Dashboard/Admin access
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
		tokenString := parts[1]

		if SaasMode {
			// JWT Validation
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return JWTSecret, nil
			})

			if err != nil || !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID := int(claims["user_id"].(float64))
				c.Set("user_id", userID)
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Claims"})
				return
			}

		} else {
			// Master Key Validation (Self-Hosted)
			storedKey, err := db.GetConfig("api_key")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			if tokenString != storedKey {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
				return
			}
			// Set user_id to 0 for self-hosted admin
			c.Set("user_id", 0)
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
			// Also allow Master Key for testing/backward compatibility in Self-Hosted
			if !SaasMode {
				storedKey, _ := db.GetConfig("api_key")
				if token != storedKey {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Agent Token"})
					return
				}
			} else {
				// In SaaS mode, strict agent token check
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Agent Token"})
				return
			}
		}

		c.Next()
	}
}
