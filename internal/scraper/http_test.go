package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestExtractMetaDescription(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "meta name description double quotes",
			body:     `<html><head><meta name="description" content="Meta Description"></head></html>`,
			expected: "Meta Description",
		},
		{
			name:     "meta name description single quotes",
			body:     `<html><head><meta name='description' content='Meta Description'></head></html>`,
			expected: "Meta Description",
		},
		{
			name:     "meta property og:description",
			body:     `<html><head><meta property="og:description" content="OG Description"></head></html>`,
			expected: "OG Description",
		},
		{
			name:     "meta description case insensitive",
			body:     `<html><head><META NAME="DESCRIPTION" CONTENT="CASE INSENSITIVE"></head></html>`,
			expected: "CASE INSENSITIVE",
		},
		{
			name:     "no meta description",
			body:     `<html><head><title>No Meta</title></head></html>`,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractMetaDescription(tc.body)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "basic title",
			body:     `<html><head><title>Hello World</title></head></html>`,
			expected: "Hello World",
		},
		{
			name:     "title case insensitive",
			body:     `<html><head><TITLE>CASE INSENSITIVE</TITLE></head></html>`,
			expected: "CASE INSENSITIVE",
		},
		{
			name:     "title with attributes",
			body:     `<html><head><title class="header">Attributes</title></head></html>`,
			expected: "Attributes",
		},
		{
			name:     "no title",
			body:     `<html><head></head></html>`,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractTitle(tc.body)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestHTTPScraper_Scrape(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		expectedExact string
		expectErr     bool
		errSubstring  string
	}{
		{
			name: "success meta description",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<html><head><meta name="description" content="Super desc"></head></html>`))
			},
			expectedExact: "Super desc",
			expectErr:     false,
		},
		{
			name: "success title fallback",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<html><head><title>Only Title</title></head></html>`))
			},
			expectedExact: "Only Title",
			expectErr:     false,
		},
		{
			name: "site returns 404 error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectErr:    true,
			errSubstring: "site returned status 404",
		},
		{
			name: "site returns 500 error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectErr:    true,
			errSubstring: "site returned status 500",
		},
		{
			name: "empty response error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<html><head></head></html>`))
			},
			expectErr:    true,
			errSubstring: "no description or title metadata found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			scraper := NewHTTPScraper()
			got, err := scraper.Scrape(context.Background(), server.URL)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil error", tc.errSubstring)
				}
				if !strings.Contains(err.Error(), tc.errSubstring) {
					t.Errorf("expected error %q to contain %q", err.Error(), tc.errSubstring)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tc.expectedExact {
					t.Errorf("expected EXACT %q, got %q", tc.expectedExact, got)
				}
			}
		})
	}
}

func TestHTTPScraper_Scrape_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	scraper := NewHTTPScraper()
	_, err := scraper.Scrape(ctx, server.URL)

	if err == nil {
		t.Fatal("expected error on timeout, got nil")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") && !strings.Contains(err.Error(), "canceled") {
		t.Errorf("expected timeout or cancellation error, got: %v", err)
	}
}
