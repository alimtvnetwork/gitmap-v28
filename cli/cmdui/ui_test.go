package cmdui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectSyntaxLanguage(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"config.json", "json"},
		{"readme.md", "markdown"},
		{"script.py", "python"},
		{"index.html", "html"},
		{"styles.css", "css"},
		{"run.sh", "shell"},
		{"unknown.xyz", "plaintext"},
	}

	for _, c := range cases {
		actual := detectSyntaxLanguage(c.path)
		if actual != c.expected {
			t.Errorf("path %s: expected %s, got %s", c.path, c.expected, actual)
		}
	}
}

func TestAPITerminalExec_LocalEcho(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	body := `{"command":"echo \"hello-gitmap\""}`
	req := httptest.NewRequest(http.MethodPost, "/api/terminal/exec", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp TerminalExecResp
	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if !resp.Success {
		t.Fatalf("expected success true, got false (error: %s, stderr: %s)", resp.Error, resp.Stderr)
	}

	if !strings.Contains(resp.Stdout, "hello-gitmap") {
		t.Fatalf("expected stdout to contain 'hello-gitmap', got %q", resp.Stdout)
	}
}

func TestAPISettings_AttributesPreservation(t *testing.T) {
	tempDir := t.TempDir()
	customSettingsPath = filepath.Join(tempDir, "ui_settings.json")
	defer func() {
		customSettingsPath = ""
	}()

	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	postBody := `{
		"theme": "dark",
		"clusterPort": 50000,
		"attributes": {
			"telegram.bot_token": "token-test-12345",
			"lap.default_hours": "48",
			"machine.name": "node-test-fleet",
			"special_repos.secrets_name": "secrets-vault"
		}
	}`
	postReq := httptest.NewRequest(http.MethodPost, "/api/settings", strings.NewReader(postBody))
	postReq.Header.Set("Content-Type", "application/json")
	postRec := httptest.NewRecorder()

	mux.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("POST expected status 200, got %d", postRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	getRec := httptest.NewRecorder()

	mux.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET expected status 200, got %d", getRec.Code)
	}

	var settings SettingsData
	err := json.Unmarshal(getRec.Body.Bytes(), &settings)

	if err != nil {
		t.Fatalf("failed to decode settings: %v", err)
	}

	if settings.Attributes["telegram.bot_token"] != "token-test-12345" {
		t.Errorf("expected telegram.bot_token preserved, got %q", settings.Attributes["telegram.bot_token"])
	}

	if settings.Attributes["lap.default_hours"] != "48" {
		t.Errorf("expected lap.default_hours preserved, got %q", settings.Attributes["lap.default_hours"])
	}

	if settings.Attributes["machine.name"] != "node-test-fleet" {
		t.Errorf("expected machine.name preserved, got %q", settings.Attributes["machine.name"])
	}

	if settings.Attributes["special_repos.secrets_name"] != "secrets-vault" {
		t.Errorf("expected special_repos.secrets_name preserved, got %q", settings.Attributes["special_repos.secrets_name"])
	}

	secondPostBody := `{
		"attributes": {
			"email.smtp_host": "smtp.example.com:587"
		}
	}`
	secondPostReq := httptest.NewRequest(http.MethodPost, "/api/settings", strings.NewReader(secondPostBody))
	secondPostReq.Header.Set("Content-Type", "application/json")
	secondPostRec := httptest.NewRecorder()

	mux.ServeHTTP(secondPostRec, secondPostReq)

	if secondPostRec.Code != http.StatusOK {
		t.Fatalf("second POST expected status 200, got %d", secondPostRec.Code)
	}

	getRec2 := httptest.NewRecorder()
	mux.ServeHTTP(getRec2, getReq)

	var settings2 SettingsData
	err = json.Unmarshal(getRec2.Body.Bytes(), &settings2)

	if err != nil {
		t.Fatalf("failed to decode merged settings: %v", err)
	}

	if settings2.Attributes["email.smtp_host"] != "smtp.example.com:587" {
		t.Errorf("expected email.smtp_host to be added, got %q", settings2.Attributes["email.smtp_host"])
	}

	if settings2.Attributes["telegram.bot_token"] != "token-test-12345" {
		t.Errorf("expected existing telegram.bot_token preserved after second POST, got %q", settings2.Attributes["telegram.bot_token"])
	}
}

