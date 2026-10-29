package cmdui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestIDESetupStatus_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ide/setup/status", nil)
	rec := httptest.NewRecorder()

	handleIDESetupStatus(rec, req)
	assertStatusSuccess(t, rec)
}

func assertStatusSuccess(t *testing.T, rec *httptest.ResponseRecorder) {
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp IDESetupStatusResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.IsSuccess || len(resp.DetectedEditors) != 3 {
		t.Fatalf("expected success with 3 editors, got: %+v", resp)
	}
}

func TestIDESetupStatus_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/ide/setup/status", nil)
	rec := httptest.NewRecorder()

	handleIDESetupStatus(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestIDESetupPresets_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ide/setup/presets", nil)
	rec := httptest.NewRecorder()

	handleIDESetupPresets(rec, req)
	assertPresetsSuccess(t, rec)
}

func assertPresetsSuccess(t *testing.T, rec *httptest.ResponseRecorder) {
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp IDESetupPresetsResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.IsSuccess || len(resp.Presets) < 6 {
		t.Fatalf("expected >= 6 presets, got: %+v", resp)
	}
}

func TestIDESetupPresets_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/ide/setup/presets", nil)
	rec := httptest.NewRecorder()

	handleIDESetupPresets(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestIDESetupTestConnection_Success(t *testing.T) {
	server := createMockProbeServer(http.StatusOK)
	defer server.Close()

	payload := `{"provider":"kimi","baseUrl":"` + server.URL + `","apiKey":"sk-test","model":"kimi-k3.1-pro"}`
	rec := postTestConnection(payload)

	var resp IDETestConnectionResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.IsSuccess || resp.HTTPStatusCode != http.StatusOK {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func createMockProbeServer(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`{"data": []}`))
	}))
}

func postTestConnection(payload string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/ide/setup/test-connection", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()
	handleIDESetupTestConnection(rec, req)

	return rec
}

func TestIDESetupTestConnection_HttpError(t *testing.T) {
	server := createMockProbeServer(http.StatusUnauthorized)
	defer server.Close()

	payload := `{"provider":"kimi","baseUrl":"` + server.URL + `","apiKey":"bad-key"}`
	rec := postTestConnection(payload)

	var resp IDETestConnectionResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.IsSuccess || resp.HTTPStatusCode != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized failure, got %+v", resp)
	}
}

func TestIDESetupTestConnection_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/ide/setup/test-connection", bytes.NewBufferString("{bad"))
	rec := httptest.NewRecorder()

	handleIDESetupTestConnection(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestIDESetupTestConnection_MissingBaseURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/ide/setup/test-connection", bytes.NewBufferString(`{"provider":"unknown"}`))
	rec := httptest.NewRecorder()

	handleIDESetupTestConnection(rec, req)

	var resp IDETestConnectionResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.IsSuccess {
		t.Fatalf("expected validation failure for missing baseUrl, got success")
	}
}

func TestIDESetupTestConnection_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ide/setup/test-connection", nil)
	rec := httptest.NewRecorder()

	handleIDESetupTestConnection(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestIDESetupApply_DryRun(t *testing.T) {
	payload := `{"target":"all","provider":"kimi","baseUrl":"https://api.moonshot.cn/v1","apiKey":"sk-test","model":"kimi-k3.1-pro","isDryRun":true}`
	rec := postApply(payload)

	var resp IDESetupApplyResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.IsSuccess || !resp.IsDryRun || len(resp.ConfiguredTargets) != 3 {
		t.Fatalf("expected dry run success with 3 targets, got %+v", resp)
	}
}

func TestIDESetupApply_LiveWithBackup(t *testing.T) {
	setupTestConfigOverrides(t)
	defer clearTargetConfigOverrides()

	payload := `{"target":"all","provider":"kimi","baseUrl":"https://api.moonshot.cn/v1","apiKey":"sk-test","model":"kimi-k3.1-pro","isDryRun":false,"isBackupEnabled":true}`
	rec := postApply(payload)

	var resp IDESetupApplyResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.IsSuccess || len(resp.ConfiguredTargets) != 3 {
		t.Fatalf("expected live apply success with 3 targets, got %+v", resp)
	}

	verifyApplyBackupFlow(t, payload)
}

func setupTestConfigOverrides(t *testing.T) {
	tmp := t.TempDir()
	setTargetConfigOverride("cline", filepath.Join(tmp, "cline.json"))
	setTargetConfigOverride("continue", filepath.Join(tmp, "continue.json"))
	setTargetConfigOverride("zed", filepath.Join(tmp, "zed.json"))
}

func postApply(payload string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/ide/setup/apply", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()
	handleIDESetupApply(rec, req)

	return rec
}

func verifyApplyBackupFlow(t *testing.T, payload string) {
	rec := postApply(payload)

	var resp IDESetupApplyResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.IsSuccess || len(resp.BackupPaths) == 0 {
		t.Fatalf("expected backup paths created on second apply, got %+v", resp)
	}
}

func TestIDESetupApply_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/ide/setup/apply", bytes.NewBufferString("{bad"))
	rec := httptest.NewRecorder()

	handleIDESetupApply(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestIDESetupApply_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ide/setup/apply", nil)
	rec := httptest.NewRecorder()

	handleIDESetupApply(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestMountIDESetupRoutes(t *testing.T) {
	mux := http.NewServeMux()
	mountAPIRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/ide/setup/presets", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from mounted route, got %d", rec.Code)
	}
}
