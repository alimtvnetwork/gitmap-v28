package cmdlogviewer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAPIPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/log/ping", nil)
	w := httptest.NewRecorder()

	handleAPIPing(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %q", body["status"])
	}
}

func TestHandleAPIDiffValidation(t *testing.T) {
	handler := handleAPIDiff(".")

	t.Run("missing fields returns bad request", func(t *testing.T) {
		payload := bytes.NewBufferString(`{"commit":""}`)
		req := httptest.NewRequest(http.MethodPost, "/api/log/diff", payload)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("malformed json returns bad request", func(t *testing.T) {
		payload := bytes.NewBufferString(`invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/api/log/diff", payload)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestHandleAPIBranchCompareValidation(t *testing.T) {
	handler := handleAPIBranchCompare(".")

	payload := bytes.NewBufferString(`{"branchA":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/log/branch-compare", payload)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleAPICommitDetailMissingHash(t *testing.T) {
	handler := handleAPICommitDetail(".")
	req := httptest.NewRequest(http.MethodGet, "/api/log/commit", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing hash, got %d", w.Code)
	}
}
