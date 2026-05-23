package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beecorrea/shortlinks/internal/database"
	"github.com/beecorrea/shortlinks/internal/structs"
	"github.com/gofiber/fiber/v2"
)

func TestHandlers(t *testing.T) {
	// Initialize in-memory SQLite database
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory db: %v", err)
	}
	defer db.Close()

	// Initialize Fiber App
	app := fiber.New()
	app.Use(NewRedirectMiddleware(db))

	adminHandler := NewAdminHandler(db)
	app.Get("/", adminHandler.ServeDashboard)
	app.Get("/api/links", adminHandler.ListLinks)
	app.Post("/api/shorten", adminHandler.CreateLink)
	app.Delete("/api/links/:key", adminHandler.DeleteLink)

	// 1. Test POST /api/shorten
	reqBody := `{"key":"testkey", "url":"https://example.com", "domain":"local"}`
	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
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
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("failed to run request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected GET /api/links to return 200 OK, got %d", resp.StatusCode)
	}
	var links []structs.Link
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		t.Fatalf("failed to decode GET /api/links response: %v", err)
	}
	if len(links) != 1 || links[0].Key != "testkey" || links[0].URL != "https://example.com" {
		t.Errorf("unexpected links list: %+v", links)
	}

	// 3. Test Host-based Redirect for existing key
	req = httptest.NewRequest("GET", "/testkey", nil)
	req.Host = "mercury.local"
	resp, err = app.Test(req)
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
	resp, err = app.Test(req)
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
	resp, err = app.Test(req)
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
