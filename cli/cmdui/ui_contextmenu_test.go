package cmdui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIContextMenuStatus(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/context-menu/status", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp ContextMenuStatusResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true, got false")
	}
}

func TestAPIContextMenuCleanPowerShellMethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/context-menu/clean-powershell", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 MethodNotAllowed, got %d", rec.Code)
	}
}

func TestAPIContextMenuCleanPowerShellPost(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/context-menu/clean-powershell", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	assertCleanResponseSuccess(t, rec)
}

func assertCleanResponseSuccess(t *testing.T, rec *httptest.ResponseRecorder) {
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if success, isOK := resp["success"].(bool); !isOK || !success {
		t.Fatalf("expected success true, got %v", resp)
	}
}

func TestUIAssetsMarkupContextMenuElements(t *testing.T) {
	assertMarkupContains(t, IndexHTML, "Context Menu &amp; Shell Integrations")
	assertMarkupContains(t, IndexHTML, "Install Context Menu")
	assertMarkupContains(t, IndexHTML, "Clean Duplicate PowerShell Menu")
	assertMarkupContains(t, IndexHTML, "gitmap pa --ssh")
	assertMarkupContains(t, IndexHTML, "gitmap nodes ping")
	assertMarkupContains(t, IndexHTML, "gitmap show-log")
	assertMarkupContains(t, IndexHTML, "btn-fullscreen")
	assertMarkupContains(t, IndexHTML, "btn-exit")
}

func assertMarkupContains(t *testing.T, html, substr string) {
	hasSubstr := strings.Contains(html, substr)
	if !hasSubstr {
		t.Errorf("expected IndexHTML to contain %q", substr)
	}
}

func TestAPIContextMenuToggleGitActionsMethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/context-menu/toggle-git-actions", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 MethodNotAllowed, got %d", rec.Code)
	}
}

func TestAPIContextMenuToggleGitActionsPost(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/context-menu/toggle-git-actions", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
