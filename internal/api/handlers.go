package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	
	// Public Routes (No Auth)
	api.GET("/metrics", GetMetrics)
	api.GET("/systems", GetSystems)
	api.POST("/systems", AddSystem)
	api.DELETE("/systems/:id", DeleteSystem)
	api.GET("/systems/:id/proxy", ProxyRequest)
}

// --- System Handlers ---

func AddSystem(c *gin.Context) {
	// Public Dashboard: User ID is always 0 (Global)
	userID := 0
	
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
	userID := 0
	
	systems, err := db.GetSystems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch systems"})
		return
	}
	c.JSON(http.StatusOK, systems)
}

func DeleteSystem(c *gin.Context) {
	userID := 0
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
	userID := 0
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

	system, err := db.GetSystem(systemID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	// Proxy to Agent
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
		c.JSON(resp.StatusCode, gin.H{"error": "Agent returned error"})
		return
	}

	var data metrics.SystemMetrics
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode agent response"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func ProxyRequest(c *gin.Context) {
	userID := 0
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

	system, err := db.GetSystem(systemID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	// Construct target URL
	targetURL := fmt.Sprintf("%s%s", system.URL, path)
	
	// Create request
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for k, v := range c.Request.Header {
		req.Header[k] = v
	}
	
	// Add Agent Auth
	req.Header.Set("Authorization", "Bearer "+system.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Failed to connect to agent: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
	c.Status(resp.StatusCode)
	
	// Copy body
	io.Copy(c.Writer, resp.Body)
}
