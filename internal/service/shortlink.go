package service

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
)

type ShortlinkService struct {
	db *database.RedirectDB
}

func NewShortlinkService(db *database.RedirectDB) *ShortlinkService {
	return &ShortlinkService{db: db}
}

func scrapeSummary(targetURL string) string {
	// Fallback description
	parsedURL, err := url.Parse(targetURL)
	fallbackHost := "target URL"
	if err == nil && parsedURL.Host != "" {
		fallbackHost = parsedURL.Host
	}
	fallbackSummary := "Shortlink redirect to " + fallbackHost

	// Set a 3-second HTTP timeout using context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return fallbackSummary
	}

	// Use a standard User-Agent header to avoid bot blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fallbackSummary
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fallbackSummary
	}

	// Limit reader to 100KB to prevent memory overflow on large resources
	limitedReader := io.LimitReader(resp.Body, 100*1024)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return fallbackSummary
	}

	bodyStr := string(bodyBytes)

	// Parse the HTML body using robust regular expressions to extract <meta name="description"> or <meta property="og:description">
	metaRegex := regexp.MustCompile(`(?i)<meta\s+([^>]+)>`)
	nameRegex := regexp.MustCompile(`(?i)\b(?:name|property)\s*=\s*["'](?:og:)?description["']`)
	contentRegex := regexp.MustCompile(`(?i)\bcontent\s*=\s*["']([^"']*)["']`)

	var description string
	metaMatches := metaRegex.FindAllStringSubmatch(bodyStr, -1)
	for _, match := range metaMatches {
		attributes := match[1]
		if nameRegex.MatchString(attributes) {
			if contentMatch := contentRegex.FindStringSubmatch(attributes); len(contentMatch) > 1 {
				description = html.UnescapeString(strings.TrimSpace(contentMatch[1]))
				if description != "" {
					break
				}
			}
		}
	}

	if description != "" {
		return description
	}

	// Fall back to <title>
	titleRegex := regexp.MustCompile(`(?i)<title\b[^>]*>(.*?)</title>`)
	titleMatch := titleRegex.FindStringSubmatch(bodyStr)
	if len(titleMatch) > 1 {
		title := html.UnescapeString(strings.TrimSpace(titleMatch[1]))
		if title != "" {
			return title
		}
	}

	return fallbackSummary
}

func (s *ShortlinkService) CreateShortlink(key, url, domain string) error {
	summary := scrapeSummary(url)
	if err := s.db.CreateShortlink(key, url, domain, summary); err != nil {
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
