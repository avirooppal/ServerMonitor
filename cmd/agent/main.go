package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/logger"
	"github.com/user/server-moni/internal/metrics"
)

type program struct {
	server *http.Server
	cfg    *Config
}

type Config struct {
	ServerURL string
	APIKey    string
}

func (p *program) Start(s service.Service) error {
	logger.Info("Starting Agent Service...")
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping Agent Service...")
	if p.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return p.server.Shutdown(ctx)
	}
	return nil
}

func (p *program) run() {
	// ... (Data Dir logic remains same) ...
	// Determine Data Directory
	var dataDir string
	if runtime.GOOS == "windows" {
		dataDir = filepath.Join(os.Getenv("ProgramData"), "ServerMonitor", "data")
	} else {
		// Linux/Mac: Use relative to executable or /var/lib
		exePath, err := os.Executable()
		if err != nil {
			logger.Error("Failed to get executable path", "error", err)
			return
		}
		dataDir = filepath.Join(filepath.Dir(exePath), "data")
	}

	// Ensure data dir exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logger.Error("Failed to create data dir", "error", err, "path", dataDir)
		return
	}
	
	if err := os.Chdir(filepath.Dir(dataDir)); err != nil {
		logger.Error("Failed to change working directory", "error", err)
	}
	
	db.InitDB()
	metrics.InitStore()

	// Initialize Collector
	collector := metrics.NewCollector()
	collector.StartBackgroundTasks()

	// Start Local Collector (Self-Monitoring / Cache)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			m := collector.Collect()
			metrics.GlobalStore.Update("local", m)
		}
	}()

	// Start Pusher if configured
	// Priority: Config (Flags) > Env
	serverURL := p.cfg.ServerURL
	if serverURL == "" {
		serverURL = os.Getenv("SERVER_URL")
	}
	
	apiKey := p.cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("API_KEY")
	}

	if serverURL != "" && apiKey != "" {
		go startPusher(collector, serverURL, apiKey)
	} else {
		logger.Warn("Push mode disabled: Missing SERVER_URL or API_KEY")
	}

	// ... (Web Server logic remains same) ...
	// Setup Web Server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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

	p.server = &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	logger.Info("Agent running", "port", port)
	if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Web server error", "error", err)
	}
}

func main() {
	logger.InitLogger()
	
	// Agent specific flags
	var flagServer, flagToken, flagService string
	flag.StringVar(&flagServer, "server", "", "Server URL")
	flag.StringVar(&flagToken, "token", "", "API Key")
	flag.StringVar(&flagService, "service", "", "Service action: install, uninstall, start, stop")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "ServerMoniAgent",
		DisplayName: "Server Monitor Agent",
		Description: "Agent for Server Monitor SaaS",
		Arguments:   []string{"-server", flagServer, "-token", flagToken},
	}

	prg := &program{
		cfg: &Config{
			ServerURL: flagServer,
			APIKey:    flagToken,
		},
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error("Failed to create service", "error", err)
		os.Exit(1)
	}

	if flagService != "" {
		if err := service.Control(s, flagService); err != nil {
			logger.Error("Service control failed", "action", flagService, "error", err)
			os.Exit(1)
		}
		logger.Info("Service action completed", "action", flagService)
		return
	}

	if err := s.Run(); err != nil {
		logger.Error("Service run failed", "error", err)
	}
}

func startPusher(c *metrics.Collector, serverURL, apiKey string) {
	logger.Info("Starting Push Mode", "url", serverURL)
	client := &http.Client{Timeout: 5 * time.Second}
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		m := c.Collect()
		
		data, err := json.Marshal(m)
		if err != nil {
			logger.Error("Error marshaling metrics", "error", err)
			continue
		}

		req, err := http.NewRequest("POST", serverURL+"/api/v1/ingest", bytes.NewBuffer(data))
		if err != nil {
			logger.Error("Error creating request", "error", err)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			logger.Error("Error pushing metrics", "error", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			logger.Warn("Error pushing metrics", "status", resp.StatusCode)
		}
	}
}
