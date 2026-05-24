package service

import (
	"context"
	"fmt"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
)

// ShortlinkService handles link management operations.
type ShortlinkService struct {
	db      *database.RedirectDB
	scraper Scraper
}

// NewShortlinkService constructs a new ShortlinkService with dependencies injected.
func NewShortlinkService(db *database.RedirectDB, scraper Scraper) *ShortlinkService {
	return &ShortlinkService{
		db:      db,
		scraper: scraper,
	}
}

// CreateShortlink persists a new shortlink with scraped web metadata context.
func (s *ShortlinkService) CreateShortlink(ctx context.Context, key, url, domain string) error {
	summary := s.scraper.Scrape(ctx, url)
	if err := s.db.CreateShortlink(key, url, domain, summary); err != nil {
		return fmt.Errorf("creating shortlink: %w", err)
	}
	return nil
}

// ListShortlinks returns all shortlinks.
func (s *ShortlinkService) ListShortlinks() ([]structs.Shortlink, error) {
	shortlinks, err := s.db.ListShortlinks()
	if err != nil {
		return nil, fmt.Errorf("listing shortlinks: %w", err)
	}
	return shortlinks, nil
}

// DeleteShortlink removes a shortlink by its key.
func (s *ShortlinkService) DeleteShortlink(key string) error {
	if err := s.db.DeleteShortlink(key); err != nil {
		return fmt.Errorf("deleting shortlink with key %q: %w", key, err)
	}
	return nil
}
