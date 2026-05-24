package scraper

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	userAgentHeader = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	acceptHeader    = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
)

var (
	metaRegex    = regexp.MustCompile(`(?i)<meta\s+([^>]+)>`)
	nameRegex    = regexp.MustCompile(`(?i)\b(?:name|property)\s*=\s*["'](?:og:)?description["']`)
	contentRegex = regexp.MustCompile(`(?i)\bcontent\s*=\s*["']([^"']*)["']`)
	titleRegex   = regexp.MustCompile(`(?i)<title\b[^>]*>(.*?)</title>`)
)

// HTTPScraper implements the service.Scraper interface using standard HTTP GET requests.
type HTTPScraper struct {
	client *http.Client
}

// NewHTTPScraper returns a new instance of HTTPScraper.
func NewHTTPScraper() *HTTPScraper {
	return &HTTPScraper{
		client: &http.Client{},
	}
}

// Scrape extracts metadata from the target URL and returns a summary description.
func (s *HTTPScraper) Scrape(ctx context.Context, targetURL string) (string, error) {
	body, err := s.fetchHTML(ctx, targetURL)
	if err != nil {
		return "", err
	}

	summary := extractMetadata(body)
	if summary == "" {
		return "", fmt.Errorf("no description or title metadata found")
	}

	return summary, nil
}

// fetchHTML fetches the HTML body from the target URL, applying timeout and size constraints.
func (s *HTTPScraper) fetchHTML(ctx context.Context, targetURL string) (string, error) {
	scrapeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := s.buildRequest(scrapeCtx, targetURL)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("site returned status %d", resp.StatusCode)
	}

	body, err := readBoundedBody(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// buildRequest constructs an http.Request with custom headers.
func (s *HTTPScraper) buildRequest(ctx context.Context, targetURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgentHeader)
	req.Header.Set("Accept", acceptHeader)

	return req, nil
}

// readBoundedBody reads up to 100KB of response body.
func readBoundedBody(rc io.ReadCloser) (string, error) {
	limitedReader := io.LimitReader(rc, 100*1024)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

// extractMetadata orchestrates metadata extraction from the HTML body string.
func extractMetadata(body string) string {
	if desc := extractMetaDescription(body); desc != "" {
		return desc
	}
	if title := extractTitle(body); title != "" {
		return title
	}
	return ""
}

// extractMetaDescription searches for description or og:description in meta tags.
func extractMetaDescription(body string) string {
	metaMatches := metaRegex.FindAllStringSubmatch(body, -1)
	for _, match := range metaMatches {
		attributes := match[1]
		if !nameRegex.MatchString(attributes) {
			continue
		}
		contentMatch := contentRegex.FindStringSubmatch(attributes)
		if len(contentMatch) <= 1 {
			continue
		}
		description := html.UnescapeString(strings.TrimSpace(contentMatch[1]))
		if description != "" {
			return description
		}
	}
	return ""
}

// extractTitle searches for the content of the <title> tag.
func extractTitle(body string) string {
	titleMatch := titleRegex.FindStringSubmatch(body)
	if len(titleMatch) <= 1 {
		return ""
	}
	return html.UnescapeString(strings.TrimSpace(titleMatch[1]))
}
