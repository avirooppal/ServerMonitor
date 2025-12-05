package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

var logger service.Logger

type program struct {
	ServerURL string
	Token     string
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	// Initialize DB (SQLite) - For Agent, we might not need full DB if just pushing?
	// But existing code uses it for Disk History.
	// Ensure data dir exists
	os.MkdirAll("data", 0755)
	db.InitDB()

	// Initialize Metric Store
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
	serverURL := os.Getenv("SERVER_URL")
	apiKey := os.Getenv("API_KEY")

	if serverURL == "" {
		serverURL = p.ServerURL
	}
	if apiKey == "" {
		apiKey = p.Token
	}

	if serverURL != "" && apiKey != "" {
		go startPusher(collector, serverURL, apiKey)
	}

	// Setup Web Server
	r := gin.Default()
	// ... (cors config) ...
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

	// Protected Metrics Endpoint (Pull Mode)
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

	// Only run web server if NOT controlling service or if running as service
	go func() {
		log.Printf("Agent running on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	// Parse flags first to get config for service arguments
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
		ServerURL: flagServer,
		Token:     flagToken,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Handle Service Control Flags
	if flagService != "" {
		switch flagService {
		case "install":
			err = s.Install()
			if err != nil {
				log.Fatalf("Failed to install: %v", err)
			}
			log.Println("Service installed")
			return
		case "uninstall":
			err = s.Uninstall()
			if err != nil {
				log.Fatalf("Failed to uninstall: %v", err)
			}
			log.Println("Service uninstalled")
			return
		case "start":
			err = s.Start()
			if err != nil {
				log.Fatalf("Failed to start: %v", err)
			}
			log.Println("Service started")
			return
		case "stop":
			err = s.Stop()
			if err != nil {
				log.Fatalf("Failed to stop: %v", err)
			}
			log.Println("Service stopped")
			return
		}
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func startPusher(c *metrics.Collector, serverURL, apiKey string) {
	log.Printf("Starting Push Mode to %s", serverURL)
	client := &http.Client{Timeout: 5 * time.Second}
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		m := c.Collect()
		
		data, err := json.Marshal(m)
		if err != nil {
			log.Printf("Error marshaling metrics: %v", err)
			continue
		}

		req, err := http.NewRequest("POST", serverURL+"/api/v1/ingest", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error pushing metrics: %v", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error pushing metrics: server returned %d", resp.StatusCode)
		}
	}
}
