package controllers

import (
	"net/url"
	"strings"

	"github.com/beecorrea/shortlinks/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ShortlinkController struct {
	Service *service.ShortlinkService
}

func NewShortlinkController(svc *service.ShortlinkService) *ShortlinkController {
	return &ShortlinkController{Service: svc}
}

// ListShortlinks handles GET /api/links.
func (c *ShortlinkController) ListShortlinks(ctx *fiber.Ctx) error {
	shortlinks, err := c.Service.ListShortlinks()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch shortlinks"})
	}
	return ctx.JSON(shortlinks)
}

type ShortenRequest struct {
	Key    string `json:"key"`
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

// CreateShortlink handles POST /api/shorten.
func (c *ShortlinkController) CreateShortlink(ctx *fiber.Ctx) error {
	var req ShortenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Key = strings.TrimSpace(req.Key)
	req.Domain = strings.TrimSpace(req.Domain)
	req.URL = strings.TrimSpace(req.URL)

	if req.Key == "" || req.Domain == "" || req.URL == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "key, domain, and url are required"})
	}

	// Validate target URL format
	u, err := url.ParseRequestURI(req.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "destination must be a valid absolute URL (e.g. https://google.com)"})
	}

	// Persist shortlink using the service
	err = c.Service.CreateShortlink(req.Key, req.URL, req.Domain)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "short key is already in use"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to persist shortlink"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "shortlink created successfully"})
}

// DeleteShortlink handles DELETE /api/links/:key.
func (c *ShortlinkController) DeleteShortlink(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "key parameter is required"})
	}

	err := c.Service.DeleteShortlink(key)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete shortlink"})
	}

	return ctx.JSON(fiber.Map{"message": "shortlink deleted successfully"})
}
