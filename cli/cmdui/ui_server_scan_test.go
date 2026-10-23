package cmdui

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScanRepoEndpoints_MethodValidation(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/scan/repos", nil)
	w := httptest.NewRecorder()

	HandleAPIScanRepos(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed for POST to /api/scan/repos, got %d", w.Code)
	}
}

func TestScanRepoEndpoints_DisableRequiresPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/scan/repos/disable", nil)
	w := httptest.NewRecorder()

	HandleAPIScanReposDisable(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed for GET to /api/scan/repos/disable, got %d", w.Code)
	}
}

func TestScanRepoEndpoints_ToggleInvalidPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/scan/repos/toggle", bytes.NewBufferString(`{}`))
	w := httptest.NewRecorder()

	HandleAPIScanReposToggle(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for empty toggle payload, got %d", w.Code)
	}
}

func TestRegisterAPIScanRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterAPIScanRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/scan/repos/disable", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 from registered mux handler, got %d", w.Code)
	}
}
