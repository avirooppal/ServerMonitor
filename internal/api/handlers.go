package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	
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
		authenticated.GET("/systems/:id/proxy", ProxyRequest)
	}
}

func AddSystem(c *gin.Context) {
	var req struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		APIKey string `json:"api_key"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := db.AddSystem(req.Name, req.URL, req.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add system"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetSystems(c *gin.Context) {
	systems, err := db.GetSystems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch systems"})
		return
	}
	c.JSON(http.StatusOK, systems)
}

func DeleteSystem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := db.DeleteSystem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete system"})
		return
	}
	c.Status(http.StatusOK)
}

func GetMetrics(c *gin.Context) {
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

	// Stream response back
	c.Status(resp.StatusCode)
	for k, v := range resp.Header {
		for _, val := range v {
			c.Writer.Header().Add(k, val)
		}
	}
	_, _ = io.Copy(c.Writer, resp.Body)
}
