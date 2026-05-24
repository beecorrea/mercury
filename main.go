package main

import (
	"fmt"
	"log"

	"github.com/beecorrea/shortlinks/internal/config"
	"github.com/beecorrea/shortlinks/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer srv.Close()

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("Starting Mercury Shortener Service with Fiber...")
	log.Printf("Dashboard accessible at http://%s", addr)
	log.Printf("Shortlink format: http://mercury.<domain>/<key> (e.g., http://mercury.%s/<key>)", cfg.Domain)

	if err := srv.Start(addr); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}

