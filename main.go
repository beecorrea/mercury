package main

import (
	"log"
	"os"

	"github.com/beecorrea/shortlinks/internal/server"
)

func main() {
	dbPath := os.Getenv("MERCURY_DB_PATH")
	if dbPath == "" {
		dbPath = "mercury.db"
	}
	srv, err := server.NewServer(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer srv.Close()

	addr := "0.0.0.0:45800"
	log.Printf("Starting Mercury Shortener Service with Fiber...")
	log.Printf("Dashboard accessible at http://%s", addr)
	log.Printf("Shortlink format: http://mercury.<domain>/<key> (e.g., http://mercury.communist.mom/<key>)")

	if err := srv.Start(addr); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
