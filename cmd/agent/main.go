package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/user/server-moni/internal/metrics"
)

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "URL of the central server")
	apiKey := flag.String("key", "", "API Key for authentication")
	interval := flag.Duration("interval", 2*time.Second, "Collection interval")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("API Key is required. Use -key <your-api-key>")
	}

	log.Printf("Starting Agent...")
	log.Printf("Server: %s", *serverURL)
	log.Printf("Interval: %v", *interval)

	collector := metrics.NewCollector()
	client := &http.Client{Timeout: 10 * time.Second}

	ticker := time.NewTicker(*interval)
	for range ticker.C {
		// Collect
		m := collector.Collect()

		// Send
		payload, err := json.Marshal(m)
		if err != nil {
			log.Printf("Error marshaling metrics: %v", err)
			continue
		}

		req, err := http.NewRequest("POST", *serverURL+"/api/v1/ingest", bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+*apiKey)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending metrics: %v", err)
			continue
		}
		
		if resp.StatusCode != http.StatusOK {
			log.Printf("Server returned error: %s", resp.Status)
		}
		resp.Body.Close()
	}
}
