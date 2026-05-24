package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestScrapeSummary(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedResult string
	}{
		{
			name: "meta name description",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<html><head><meta name="description" content="This is description!"></head></html>`))
			},
			expectedResult: "This is description!",
		},
		{
			name: "meta og description",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<html><head><meta property="og:description" content="This is og description!"></head></html>`))
			},
			expectedResult: "This is og description!",
		},
		{
			name: "fallback to title",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<html><head><title>This is title!</title></head></html>`))
			},
			expectedResult: "This is title!",
		},
		{
			name: "fallback to redirect message",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<html><head></head></html>`))
			},
			expectedResult: "Shortlink redirect to ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			got := scrapeSummary(server.URL)
			if strings.HasPrefix(tc.expectedResult, "Shortlink redirect to ") {
				if !strings.HasPrefix(got, "Shortlink redirect to ") {
					t.Errorf("expected summary to start with 'Shortlink redirect to ', got %q", got)
				}
			} else {
				if got != tc.expectedResult {
					t.Errorf("expected summary %q, got %q", tc.expectedResult, got)
				}
			}
		})
	}
}
