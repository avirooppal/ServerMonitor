package main

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/user/server-moni/internal/api"
	"github.com/user/server-moni/internal/config"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/logger"
	"github.com/user/server-moni/internal/metrics"
)

//go:embed all:dist
var staticFiles embed.FS

type program struct {
	server *http.Server
}

func (p *program) Start(s service.Service) error {
	logger.Info("Starting Service...")
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping Service...")
	if p.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return p.server.Shutdown(ctx)
	}
	return nil
}

func (p *program) run() {
	// Initialize Components
	db.InitDB()
	metrics.InitStore()

	// Start Local Collector
	go func() {
		collector := metrics.NewCollector()
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			m := collector.Collect()
			metrics.GlobalStore.Update("local", m)
		}
	}()

	// Setup Web Server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// Custom Logger Middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"ip", c.ClientIP(),
		)
	})

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	api.RegisterRoutes(r)
	r.Static("/downloads", "./downloads")
	// Serve install.sh dynamically
	r.GET("/install.sh", func(c *gin.Context) {
		content, err := os.ReadFile("./downloads/install_agent_linux.sh")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load install script")
			return
		}
		
		// Determine Server URL (Scheme + Host)
		scheme := "http"
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		host := c.Request.Host
		serverURL := scheme + "://" + host

		// Inject URL
		script := strings.Replace(string(content), "__SERVER_URL__", serverURL, -1)
		
		c.Header("Content-Type", "text/x-shellscript")
		c.String(http.StatusOK, script)
	})

	// Frontend Serving
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		logger.Warn("Frontend dist not found in embed, serving API only.")
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Server Monitor API is running (Frontend not found)"})
		})
	} else {
		logger.Info("Frontend dist found, serving SPA.")
		assetsFS, _ := fs.Sub(distFS, "assets")
		r.StaticFS("/assets", http.FS(assetsFS))
		r.GET("/vite.svg", func(c *gin.Context) {
			c.FileFromFS("vite.svg", http.FS(distFS))
		})
		r.GET("/get-key.sh", func(c *gin.Context) {
			c.FileFromFS("get-key.sh", http.FS(distFS))
		})
		r.GET("/", func(c *gin.Context) {
			c.FileFromFS("index.html", http.FS(distFS))
		})
		r.NoRoute(func(c *gin.Context) {
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

	p.server = &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: r,
	}

	logger.Info("Server starting", "port", config.AppConfig.Port)
	if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed", "error", err)
	}
}

func main() {
	logger.InitLogger()
	config.Load()

	svcConfig := &service.Config{
		Name:        "ServerMoni",
		DisplayName: "Server Monitor",
		Description: "Server Monitor Backend & Dashboard",
		Arguments:   []string{"-port", config.AppConfig.Port},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error("Failed to create service", "error", err)
		os.Exit(1)
	}

	if config.AppConfig.Service != "" {
		if err := service.Control(s, config.AppConfig.Service); err != nil {
			logger.Error("Service control failed", "action", config.AppConfig.Service, "error", err)
			os.Exit(1)
		}
		logger.Info("Service action completed", "action", config.AppConfig.Service)
		return
	}

	if err := s.Run(); err != nil {
		logger.Error("Service run failed", "error", err)
	}
}
