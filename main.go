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

	port := ":45800"
	log.Printf("Starting Mercury Shortener Service with Fiber...")
	log.Printf("Dashboard accessible at http://localhost%s", port)
	log.Printf("Shortlink format: http://mercury.<domain>%s/<key>", port)

	if err := srv.Start(port); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
