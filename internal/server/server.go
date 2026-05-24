package server

import (
	"fmt"

	"github.com/beecorrea/shortlinks/internal/config"
	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/scraper"
	"github.com/beecorrea/shortlinks/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
	db  *database.RedirectDB
}

func NewServer(cfg *config.Config) (*Server, error) {
	db, err := database.NewRedirectDB(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("initializing database: %w", err)
	}

	redirectSvc := service.NewRedirectService(db, cfg.Domain)
	scraperImpl := scraper.NewHTTPScraper()
	shortlinkSvc := service.NewShortlinkService(db, scraperImpl)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	s := &Server{
		App: app,
		db:  db,
	}

	s.setupRoutes(cfg, redirectSvc, shortlinkSvc)

	return s, nil
}

func (s *Server) Start(addr string) error {
	return s.App.Listen(addr)
}

func (s *Server) Close() error {
	return s.db.Close()
}

