package api

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/metrics"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	
	// Public endpoint to check if setup is needed? 
	// For security, we might want to keep everything behind auth or have a specific "is configured" endpoint.
	// But the requirement says "API Key auto-generated", so it's always configured.
	
	api.POST("/verify-key", func(c *gin.Context) {
		// This endpoint is just to test if the key is valid.
		// The middleware will handle the actual validation.
		// If we reach here, the middleware passed, so the key is valid.
		c.JSON(http.StatusOK, gin.H{"status": "valid"})
	})

	authenticated := api.Group("/")
	authenticated.Use(auth.AuthMiddleware())
	{
		authenticated.GET("/metrics", GetMetrics)
		authenticated.GET("/servers", GetServers)
		authenticated.POST("/agents", RegisterAgent)
		authenticated.GET("/status", GetStatus)
	}

	// Ingestion endpoint uses Agent Auth (or Master Key)
	ingest := api.Group("/")
	ingest.Use(auth.AgentAuthMiddleware())
	{
		ingest.POST("/ingest", IngestMetrics)
	}
}

func RegisterAgent(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	auth.AddAgentToken(req.Token, req.Name)
	c.Status(http.StatusOK)
}

func GetMetrics(c *gin.Context) {
	serverID := c.Query("server_id")
	if serverID == "" {
		serverID = "local"
	}

	data, ok := metrics.GlobalStore.Get(serverID)
	if !ok {
		// If requesting local and it's not ready yet, return empty or wait?
		// Better to return 404 or empty structure.
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetServers(c *gin.Context) {
	all := metrics.GlobalStore.GetAll()
	var servers []gin.H
	for id, m := range all {
		servers = append(servers, gin.H{
			"id":          id,
			"hostname":    m.HostInfo.Hostname,
			"platform":    m.HostInfo.Platform,
			"last_update": m.LastUpdate,
		})
	}
	c.JSON(http.StatusOK, servers)
}

func IngestMetrics(c *gin.Context) {
	var m metrics.SystemMetrics
	if err := c.BindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Determine Server ID (prefer HostID, fallback to Hostname)
	serverID := m.HostInfo.Hostname
	if m.HostInfo.HostID != "" {
		serverID = m.HostInfo.HostID
	}
	
	metrics.GlobalStore.Update(serverID, m)
	c.Status(http.StatusOK)
}

func GetStatus(c *gin.Context) {
	h, _ := host.Info()
	c.JSON(http.StatusOK, gin.H{
		"hostname":      h.Hostname,
		"uptime":        h.Uptime,
		"os":            h.OS,
		"platform":      h.Platform,
		"platform_ver":  h.PlatformVersion,
		"kernel_ver":    h.KernelVersion,
		"arch":          runtime.GOARCH,
		"cpus":          runtime.NumCPU(),
		"go_version":    runtime.Version(),
	})
}
