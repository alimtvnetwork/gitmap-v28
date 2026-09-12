package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTerminalCommandExecAPI(t *testing.T) {
	mux := http.NewServeMux()
	mountTerminalHandlers(mux)

	reqBody := []byte(`{"command": "echo API_TEST_SUCCESS"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/command/exec", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	if !bytes.Contains(rec.Body.Bytes(), []byte("API_TEST_SUCCESS")) {
		t.Errorf("expected response to contain output, got %s", rec.Body.String())
	}
}
