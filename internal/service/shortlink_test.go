package service

import (
	"testing"

	"github.com/beecorrea/shortlinks/internal/database"
)

func TestShortlinkService(t *testing.T) {
	// Initialize in-memory SQLite database
	db, err := database.NewRedirectDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	defer db.Close()

	svc := NewShortlinkService(db)

	// Test CreateShortlink
	err = svc.CreateShortlink("g", "https://google.com", "local")
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
