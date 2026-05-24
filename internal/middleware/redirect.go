package middleware

import (
	"github.com/beecorrea/shortlinks/internal/service"
	"github.com/beecorrea/shortlinks/internal/templates"
	"github.com/gofiber/fiber/v2"
)

func NewRedirectMiddleware(svc *service.RedirectService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		host := ctx.Hostname()
		path := ctx.Path()

		key := svc.GetRedirectKey(host, path)
		if key == "" {
			return ctx.Next()
		}

		shortlink, err := svc.GetShortlinkByKey(key)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		if shortlink != nil {
			ctx.Set("Cache-Control", "no-store, no-cache, must-revalidate")
			return ctx.Redirect(shortlink.URL, fiber.StatusFound)
		}

		return ctx.Status(fiber.StatusNotFound).Send(templates.NotFoundHTML)
	}
}
