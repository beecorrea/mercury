package scraper

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
func (s *HTTPScraper) Scrape(ctx context.Context, targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	fallbackHost := "target URL"
	if err == nil && parsedURL.Host != "" {
		fallbackHost = parsedURL.Host
	}
	defaultFallback := "Shortlink redirect to " + fallbackHost

	// Limit network timeout to 3 seconds, nested under caller's context
	scrapeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(scrapeCtx, "GET", targetURL, nil)
	if err != nil {
		return defaultFallback
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return defaultFallback
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Could not scrape: Site returned status %d", resp.StatusCode)
	}

	// Read up to 100KB to bound memory usage
	limitedReader := io.LimitReader(resp.Body, 100*1024)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return defaultFallback
	}

	bodyStr := string(bodyBytes)
	if summary := extractMetadata(bodyStr); summary != "" {
		return summary
	}

	return defaultFallback
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

var (
	metaRegex    = regexp.MustCompile(`(?i)<meta\s+([^>]+)>`)
	nameRegex    = regexp.MustCompile(`(?i)\b(?:name|property)\s*=\s*["'](?:og:)?description["']`)
	contentRegex = regexp.MustCompile(`(?i)\bcontent\s*=\s*["']([^"']*)["']`)
	titleRegex   = regexp.MustCompile(`(?i)<title\b[^>]*>(.*?)</title>`)
)

// extractMetaDescription searches for description or og:description in meta tags.
func extractMetaDescription(body string) string {
	metaMatches := metaRegex.FindAllStringSubmatch(body, -1)
	for _, match := range metaMatches {
		attributes := match[1]
		if nameRegex.MatchString(attributes) {
			if contentMatch := contentRegex.FindStringSubmatch(attributes); len(contentMatch) > 1 {
				description := html.UnescapeString(strings.TrimSpace(contentMatch[1]))
				if description != "" {
					return description
				}
			}
		}
	}
	return ""
}

// extractTitle searches for the content of the <title> tag.
func extractTitle(body string) string {
	titleMatch := titleRegex.FindStringSubmatch(body)
	if len(titleMatch) > 1 {
		return html.UnescapeString(strings.TrimSpace(titleMatch[1]))
	}
	return ""
}
