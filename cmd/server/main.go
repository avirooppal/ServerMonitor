package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/server-moni/internal/api"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

//go:embed all:dist
var staticFiles embed.FS

func main() {
	// Initialize DB
	db.InitDB()

	// Initialize Auth (Generate Key if needed)


	// Initialize Metric Store
	metrics.InitStore()

	// Start Local Collector (Self-Monitoring)
	go func() {
		collector := metrics.NewCollector()
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			m := collector.Collect()
			metrics.GlobalStore.Update("local", m)
		}
	}()

	// Setup Web Server
	r := gin.Default()

	// CORS for development (allow all for now, tighten later if needed)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
	api.RegisterRoutes(r)

	// Serve Downloads (Agent Binaries & Scripts)
	r.Static("/downloads", "./downloads")
	// Also serve install.sh directly from downloads for convenience
	r.StaticFile("/install.sh", "./downloads/install_agent_linux.sh")


	// Serve Frontend (Embedded)
	// We assume the frontend is built into 'dist' folder and embedded.
	// If 'dist' is not found (dev mode), we might want to skip or serve from local dir.
	
	// Check if we are running in dev mode without embed
	// For this implementation, we will try to serve from embed.
	// If the embed is empty (during dev before build), we might fallback.
	
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Println("Frontend dist not found in embed, serving API only.")
	} else {
		assetsFS, _ := fs.Sub(distFS, "assets")
		r.StaticFS("/assets", http.FS(assetsFS))
		r.GET("/vite.svg", func(c *gin.Context) {
			c.FileFromFS("vite.svg", http.FS(distFS))
		})
		r.GET("/install.sh", func(c *gin.Context) {
			c.FileFromFS("install.sh", http.FS(distFS))
		})
		r.GET("/get-key.sh", func(c *gin.Context) {
			c.FileFromFS("get-key.sh", http.FS(distFS))
		})
		r.NoRoute(func(c *gin.Context) {
			// Serve index.html for SPA routing
			file, err := distFS.Open("index.html")
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Frontend not found"})
				return
			}
			defer file.Close()
			stat, _ := file.Stat()
			http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), file.(io.ReadSeeker))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
