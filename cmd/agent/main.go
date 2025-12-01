package main

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

func main() {
	// Optimize for Low RAM
	debug.SetGCPercent(50) // Aggressive GC (default is 100)

	// Initialize DB (SQLite)
	db.InitDB()

	// Initialize Auth (Check for AGENT_SECRET env var)
	if secret := os.Getenv("AGENT_SECRET"); secret != "" {
		auth.SetStaticKey(secret)
		log.Printf("Agent initialized with Static Token: %s", secret)
	} else {
		auth.InitAuth()
	}

	// Initialize Metric Store (In-Memory for now)
	metrics.InitStore()

	// Initialize Collector
	collector := metrics.NewCollector()
	// collector.StartBackgroundTasks() // Disable heavy background tasks (Disk Usage) for low RAM

	// Start Collector
	go func() {
		// Default interval 5s (slower for less CPU/RAM churn)
		interval := 5 * time.Second
		if val := os.Getenv("COLLECTION_INTERVAL_SECONDS"); val != "" {
			if i, err := strconv.Atoi(val); err == nil {
				interval = time.Duration(i) * time.Second
			}
		}

		ticker := time.NewTicker(interval)
		for range ticker.C {
			m := collector.Collect()
			// Update "local" metrics
			metrics.GlobalStore.Update("local", m)
			
			// Force GC every few cycles if needed, but SetGCPercent should handle it.
			// runtime.GC() // Optional: Manual GC if really tight
		}
	}()

	// Setup Web Server
	gin.SetMode(gin.ReleaseMode) // Release mode for less memory
	r := gin.New() // Use New() instead of Default() to avoid some middleware overhead if desired
	r.Use(gin.Recovery())

	// CORS - Allow All
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
	api := r.Group("/api/v1")
	
	// Public check
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Metrics Endpoint (Protected by Agent Token if configured, or Open)
	api.GET("/metrics", auth.AgentAuthMiddleware(), func(c *gin.Context) {
		data, ok := metrics.GlobalStore.Get("local")
		if !ok {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Collecting metrics..."})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	// ... (Keep other endpoints but ensure they use AgentAuthMiddleware)
	// Docker Logs
	api.GET("/docker/containers/:id/logs", auth.AgentAuthMiddleware(), func(c *gin.Context) {
		id := c.Param("id")
		tail := c.DefaultQuery("tail", "100")
		logs, err := collector.GetContainerLogs(id, tail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"logs": logs})
	})

	// Fail2Ban Stats
	api.GET("/security/fail2ban", auth.AgentAuthMiddleware(), func(c *gin.Context) {
		path := c.DefaultQuery("path", "/var/log/fail2ban.log")
		stats, err := collector.GetFail2BanStats(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	})

	// Disk Usage (Disabled/Cached)
	api.GET("/disk/usage", auth.AgentAuthMiddleware(), func(c *gin.Context) {
		// Return empty or cached
		usage := collector.GetCachedDiskUsage()
		c.JSON(http.StatusOK, usage)
	})

	// Auth Logs
	api.GET("/security/logins", auth.AgentAuthMiddleware(), func(c *gin.Context) {
		path := c.DefaultQuery("path", "/var/log/auth.log")
		logs, err := collector.GetAuthLogs(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, logs)
	})

	// Disk History
	api.GET("/disk/history", auth.AgentAuthMiddleware(), func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "30")
		limit, _ := strconv.Atoi(limitStr)
		history, err := db.GetDiskHistory(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, history)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Agent running on port %s (Low RAM Mode)", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

