package server

import (
	"github.com/beecorrea/shortlinks/internal/config"
	"github.com/beecorrea/shortlinks/internal/controllers"
	"github.com/beecorrea/shortlinks/internal/middleware"
	"github.com/beecorrea/shortlinks/internal/service"
)

func (s *Server) setupRoutes(cfg *config.Config, redirectSvc *service.RedirectService, shortlinkSvc *service.ShortlinkService) {
	// Apply redirect middleware
	s.App.Use(middleware.NewRedirectMiddleware(redirectSvc))

	// Initialize Controllers
	shortlinkController := controllers.NewShortlinkController(shortlinkSvc, cfg.Domain, cfg.Port)
	dashboardController := controllers.NewDashboardController()

	// Register Routes
	s.App.Get("/", dashboardController.ServeDashboard)
	s.App.Get("/api/links", shortlinkController.ListShortlinks)
	s.App.Post("/api/shorten", shortlinkController.CreateShortlink)
	s.App.Delete("/api/links/:key", shortlinkController.DeleteShortlink)
}
