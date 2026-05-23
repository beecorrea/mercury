package controllers

import (
	"github.com/beecorrea/shortlinks/internal/templates"
	"github.com/gofiber/fiber/v2"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// ServeDashboard serves the embedded HTML admin dashboard.
func (c *DashboardController) ServeDashboard(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html; charset=utf-8")
	ctx.Set("Cache-Control", "no-store, no-cache, must-revalidate")
	return ctx.Send(templates.IndexHTML)
}
