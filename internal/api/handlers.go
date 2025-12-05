package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// Public Endpoints
	api.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	api.POST("/register", Register)
	api.POST("/login", Login)
	
	// Ingestion Endpoint (Public, but requires Valid API Key in Header)
	// We use a custom middleware or just check in handler.
	// Since we want to support "One Command" where agent sends API Key,
	// we should check if the API Key belongs to a User or is a System Key.
	// For now, let's assume the Agent sends the User's API Key.
	api.POST("/ingest", IngestMetrics)

	// Authenticated User Endpoints
	authenticated := api.Group("/")
	authenticated.Use(auth.AuthMiddleware())
	{
		authenticated.POST("/verify-key", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "valid"})
		})
		authenticated.GET("/metrics", GetMetrics)
		authenticated.GET("/systems", GetSystems)
		authenticated.POST("/systems", AddSystem)
		authenticated.DELETE("/systems/:id", DeleteSystem)
		authenticated.GET("/systems/:id/proxy/*path", ProxyRequest)
	}
}

// --- Auth Handlers ---

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	apiKey, err := auth.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": apiKey})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	apiKey, err := auth.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": apiKey})
}

// --- System Handlers ---

func AddSystem(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := user.(*auth.User).ID
	userAPIKey := user.(*auth.User).APIKey

	var req struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		APIKey string `json:"api_key"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For Push agents, we use the User's API Key.
	// But for Pull agents (which we are adding here), we need the AGENT'S API Key.
	// The frontend sends the Agent's API Key in the request body.
	// If it's empty, we might default to User's key (for push), but let's prefer the one sent.
	systemAPIKey := req.APIKey
	systemAPIKey := req.APIKey
	if systemAPIKey == "" {
		systemAPIKey = userAPIKey
	}
	
	log.Printf("Adding System: Name=%s, URL=%s, Key=%s", req.Name, req.URL, systemAPIKey)

	id, err := db.AddSystem(userID, req.Name, req.URL, systemAPIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add system"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "name": req.Name, "url": req.URL, "api_key": userAPIKey})
}

func GetSystems(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*auth.User).ID

	systems, err := db.GetSystems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch systems"})
		return
	}
	c.JSON(http.StatusOK, systems)
}

func DeleteSystem(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*auth.User).ID
	idStr := c.Param("id")
	
	// Convert to int to be safe
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := db.DeleteSystem(userID, id); err != nil {
		log.Printf("DeleteSystem Failed: UserID=%d, ID=%d, Error=%v", userID, id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete system"})
		return
	}
	c.Status(http.StatusOK)
}

// --- Metrics Handlers ---

func IngestMetrics(c *gin.Context) {
	// 1. Validate API Key
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API Key"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key format"})
		return
	}
	apiKey := parts[1]

	// Check if this API Key belongs to a User
	_, err := auth.GetUserByAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
		return
	}

	// 2. Parse Metrics
	var data metrics.SystemMetrics
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metrics data"})
		return
	}

	// 3. Store Metrics
	// We need a unique ID for this agent.
	// The agent should send a unique ID, or we use Hostname.
	// For now, let's use the API Key + Hostname as the key in GlobalStore?
	// Or just API Key if we assume 1 User = 1 Server? No, user has multiple servers.
	// The Store needs to support multiple servers per user.
	// Let's use the API Key for now to verify connectivity, but ideally we need a unique Agent ID.
	// Update: The Agent sends "HostInfo". We can use Hostname.
	// But Hostnames can collide.
	// Let's use "APIKey:Hostname" as the key in GlobalStore.
	
	storeKey := fmt.Sprintf("%s:%s", apiKey, data.HostInfo.Hostname)
	
	// Also, we need to auto-register the system if it doesn't exist?
	// The "One Command" adds the system to the DB via the Frontend *before* the agent starts?
	// No, the user runs the command, the agent starts pushing.
	// If the system is not in the DB, it won't show up in the dashboard.
	// We should AUTO-REGISTER the system if it's new!
	// This is a great feature for "Zero Touch".
	
	// Check if system exists for this user with this hostname?
	// We don't have hostname in DB. We have Name.
	// Let's check if we have a system with this API Key and Name = Hostname.
	// If not, create it.
	
	// For now, let's just store it. The Frontend will query by System ID.
	// If the system is in the DB, we can map it.
	
	metrics.GlobalStore.Update(storeKey, data)
	
	// Auto-registration logic (Optional but cool)
	// systems, _ := db.GetSystems(user.ID)
	// found := false
	// for _, s := range systems {
	// 	if s.Name == data.HostInfo.Hostname { found = true; break }
	// }
	// if !found {
	// 	db.AddSystem(user.ID, data.HostInfo.Hostname, "dynamic", apiKey)
	// }

	c.Status(http.StatusOK)
}

func GetMetrics(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*auth.User).ID

	systemIDStr := c.Query("system_id")
	if systemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "system_id is required"})
		return
	}

	system, err := db.GetSystem(userID, systemIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	// Push Mode
	if system.URL == "push" || system.URL == "dynamic" {
		// We need to reconstruct the store key.
		// We assumed "APIKey:Hostname".
		// But we don't know the Hostname here, we only have the System Name.
		// If we enforced Name == Hostname, it works.
		storeKey := fmt.Sprintf("%s:%s", system.APIKey, system.Name)
		
		// Fallback: Try just API Key (for single server legacy)
		data, ok := metrics.GlobalStore.Get(storeKey)
		if !ok {
			// Try just API Key
			data, ok = metrics.GlobalStore.Get(system.APIKey)
			if !ok {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No data received yet"})
				return
			}
		}
		c.JSON(http.StatusOK, data)
		return
	}

	// Pull Mode
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", system.URL+"/api/v1/metrics", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+system.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Failed to connect to agent: %v", err)})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// If Agent returns 401, don't return 401 to frontend (triggers logout).
		// Return 502 Bad Gateway instead.
		status := resp.StatusCode
		if status == http.StatusUnauthorized {
			status = http.StatusBadGateway
		}
		c.JSON(status, gin.H{"error": fmt.Sprintf("Agent returned error: %d", resp.StatusCode)})
		return
	}

	io.Copy(c.Writer, resp.Body)
}

func ProxyRequest(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(*auth.User).ID
	
	systemIDStr := c.Param("id")
	path := c.Query("path")

	system, err := db.GetSystem(userID, systemIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	if system.URL == "push" || system.URL == "dynamic" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot proxy to push agent"})
		return
	}

	// Proxy to Agent
	client := &http.Client{Timeout: 10 * time.Second}
	targetURL := system.URL + "/api/v1" + path
	
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+system.APIKey)

	// Copy query params
	q := req.URL.Query()
	for k, v := range c.Request.URL.Query() {
		if k != "path" {
			for _, val := range v {
				q.Add(k, val)
			}
		}
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Failed to connect to agent: %v", err)})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	for k, v := range resp.Header {
		for _, val := range v {
			c.Writer.Header().Add(k, val)
		}
	}
	_, _ = io.Copy(c.Writer, resp.Body)
}
