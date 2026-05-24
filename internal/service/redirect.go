package service

import (
	"fmt"
	"strings"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
)

type RedirectService struct {
	db *database.RedirectDB
}

func NewRedirectService(db *database.RedirectDB) *RedirectService {
	return &RedirectService{db: db}
}

func (s *RedirectService) GetRedirectKey(host, path string) string {
	if host != "mercury" && !strings.HasPrefix(host, "mercury.") {
		return ""
	}

	if strings.HasPrefix(path, "/api/") {
		return ""
	}

	return strings.Trim(path, "/")
}

func (s *RedirectService) GetShortlinkByKey(key string) (*structs.Shortlink, error) {
	shortlink, err := s.db.GetShortlinkByKey(key)
	if err != nil {
		return nil, fmt.Errorf("getting shortlink by key %q: %w", key, err)
	}
	return shortlink, nil
}
