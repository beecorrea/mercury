package database

import (
	"testing"
)

func TestDBOperations(t *testing.T) {
	// Initialize in-memory SQLite database for testing
	db, err := NewRedirectDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	defer db.Close()

	// Test CreateShortlink
	err = db.CreateShortlink("g", "https://google.com", "local")
	if err != nil {
		t.Fatalf("failed to create shortlink: %v", err)
	}

	// Test GetShortlinkByKey
	s, err := db.GetShortlinkByKey("g")
	if err != nil {
		t.Fatalf("failed to get shortlink: %v", err)
	}
	if s == nil {
		t.Fatal("expected shortlink to be found, got nil")
	}
	if s.URL != "https://google.com" || s.Domain != "local" {
		t.Errorf("expected url=https://google.com, domain=local; got url=%s, domain=%s", s.URL, s.Domain)
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be non-zero time")
	}
	t.Logf("CreatedAt value: %v", s.CreatedAt)

	// Test Duplicate CreateShortlink (should fail due to UNIQUE constraint)
	err = db.CreateShortlink("g", "https://github.com", "local")
	if err == nil {
		t.Error("expected error creating duplicate key, got nil")
	}

	// Test ListShortlinks
	shortlinks, err := db.ListShortlinks()
	if err != nil {
		t.Fatalf("failed to list shortlinks: %v", err)
	}
	if len(shortlinks) != 1 {
		t.Errorf("expected 1 shortlink, got %d", len(shortlinks))
	}

	// Test DeleteShortlink
	err = db.DeleteShortlink("g")
	if err != nil {
		t.Fatalf("failed to delete shortlink: %v", err)
	}

	// Test GetShortlinkByKey after delete
	s, err = db.GetShortlinkByKey("g")
	if err != nil {
		t.Fatalf("failed to get shortlink after delete: %v", err)
	}
	if s != nil {
		t.Error("expected shortlink to be deleted, but it was found")
	}
}
