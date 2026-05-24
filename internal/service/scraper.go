package service

import "context"

// Scraper defines the contract for extracting metadata summaries from target URLs.
type Scraper interface {
	// Scrape extracts a description/summary or title from the target URL.
	// It accepts a context for cancellation and timeout propagation.
	// It must always return a fallback description, never an empty string.
	Scrape(ctx context.Context, targetURL string) string
}
