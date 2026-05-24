package service

import (
	"context"
	"log"
	"time"

	"github.com/beecorrea/shortlinks/internal/database"
)

// ScrapeRetrier periodically retries failed scrapes for shortlinks whose initial
// scrape did not succeed. It stops retrying a link once its attempt count reaches
// maxAttempts, preventing unbounded retries against permanently unreachable URLs.
//
// NOTE: This retrier is designed for single-replica deployments only. With multiple
// replicas sharing the same SQLite database, each instance runs its own retry loop
// independently, causing duplicate scrapes per tick and faster-than-expected
// exhaustion of the attempt cap. A PostgreSQL migration with SELECT … FOR UPDATE
// SKIP LOCKED is required to safely support horizontal scaling.
type ScrapeRetrier struct {
	db          *database.RedirectDB
	scraper     Scraper
	interval    time.Duration
	maxAttempts int
}

// NewScrapeRetrier constructs a ScrapeRetrier with the given dependencies and configuration.
// If interval is zero it defaults to 30 seconds.
func NewScrapeRetrier(db *database.RedirectDB, scraper Scraper, interval time.Duration, maxAttempts int) *ScrapeRetrier {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &ScrapeRetrier{
		db:          db,
		scraper:     scraper,
		interval:    interval,
		maxAttempts: maxAttempts,
	}
}

// Start runs the retry loop, waking every r.interval to process failed scrapes.
// It exits cleanly when ctx is cancelled.
func (r *ScrapeRetrier) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.retryOnce(ctx)
		}
	}
}

// retryOnce fetches all eligible failed links and attempts a re-scrape for each.
// The attempt counter is always incremented, whether the scrape succeeds or not,
// so a permanently unreachable link is eventually excluded from future cycles.
func (r *ScrapeRetrier) retryOnce(ctx context.Context) {
	links, err := r.db.ListFailedScrapes(r.maxAttempts)
	if err != nil {
		log.Printf("scrape retrier: listing failed scrapes: %v", err)
		return
	}

	for _, link := range links {
		if err := r.db.IncrementScrapeAttempts(link.Key); err != nil {
			log.Printf("scrape retrier: incrementing attempts for %q: %v", link.Key, err)
			continue
		}

		summary, err := r.scraper.Scrape(ctx, link.URL)
		if err != nil {
			log.Printf("scrape retrier: re-scrape failed for %q (attempt %d/%d): %v", link.Key, link.ScrapeAttempts+1, r.maxAttempts, err)
			continue
		}

		if err := r.db.UpdateSummary(link.Key, summary); err != nil {
			log.Printf("scrape retrier: updating summary for %q: %v", link.Key, err)
			continue
		}

		log.Printf("scrape retrier: successfully updated summary for %q", link.Key)
	}
}
