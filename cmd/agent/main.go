package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

func main() {
	// Initialize DB (SQLite)
	db.InitDB()

	// Initialize Auth (Generate API Key)
	auth.InitAuth()

	// Initialize Metric Store (In-Memory for now)
	metrics.InitStore()

	// Start Collector
	go func() {
		collector := metrics.NewCollector()
		// Default interval 2s, can be env var
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			m := collector.Collect()
			// Update "local" metrics
			metrics.GlobalStore.Update("local", m)
		}
	}()

	// Setup Web Server
	r := gin.Default()

	// CORS - Allow All (User will likely access from cloud dashboard domain)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // TODO: Restrict to user's dashboard domain in prod
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

	// Protected Metrics Endpoint
	api.GET("/metrics", auth.AuthMiddleware(), func(c *gin.Context) {
		data, ok := metrics.GlobalStore.Get("local")
		if !ok {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Collecting metrics..."})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Agent running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
