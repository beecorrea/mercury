package service

import (
	"testing"

	"github.com/beecorrea/shortlinks/internal/database"
)

func TestRedirectService(t *testing.T) {
	// Initialize in-memory SQLite database
	db, err := database.NewRedirectDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	defer db.Close()

	// Seed some test data directly in DB
	err = db.CreateShortlink("g", "https://google.com", "local")
	if err != nil {
		t.Fatalf("failed to create shortlink: %v", err)
	}

	svc := NewRedirectService(db)

	// Test GetShortlinkByKey
	s, err := svc.GetShortlinkByKey("g")
	if err != nil {
		t.Fatalf("failed to get shortlink: %v", err)
	}
	if s == nil {
		t.Fatal("expected shortlink to be found, got nil")
	}
	if s.URL != "https://google.com" || s.Domain != "local" {
		t.Errorf("unexpected shortlink: %+v", s)
	}

	// Test Get missing key
	s, err = svc.GetShortlinkByKey("missing")
	if err != nil {
		t.Fatalf("failed to get missing key: %v", err)
	}
	if s != nil {
		t.Errorf("expected nil for missing key, got: %+v", s)
	}
}
