package cmdui

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScanRepoEndpoints_MethodValidation(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/scan/repos", nil)
	w := httptest.NewRecorder()

	handleAPIScanRepos(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed for POST to /api/scan/repos, got %d", w.Code)
	}
}

func TestScanRepoEndpoints_DisableRequiresPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/scan/disable", nil)
	w := httptest.NewRecorder()

	handleAPIScanDisable(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed for GET to /api/scan/disable, got %d", w.Code)
	}
}

func TestScanRepoEndpoints_ToggleInvalidPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/scan/disable", bytes.NewBufferString(`{invalid`))
	w := httptest.NewRecorder()

	handleAPIScanDisable(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with failed envelope for invalid JSON, got %d", w.Code)
	}
}

func TestRegisterAPIScanRoutes(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/scan/disable", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 from registered mux handler, got %d", w.Code)
	}
}
