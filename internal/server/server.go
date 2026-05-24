package server

import (
	"context"
	"fmt"

	"github.com/beecorrea/shortlinks/internal/config"
	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/scraper"
	"github.com/beecorrea/shortlinks/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App    *fiber.App
	db     *database.RedirectDB
	cancel context.CancelFunc
}

func NewServer(cfg *config.Config) (*Server, error) {
	db, err := database.NewRedirectDB(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("initializing database: %w", err)
	}

	redirectSvc := service.NewRedirectService(db, cfg.Domain)
	scraperImpl := scraper.NewHTTPScraper()
	shortlinkSvc := service.NewShortlinkService(db, scraperImpl)

	retrier := service.NewScrapeRetrier(db, scraperImpl, cfg.ScrapeRetryInterval, cfg.ScrapeMaxAttempts)
	ctx, cancel := context.WithCancel(context.Background())
	go retrier.Start(ctx)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	s := &Server{
		App:    app,
		db:     db,
		cancel: cancel,
	}

	s.setupRoutes(cfg, redirectSvc, shortlinkSvc)

	return s, nil
}

func (s *Server) Start(addr string) error {
	return s.App.Listen(addr)
}

func (s *Server) Close() error {
	s.cancel()
	return s.db.Close()
}

