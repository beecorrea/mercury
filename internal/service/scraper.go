package service

import "context"

// Scraper defines the contract for extracting metadata summaries from target URLs.
type Scraper interface {
	// Scrape extracts a description/summary or title from the target URL.
	// It accepts a context for cancellation and timeout propagation.
	// It returns the extracted string and an error if the scrape failed.
	Scrape(ctx context.Context, targetURL string) (string, error)
}
