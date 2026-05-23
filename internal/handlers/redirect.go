package handlers

import (
	"strings"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/templates"
	"github.com/gofiber/fiber/v2"
)

func NewRedirectMiddleware(db *database.RedirectDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := c.Hostname()

		// Check if Host starts with "mercury."
		if strings.HasPrefix(host, "mercury.") {
			path := strings.Trim(c.Path(), "/")

			// If path contains sub-routes, use the first segment as the key
			key := path
			if idx := strings.Index(path, "/"); idx != -1 {
				key = path[:idx]
			}

			if key != "" {
				link, err := db.GetLinkByKey(key)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
				}

				if link != nil {
					c.Set("Cache-Control", "no-store, no-cache, must-revalidate")
					return c.Redirect(link.URL, fiber.StatusFound)
				}
			}

			// Key not found - display premium 404 page
			c.Status(fiber.StatusNotFound)
			c.Set("Content-Type", "text/html; charset=utf-8")
			return c.Send(templates.NotFoundHTML)
		}

		return c.Next()
	}
}
