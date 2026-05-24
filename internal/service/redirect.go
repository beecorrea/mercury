package service

import (
	"fmt"
	"strings"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
)

type RedirectService struct {
	db     *database.RedirectDB
	domain string
}

func NewRedirectService(db *database.RedirectDB, domain string) *RedirectService {
	return &RedirectService{db: db, domain: domain}
}

// GetRedirectKey extracts the shortlink key from a request.
// It only handles requests whose hostname matches the configured domain or a subdomain of it.
// API routes (/api/*) are always bypassed.
func (s *RedirectService) GetRedirectKey(host, path string) string {
	if host != s.domain && !strings.HasSuffix(host, "."+s.domain) {
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
