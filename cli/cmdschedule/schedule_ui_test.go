package cmdschedule

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScheduleUISPA_OfflineNoCDNDependencies(t *testing.T) {
	html := BuildScheduleDashboardHTML()
	if len(html) == 0 {
		t.Fatal("expected non-empty schedule HTML")
	}

	forbidden := []string{
		"http://",
		"https://",
		"unpkg.com",
		"cdnjs.cloudflare.com",
		"cdn.jsdelivr.net",
		"fonts.googleapis.com",
	}

	for _, pattern := range forbidden {
		if strings.Contains(strings.ToLower(html), pattern) {
			t.Errorf("found external CDN pattern %q in schedule HTML", pattern)
		}
	}
}

func TestScheduleUISPA_RouteDelivery(t *testing.T) {
	mux := BuildScheduleUIMux()
	routes := []string{"/", "/ui"}

	for _, route := range routes {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK for %s, got %d", route, rec.Code)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "GitMap Schedule & Notifier Settings") {
			t.Errorf("expected header title in route %s", route)
		}
	}
}

func TestScheduleUI_TasksEndpoint(t *testing.T) {
	mux := BuildScheduleUIMux()
	req := httptest.NewRequest(http.MethodGet, "/api/schedule/tasks", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var tasks []ScheduleTaskItem
	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("failed to decode tasks: %v", err)
	}

	if len(tasks) == 0 {
		t.Error("expected non-empty task list")
	}
}

func TestScheduleUI_ToggleEndpoint(t *testing.T) {
	mux := BuildScheduleUIMux()
	payload := TaskToggleRequest{
		Name:      "daily-fleet-sync",
		IsEnabled: false,
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/schedule/toggle", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var resp TaskToggleResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "" || resp.IsEnabled {
		t.Errorf("expected success with IsEnabled false after toggle, got error=%s isEnabled=%v", resp.Error, resp.IsEnabled)
	}
}

func TestScheduleUI_WatchersEndpoints(t *testing.T) {
	mux := BuildScheduleUIMux()

	// Test GET watchers
	reqGet := httptest.NewRequest(http.MethodGet, "/api/schedule/watchers", nil)
	recGet := httptest.NewRecorder()
	mux.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on GET watchers, got %d", recGet.Code)
	}

	// Test POST valid watcher
	payload := WatcherCreateRequest{
		Path:            "./test",
		Events:          []string{"modify", "delete"},
		DebounceSeconds: 5,
		ActionCommand:   "gitmap pe -t",
		IsEnabled:       true,
	}
	bodyBytes, _ := json.Marshal(payload)
	reqPost := httptest.NewRequest(http.MethodPost, "/api/schedule/watchers", bytes.NewReader(bodyBytes))
	recPost := httptest.NewRecorder()
	mux.ServeHTTP(recPost, reqPost)

	if recPost.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on POST watcher, got %d", recPost.Code)
	}

	var resp WatcherCreateResponse
	if err := json.Unmarshal(recPost.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode watcher response: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("expected watcher response without error, got %s", resp.Error)
	}
}

func TestScheduleUI_WatcherEmptyPathValidation(t *testing.T) {
	mux := BuildScheduleUIMux()
	payload := WatcherCreateRequest{
		Path: "",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/schedule/watchers", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestScheduleUI_NotifiersEndpoint(t *testing.T) {
	mux := BuildScheduleUIMux()
	req := httptest.NewRequest(http.MethodGet, "/api/schedule/notifiers", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var notifiers []NotifierRecipientItem
	if err := json.Unmarshal(rec.Body.Bytes(), &notifiers); err != nil {
		t.Fatalf("failed to decode notifiers: %v", err)
	}

	if len(notifiers) == 0 {
		t.Error("expected non-empty notifiers list")
	}
}

func TestScheduleUI_PortProbing(t *testing.T) {
	listener, port, err := ScanAvailableScheduleListener("127.0.0.1", DefaultScheduleUIPort)
	if err != nil {
		t.Fatalf("port probing failed: %v", err)
	}
	defer listener.Close()

	if port < DefaultScheduleUIPort {
		t.Errorf("expected port >= %d, got %d", DefaultScheduleUIPort, port)
	}
}
