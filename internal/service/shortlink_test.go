package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
)

type mockScraper struct {
	summary string
	err     error
}

func (m *mockScraper) Scrape(ctx context.Context, targetURL string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.summary, nil
}

func TestShortlinkService(t *testing.T) {
	// Initialize in-memory SQLite database
	db, err := database.NewRedirectDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	defer db.Close()

	mockScr := &mockScraper{summary: "Mock Description"}
	svc := NewShortlinkService(db, mockScr)

	// Test CreateShortlink (success path)
	err = svc.CreateShortlink(context.Background(), structs.Shortlink{
		Key:    "g",
		URL:    "https://google.com",
		Domain: "local",
	})
	if err != nil {
		t.Fatalf("failed to create shortlink: %v", err)
	}

	// Test ListShortlinks
	links, err := svc.ListShortlinks()
	if err != nil {
		t.Fatalf("failed to list shortlinks: %v", err)
	}
	if len(links) != 1 || links[0].Key != "g" {
		t.Errorf("unexpected shortlinks list: %+v", links)
	}
	if links[0].Summary != "Mock Description" {
		t.Errorf("expected summary %q, got %q", "Mock Description", links[0].Summary)
	}

	// Test CreateShortlink (fallback path on scraper error, e.g. 404)
	mockScrErr := &mockScraper{err: fmt.Errorf("site returned status 404")}
	svcErr := NewShortlinkService(db, mockScrErr)
	err = svcErr.CreateShortlink(context.Background(), structs.Shortlink{
		Key:    "brokenkey",
		URL:    "https://brokenlink.com/nonexistent",
		Domain: "local",
	})
	if err != nil {
		t.Fatalf("failed to create shortlink with broken link: %v", err)
	}

	links, err = svcErr.ListShortlinks()
	if err != nil {
		t.Fatalf("failed to list shortlinks: %v", err)
	}
	// We have 2 shortlinks now: "g" and "brokenkey"
	var brokenLink structs.Shortlink
	for _, l := range links {
		if l.Key == "brokenkey" {
			brokenLink = l
		}
	}
	expectedFallbackSummary := "Could not scrape brokenlink.com: site returned status 404"
	if brokenLink.Summary != expectedFallbackSummary {
		t.Errorf("expected fallback summary %q, got %q", expectedFallbackSummary, brokenLink.Summary)
	}

	// Test DeleteShortlink
	err = svc.DeleteShortlink("g")
	if err != nil {
		t.Fatalf("failed to delete shortlink: %v", err)
	}

	// Test ListShortlinks after delete
	links, err = svc.ListShortlinks()
	if err != nil {
		t.Fatalf("failed to list shortlinks after delete: %v", err)
	}
	if len(links) != 1 || links[0].Key != "brokenkey" {
		t.Errorf("expected only 'brokenkey' shortlink remaining, got %d links", len(links))
	}
}
