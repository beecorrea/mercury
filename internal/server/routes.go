package server

import (
	"github.com/beecorrea/shortlinks/internal/controllers"
	"github.com/beecorrea/shortlinks/internal/middleware"
	"github.com/beecorrea/shortlinks/internal/service"
)

func (s *Server) setupRoutes(redirectSvc *service.RedirectService, shortlinkSvc *service.ShortlinkService) {
	// Apply redirect middleware
	s.App.Use(middleware.NewRedirectMiddleware(redirectSvc))

	// Initialize Controllers
	shortlinkController := controllers.NewShortlinkController(shortlinkSvc)
	dashboardController := controllers.NewDashboardController()

	// Register Routes
	dashboard := s.App.Group("/")
	dashboard.Get("/", dashboardController.ServeDashboard)

	api := s.App.Group("/api")
	api.Get("/links", shortlinkController.ListShortlinks)
	api.Post("/shorten", shortlinkController.CreateShortlink)
	api.Delete("/links/:key", shortlinkController.DeleteShortlink)
}
