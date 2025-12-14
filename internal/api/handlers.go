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

	// Public Auth Routes
	api.POST("/auth/register", Register)
	api.POST("/auth/login", Login)
	
	// Health Check
	r.GET("/health", HealthCheck)
	api.GET("/health", HealthCheck)
	
	// Ingestion (Agent Push) - Validates System API Key internally
	api.POST("/ingest", IngestMetrics)

	// Protected Routes (User UI)
	protected := api.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/systems", GetSystems)
		protected.POST("/systems", AddSystem)
		protected.DELETE("/systems/:id", DeleteSystem)
		protected.GET("/metrics", GetMetrics)
		protected.GET("/systems/:id/proxy", ProxyRequest)
		protected.POST("/auth/logout", Logout)
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
}

// Auth Handlers

func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := auth.Register(req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user (email might be taken)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "registered"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 {
			auth.Logout(parts[1])
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}

// System Handlers

func AddSystem(c *gin.Context) {
	userID := c.GetInt("userID")
	var req struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		APIKey string `json:"api_key"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := db.AddSystem(userID, req.Name, req.URL, strings.TrimSpace(req.APIKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add system"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetSystems(c *gin.Context) {
	userID := c.GetInt("userID")
	systems, err := db.GetSystems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch systems"})
		return
	}
	c.JSON(http.StatusOK, systems)
}

func DeleteSystem(c *gin.Context) {
	userID := c.GetInt("userID")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := db.DeleteSystem(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete system"})
		return
	}
	c.Status(http.StatusOK)
}

func GetMetrics(c *gin.Context) {
	userID := c.GetInt("userID")
	systemIDStr := c.Query("system_id")
	if systemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "system_id is required"})
		return
	}

	systemID, err := strconv.Atoi(systemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system_id"})
		return
	}

	system, err := db.GetSystem(systemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	// Verify Ownership
	if system.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if this is a Push-based system (URL is "push" or "dynamic" or empty/placeholder)
	// For SaaS, we assume mostly Push.
	if system.URL == "push" || system.URL == "dynamic" || true { // Force Push check for now as default
		data, ok := metrics.GlobalStore.Get(strconv.Itoa(system.ID))
		if !ok {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No metrics received yet"})
			return
		}
		c.JSON(http.StatusOK, data)
		return
	}

	// Proxy logic (Legacy Pull) - Kept for backward compatibility if needed
	// ...
}

func IngestMetrics(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	apiKey := parts[1]

	// Validate API Key
	system, err := db.GetSystemByAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
		return
	}

	var metricsData metrics.SystemMetrics
	if err := c.BindJSON(&metricsData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update Store
	metrics.GlobalStore.Update(strconv.Itoa(system.ID), metricsData)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ProxyRequest(c *gin.Context) {
	userID := c.GetInt("userID")
	systemIDStr := c.Param("id")
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	systemID, err := strconv.Atoi(systemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system_id"})
		return
	}

	system, err := db.GetSystem(systemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	if system.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Proxy only works if Agent is directly accessible (Pull mode).
	// For Push mode (SaaS), we can't proxy unless we implement a reverse tunnel.
	// For now, return error if Push mode.
	if system.URL == "push" || system.URL == "dynamic" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Proxy not supported for Push agents yet"})
		return
	}

	// Proxy to Agent (Legacy)
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
