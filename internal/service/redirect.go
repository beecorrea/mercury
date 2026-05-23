package service

import (
	"fmt"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
)

type RedirectService struct {
	db *database.RedirectDB
}

func NewRedirectService(db *database.RedirectDB) *RedirectService {
	return &RedirectService{db: db}
}

func (s *RedirectService) GetShortlinkByKey(key string) (*structs.Shortlink, error) {
	shortlink, err := s.db.GetShortlinkByKey(key)
	if err != nil {
		return nil, fmt.Errorf("getting shortlink by key %q: %w", key, err)
	}
	return shortlink, nil
}
