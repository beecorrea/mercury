package main

import (
	"log"

	"github.com/beecorrea/shortlinks/internal/server"
)

func main() {
	dbPath := "mercury.db"
	srv, err := server.NewServer(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer srv.Close()

	addr := "0.0.0.0:45800"
	log.Printf("Starting Mercury Shortener Service with Fiber...")
	log.Printf("Dashboard accessible at http://localhost:45800")
	log.Printf("Shortlink format: http://mercury.<domain>:45800/<key>")

	if err := srv.Start(addr); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
