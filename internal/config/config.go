package config

import (
	"flag"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	Service     string // install, uninstall, start, stop
}

var AppConfig Config

func Load() {
	// Defaults
	AppConfig.Port = "8080"
	AppConfig.DatabaseURL = "server-moni.db"

	// Flags
	flag.StringVar(&AppConfig.Port, "port", "8080", "Server Port")
	flag.StringVar(&AppConfig.Service, "service", "", "Service action: install, uninstall, start, stop")
	flag.Parse()

	// Env Overrides
	if envPort := os.Getenv("PORT"); envPort != "" {
		AppConfig.Port = envPort
	}
}
