package service

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/beecorrea/shortlinks/internal/database"
)

func TestScrapeRetrier(t *testing.T) {
	db, err := database.NewRedirectDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	defer db.Close()

	// Seed a failed link.
	if err := db.CreateShortlink("bad", "https://bad.example.com", "local", "Could not scrape bad.example.com: connection refused"); err != nil {
		t.Fatalf("failed to seed shortlink: %v", err)
	}

	t.Run("success path: updates summary on successful re-scrape", func(t *testing.T) {
		scraper := &mockScraper{summary: "Recovered description"}
		retrier := NewScrapeRetrier(db, scraper, time.Minute, 5)

		retrier.retryOnce(context.Background())

		links, err := db.ListShortlinks()
		if err != nil {
			t.Fatalf("listing shortlinks: %v", err)
		}
		var found bool
		for _, l := range links {
			if l.Key == "bad" {
				found = true
				if l.Summary != "Recovered description" {
					t.Errorf("expected updated summary, got %q", l.Summary)
				}
				if l.ScrapeAttempts != 1 {
					t.Errorf("expected scrape_attempts=1, got %d", l.ScrapeAttempts)
				}
			}
		}
		if !found {
			t.Error("link 'bad' not found after retry")
		}
	})

	// Reset the link back to a failed state for the next sub-test.
	if err := db.UpdateSummary("bad", "Could not scrape bad.example.com: connection refused"); err != nil {
		t.Fatalf("resetting summary: %v", err)
	}

	t.Run("failure path: increments attempts but keeps failed summary", func(t *testing.T) {
		scraper := &mockScraper{err: errors.New("still unreachable")}
		retrier := NewScrapeRetrier(db, scraper, time.Minute, 5)

		retrier.retryOnce(context.Background())

		links, err := db.ListShortlinks()
		if err != nil {
			t.Fatalf("listing shortlinks: %v", err)
		}
		for _, l := range links {
			if l.Key == "bad" {
				if l.ScrapeAttempts < 2 {
					t.Errorf("expected scrape_attempts >= 2, got %d", l.ScrapeAttempts)
				}
				if l.Summary != "Could not scrape bad.example.com: connection refused" {
					t.Errorf("expected summary unchanged, got %q", l.Summary)
				}
			}
		}
	})

	t.Run("cap enforcement: exhausted link is excluded from retry batch", func(t *testing.T) {
		var scrapeCount atomic.Int32
		scraper := &countingScraper{count: &scrapeCount, summary: "Recovered"}
		// maxAttempts=2 — the link already has >= 2 attempts from previous sub-tests.
		retrier := NewScrapeRetrier(db, scraper, time.Minute, 2)

		retrier.retryOnce(context.Background())

		if n := scrapeCount.Load(); n != 0 {
			t.Errorf("expected no scrape calls for exhausted link, got %d", n)
		}
	})

	t.Run("Start: exits cleanly on context cancellation", func(t *testing.T) {
		scraper := &mockScraper{}
		retrier := NewScrapeRetrier(db, scraper, time.Hour, 5)

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() {
			retrier.Start(ctx)
			close(done)
		}()

		cancel()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Error("Start did not exit after context cancellation")
		}
	})
}

// countingScraper records how many times Scrape is called.
type countingScraper struct {
	count   *atomic.Int32
	summary string
}

func (c *countingScraper) Scrape(_ context.Context, _ string) (string, error) {
	c.count.Add(1)
	return c.summary, nil
}
