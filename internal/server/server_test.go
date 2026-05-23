package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beecorrea/shortlinks/internal/structs"
)

func TestServer(t *testing.T) {
	// Initialize Server with in-memory SQLite database
	srv, err := NewServer(":memory:")
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer srv.Close()

	// 1. Test POST /api/shorten
	reqBody := `{"key":"testkey", "url":"https://example.com", "domain":"local"}`
	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.App.Test(req)
	if err != nil {
		t.Fatalf("failed to run request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("expected POST /api/shorten to return 201 Created, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// 2. Test GET /api/links
	req = httptest.NewRequest("GET", "/api/links", nil)
	resp, err = srv.App.Test(req)
	if err != nil {
		t.Fatalf("failed to run request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected GET /api/links to return 200 OK, got %d", resp.StatusCode)
	}
	var shortlinks []structs.Shortlink
	if err := json.NewDecoder(resp.Body).Decode(&shortlinks); err != nil {
		t.Fatalf("failed to decode GET /api/links response: %v", err)
	}
	if len(shortlinks) != 1 || shortlinks[0].Key != "testkey" || shortlinks[0].URL != "https://example.com" {
		t.Errorf("unexpected shortlinks list: %+v", shortlinks)
	}

	// 3. Test Host-based Redirect for existing key
	req = httptest.NewRequest("GET", "/testkey", nil)
	req.Host = "mercury.local"
	resp, err = srv.App.Test(req)
	if err != nil {
		t.Fatalf("failed to run request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("expected redirect to return 302 Found, got %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if loc != "https://example.com" {
		t.Errorf("expected redirect location https://example.com, got %s", loc)
	}

	// 4. Test Host-based Redirect for missing key
	req = httptest.NewRequest("GET", "/missingkey", nil)
	req.Host = "mercury.local"
	resp, err = srv.App.Test(req)
	if err != nil {
		t.Fatalf("failed to run request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected missing key redirect to return 404 Not Found, got %d", resp.StatusCode)
	}

	// 5. Test Non-Mercury Host serves Dashboard
	req = httptest.NewRequest("GET", "/", nil)
	req.Host = "localhost:45800"
	resp, err = srv.App.Test(req)
	if err != nil {
		t.Fatalf("failed to run request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected main page on localhost to return 200 OK, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("Mercury")) {
		t.Error("expected served dashboard to contain 'Mercury'")
	}
}
