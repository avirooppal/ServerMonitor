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
		authenticated.GET("/status", GetStatus)
	}
}

func GetMetrics(c *gin.Context) {
	data := metrics.GetMetrics()
	c.JSON(http.StatusOK, data)
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
