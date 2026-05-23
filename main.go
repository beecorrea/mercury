package main

import (
	"log"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	dbPath := "mercury.db"
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Host Redirect Middleware
	app.Use(handlers.NewRedirectMiddleware(db))

	adminHandler := handlers.NewAdminHandler(db)

	// Dashboard
	app.Get("/", adminHandler.ServeDashboard)

	// API
	app.Get("/api/links", adminHandler.ListLinks)
	app.Post("/api/shorten", adminHandler.CreateLink)
	app.Delete("/api/links/:key", adminHandler.DeleteLink)

	port := ":45800"
	log.Printf("Starting Mercury Shortener Service with Fiber...")
	log.Printf("Dashboard accessible at http://localhost%s", port)
	log.Printf("Shortlink format: http://mercury.<domain>%s/<key>", port)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
