package cmdmail

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMailUISPA_OfflineNoCDNDependencies(t *testing.T) {
	html := BuildMailComposerHTML()
	if len(html) == 0 {
		t.Fatal("expected non-empty HTML")
	}

	forbiddenPatterns := []string{
		"http://",
		"https://",
		"unpkg.com",
		"cdnjs.cloudflare.com",
		"cdn.jsdelivr.net",
		"fonts.googleapis.com",
	}

	for _, pattern := range forbiddenPatterns {
		if strings.Contains(strings.ToLower(html), pattern) {
			t.Errorf("found external CDN pattern %q in mail UI HTML", pattern)
		}
	}
}

func TestMailUISPA_RouteDelivery(t *testing.T) {
	mux := BuildMailUIMux()
	routes := []string{"/", "/ui"}

	for _, route := range routes {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK for %s, got %d", route, rec.Code)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "GitMap Mail Composer") {
			t.Errorf("expected page title in route %s", route)
		}
	}
}

func TestMailUISPA_NotFoundRoute(t *testing.T) {
	mux := BuildMailUIMux()
	req := httptest.NewRequest(http.MethodGet, "/random-unknown-page", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown route, got %d", rec.Code)
	}
}

func TestMailUI_TemplatesEndpoint(t *testing.T) {
	mux := BuildMailUIMux()
	req := httptest.NewRequest(http.MethodGet, "/api/mail/templates", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var templates []MailTemplate
	if err := json.Unmarshal(rec.Body.Bytes(), &templates); err != nil {
		t.Fatalf("failed to decode templates: %v", err)
	}

	if len(templates) == 0 {
		t.Error("expected non-empty templates list")
	}
}

func TestMailUI_AccountsEndpoint(t *testing.T) {
	mux := BuildMailUIMux()
	req := httptest.NewRequest(http.MethodGet, "/api/mail/accounts", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var accounts []MailAccount
	if err := json.Unmarshal(rec.Body.Bytes(), &accounts); err != nil {
		t.Fatalf("failed to decode accounts: %v", err)
	}

	if len(accounts) == 0 {
		t.Error("expected non-empty accounts list")
	}
}

func TestMailUI_SendEndpoint_Success(t *testing.T) {
	mux := BuildMailUIMux()
	payload := MailSendRequest{
		Account:      "default",
		To:           []string{"team@example.com"},
		Subject:      "Weekly Sync",
		BodyMarkdown: "## Summary",
		IsHTML:       true,
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/mail/send", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var resp MailSendResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !resp.IsSuccess {
		t.Errorf("expected IsSuccess true, got false")
	}
	if resp.RecipientCount != 1 {
		t.Errorf("expected RecipientCount 1, got %d", resp.RecipientCount)
	}
}

func TestMailUI_SendEndpoint_MissingRecipients(t *testing.T) {
	mux := BuildMailUIMux()
	payload := MailSendRequest{
		Account: "default",
		To:      []string{},
		Subject: "Empty",
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/mail/send", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestMailUI_PortProbing(t *testing.T) {
	listener, port, err := ScanAvailableMailListener("127.0.0.1", DefaultMailUIPort)
	if err != nil {
		t.Fatalf("port probing failed: %v", err)
	}
	defer listener.Close()

	if port < DefaultMailUIPort {
		t.Errorf("expected port >= %d, got %d", DefaultMailUIPort, port)
	}
}
