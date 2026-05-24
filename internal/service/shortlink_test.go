package service

import (
	"context"
	"testing"

	"github.com/beecorrea/shortlinks/internal/database"
)

type mockScraper struct {
	summary string
}

func (m *mockScraper) Scrape(ctx context.Context, targetURL string) string {
	return m.summary
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

	// Test CreateShortlink
	err = svc.CreateShortlink(context.Background(), "g", "https://google.com", "local")
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
	if len(links) != 0 {
		t.Errorf("expected 0 shortlinks, got %d", len(links))
	}
}
