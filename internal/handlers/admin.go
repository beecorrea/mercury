package handlers

import (
	"net/url"
	"strings"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/templates"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	DB *database.DB
}

func NewAdminHandler(db *database.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

// ServeDashboard serves the embedded HTML admin dashboard.
func (h *AdminHandler) ServeDashboard(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(templates.IndexHTML)
}

// ListLinks handles GET /api/links.
func (h *AdminHandler) ListLinks(c *fiber.Ctx) error {
	links, err := h.DB.ListLinks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch links"})
	}
	return c.JSON(links)
}

type ShortenRequest struct {
	Key    string `json:"key"`
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

// CreateLink handles POST /api/shorten.
func (h *AdminHandler) CreateLink(c *fiber.Ctx) error {
	var req ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Key = strings.TrimSpace(req.Key)
	req.Domain = strings.TrimSpace(req.Domain)
	req.URL = strings.TrimSpace(req.URL)

	if req.Key == "" || req.Domain == "" || req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "key, domain, and url are required"})
	}

	// Validate target URL format
	u, err := url.ParseRequestURI(req.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "destination must be a valid absolute URL (e.g. https://google.com)"})
	}

	// Persist link in SQLite
	err = h.DB.CreateLink(req.Key, req.URL, req.Domain)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "short key is already in use"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to persist link"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "link created successfully"})
}

// DeleteLink handles DELETE /api/links/:key.
func (h *AdminHandler) DeleteLink(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "key parameter is required"})
	}

	err := h.DB.DeleteLink(key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete link"})
	}

	return c.JSON(fiber.Map{"message": "link deleted successfully"})
}
