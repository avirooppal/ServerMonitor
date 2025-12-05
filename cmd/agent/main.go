package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/user/server-moni/internal/auth"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

var logger service.Logger

type program struct{}

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

	// Initialize Auth
	auth.InitAuth()

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

	// Check flags if env vars are empty
	var flagServer, flagToken string
	flag.StringVar(&flagServer, "server", "", "Server URL")
	flag.StringVar(&flagToken, "token", "", "API Key")
	flag.Parse()

	if serverURL == "" {
		serverURL = flagServer
	}
	if apiKey == "" {
		apiKey = flagToken
	}

	if serverURL != "" && apiKey != "" {
		go startPusher(collector, serverURL, apiKey)
	}

	// Setup Web Server
	r := gin.Default()
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

	// Other endpoints... (omitted for brevity, can add back if needed)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Agent running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "ServerMoniAgent",
		DisplayName: "Server Monitor Agent",
		Description: "Agent for Server Monitor SaaS",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Handle Service Control Flags
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
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
