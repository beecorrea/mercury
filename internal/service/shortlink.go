package service

import (
	"fmt"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
)

type ShortlinkService struct {
	db *database.RedirectDB
}

func NewShortlinkService(db *database.RedirectDB) *ShortlinkService {
	return &ShortlinkService{db: db}
}

func (s *ShortlinkService) CreateShortlink(key, url, domain string) error {
	if err := s.db.CreateShortlink(key, url, domain); err != nil {
		return fmt.Errorf("creating shortlink: %w", err)
	}
	return nil
}

func (s *ShortlinkService) ListShortlinks() ([]structs.Shortlink, error) {
	shortlinks, err := s.db.ListShortlinks()
	if err != nil {
		return nil, fmt.Errorf("listing shortlinks: %w", err)
	}
	return shortlinks, nil
}

func (s *ShortlinkService) DeleteShortlink(key string) error {
	if err := s.db.DeleteShortlink(key); err != nil {
		return fmt.Errorf("deleting shortlink with key %q: %w", key, err)
	}
	return nil
}
