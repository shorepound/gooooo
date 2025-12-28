package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shorepound/gooooo/internal/store"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	RegisterRoutes(r, store.New())
	return r
}

func TestHealth(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestCRUDFlow(t *testing.T) {
	r := setupRouter()

	// Create
	payload := `{"name":"foo","description":"bar"}`
	req := httptest.NewRequest("POST", "/items", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", rr.Code)
	}
	var created store.Item
	if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
		t.Fatalf("create decode: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("create: expected non-zero id")
	}

	// List
	req = httptest.NewRequest("GET", "/items", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", rr.Code)
	}
	var items []store.Item
	if err := json.NewDecoder(rr.Body).Decode(&items); err != nil {
		t.Fatalf("list decode: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("list: expected 1 item, got %d", len(items))
	}

	// Get
	req = httptest.NewRequest("GET", fmt.Sprintf("/items/%d", created.ID), nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", rr.Code)
	}

	// Update
	upd := `{"name":"foo2","description":"bar2"}`
	req = httptest.NewRequest("PUT", fmt.Sprintf("/items/%d", created.ID), strings.NewReader(upd))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", rr.Code)
	}

	// Delete
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/items/%d", created.ID), nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", rr.Code)
	}

	// Get should be not found
	req = httptest.NewRequest("GET", fmt.Sprintf("/items/%d", created.ID), nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("get after delete: expected 404, got %d", rr.Code)
	}
}
