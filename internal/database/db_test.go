package database

import (
	"testing"
)

func TestDBOperations(t *testing.T) {
	// Initialize in-memory SQLite database for testing
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	defer db.Close()

	// Test CreateLink
	err = db.CreateLink("g", "https://google.com", "local")
	if err != nil {
		t.Fatalf("failed to create link: %v", err)
	}

	// Test GetLinkByKey
	link, err := db.GetLinkByKey("g")
	if err != nil {
		t.Fatalf("failed to get link: %v", err)
	}
	if link == nil {
		t.Fatal("expected link to be found, got nil")
	}
	if link.URL != "https://google.com" || link.Domain != "local" {
		t.Errorf("expected url=https://google.com, domain=local; got url=%s, domain=%s", link.URL, link.Domain)
	}
	if link.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be non-zero time")
	}
	t.Logf("CreatedAt value: %v", link.CreatedAt)

	// Test Duplicate CreateLink (should fail due to UNIQUE constraint)
	err = db.CreateLink("g", "https://github.com", "local")
	if err == nil {
		t.Error("expected error creating duplicate key, got nil")
	}

	// Test ListLinks
	links, err := db.ListLinks()
	if err != nil {
		t.Fatalf("failed to list links: %v", err)
	}
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}

	// Test DeleteLink
	err = db.DeleteLink("g")
	if err != nil {
		t.Fatalf("failed to delete link: %v", err)
	}

	// Test GetLinkByKey after delete
	link, err = db.GetLinkByKey("g")
	if err != nil {
		t.Fatalf("failed to get link after delete: %v", err)
	}
	if link != nil {
		t.Error("expected link to be deleted, but it was found")
	}
}
