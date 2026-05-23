package middleware

import (
	"strings"

	"github.com/beecorrea/shortlinks/internal/service"
	"github.com/beecorrea/shortlinks/internal/templates"
	"github.com/gofiber/fiber/v2"
)

func NewRedirectMiddleware(svc *service.RedirectService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		host := ctx.Hostname()

		// Check if Host starts with "mercury."
		if host == "mercury" || strings.HasPrefix(host, "mercury.") {
			path := ctx.Path()
			if strings.HasPrefix(path, "/api/") {
				return ctx.Next()
			}

			trimmedPath := strings.Trim(path, "/")

			// If path contains sub-routes, use the first segment as the key
			key := trimmedPath
			if idx := strings.Index(trimmedPath, "/"); idx != -1 {
				key = trimmedPath[:idx]
			}

			if key != "" {
				shortlink, err := svc.GetShortlinkByKey(key)
				if err != nil {
					return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
				}

				if shortlink != nil {
					ctx.Set("Cache-Control", "no-store, no-cache, must-revalidate")
					return ctx.Redirect(shortlink.URL, fiber.StatusFound)
				}

				// Key not found - display premium 404 page
				ctx.Status(fiber.StatusNotFound)
				ctx.Set("Content-Type", "text/html; charset=utf-8")
				return ctx.Send(templates.NotFoundHTML)
			}
		}

		return ctx.Next()
	}
}
