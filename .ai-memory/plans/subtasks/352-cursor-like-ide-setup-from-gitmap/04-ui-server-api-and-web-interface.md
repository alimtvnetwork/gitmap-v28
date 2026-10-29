# Subtask 04: Web UI Server REST Endpoints, API Route Mounting, and Settings Visualizer

> **Subtask ID:** 352-04  
> **Sequence ID:** 352  
> **Status:** PENDING  
> **Module Path:** `cli/cmdui/`  
> **Parent Plan:** `.ai-memory/plans/pending/352-cursor-like-ide-setup-from-gitmap.md`  
> **Specification Reference:** `02-spec/21-app/352-cursor-like-ide-setup-from-gitmap/02-ui-server-and-web-integration-spec.md`  

---

## 1. Objective & Scope

Integrate the interactive Cursor-like IDE configuration experience into the GitMap Fleet Web UI (`gitmap ui`):
1. Author REST API handlers in `cli/cmdui/ui_server_ide_setup.go` supporting editor detection (`/status`), preset catalogs (`/presets`), live connection testing (`/test-connection`), and atomic configuration application (`/apply`).
2. Mount the new API routes in `cli/cmdui/ui_server.go` within `mountAPIRoutes`.
3. Author frontend markup and reactive JavaScript in `cli/cmdui/ui_assets_markup_ide_setup.go` and `cli/cmdui/ui_assets_js_ide_setup.go` rendering a tabbed IDE workspace (Cline Composer, Continue.dev Inline & Tab, Zed AI, and Universal All Setup).
4. Provide real-time provider ping testing, credential masking, diff previews, and success notifications.

---

## 2. REST API Request & Response Contracts

All endpoints return standard GitMap `Result[T]` JSON envelopes:

### 2.1 GET `/api/ide/setup/status`
Detects installed editors and existing configuration files.

**Response Payload:**
```json
{
  "isSuccess": true,
  "data": {
    "cline": {
      "isInstalled": true,
      "configPath": "C:\\Users\\...\\cline_mcp_settings.json",
      "hasConfig": true,
      "activeProvider": "openai-compatible"
    },
    "continue": {
      "isInstalled": true,
      "configPath": "C:\\Users\\...\\.continue\\config.json",
      "hasConfig": true,
      "activeProvider": "openai"
    },
    "zed": {
      "isInstalled": true,
      "configPath": "C:\\Users\\...\\Zed\\settings.json",
      "hasConfig": true,
      "activeProvider": "openai"
    }
  }
}
```

### 2.2 POST `/api/ide/setup/test-connection`
Validates provider endpoint connectivity and measures round-trip response latency.

**Request Payload:**
```json
{
  "provider": "kimi",
  "baseUrl": "https://api.moonshot.cn/v1",
  "apiKey": "sk-...",
  "model": "kimi-k3.1-pro"
}
```

**Response Payload:**
```json
{
  "isSuccess": true,
  "data": {
    "isSuccess": true,
    "latencyMs": 138,
    "model": "kimi-k3.1-pro",
    "httpStatusCode": 200,
    "errorMessage": ""
  }
}
```

### 2.3 POST `/api/ide/setup/apply`
Applies AI configuration to selected IDE targets with pre-flight backup.

**Request Payload:**
```json
{
  "target": "all",
  "provider": "kimi",
  "baseUrl": "https://api.moonshot.cn/v1",
  "apiKey": "sk-...",
  "model": "kimi-k3.1-pro",
  "isDryRun": false,
  "isBackupEnabled": true
}
```

**Response Payload:**
```json
{
  "isSuccess": true,
  "data": {
    "isSuccess": true,
    "configuredTargets": ["cline", "continue", "zed"],
    "backupPaths": [
      "C:\\Users\\...\\.continue\\config.json.bak-20261029-153000"
    ],
    "isDryRun": false
  }
}
```

---

## 3. Server REST API Handlers (`cli/cmdui/ui_server_ide_setup.go`)

Flattened handlers strictly adhering to <= 15 line functions and positive booleans:

```go
package cmdui

import (
	"encoding/json"
	"net/http"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/setup"
)

func handleIDESetupStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	resolver := setup.NewPathResolver()
	statusSummary := collectIDEStatusSummary(resolver)

	writeJSONResponse(w, statusSummary)
}

func handleIDESetupPresets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	presets := setup.GetProviderPresets()
	writeJSONResponse(w, presets)
}

func handleIDESetupTestConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	var req setup.TestConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result := setup.TestConnection(req)
	writeJSONResponse(w, result)
}

func handleIDESetupApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	var opts setup.IDESetupOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := executeApplySetup(opts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, result)
}
```

---

## 4. Route Mounting (`cli/cmdui/ui_server.go`)

Wire endpoints into `mountAPIRoutes`:

```go
func mountAPIRoutes(mux *http.ServeMux) {
	// Existing API routes...
	mux.HandleFunc("/api/ide/setup/status", handleIDESetupStatus)
	mux.HandleFunc("/api/ide/setup/presets", handleIDESetupPresets)
	mux.HandleFunc("/api/ide/setup/test-connection", handleIDESetupTestConnection)
	mux.HandleFunc("/api/ide/setup/apply", handleIDESetupApply)
}
```

---

## 5. Web UI Frontend Markup & Components (`cli/cmdui/ui_assets_markup_ide_setup.go`)

Author semantic HTML embedded template for the tabbed settings visualizer:

```html
<div class="ide-setup-container">
  <div class="ide-setup-header">
    <h2>⚡ Cursor-Grade AI IDE Setup</h2>
    <p>Configure open-source and native editors into Cursor-equivalent AI environments powered by Moonshot Kimi K3.1 Pro or your preferred LLM provider.</p>
  </div>

  <div class="ide-tab-bar">
    <button class="ide-tab-btn active" data-tab="all">Universal Setup (All 3)</button>
    <button class="ide-tab-btn" data-tab="cline">Composer (VS Code + Cline)</button>
    <button class="ide-tab-btn" data-tab="continue">Inline & Tab (Continue.dev)</button>
    <button class="ide-tab-btn" data-tab="zed">Native Fast Editor (Zed AI)</button>
  </div>

  <div class="ide-config-panel">
    <div class="form-group">
      <label for="provider-select">AI Provider Preset</label>
      <select id="provider-select" class="form-control">
        <option value="kimi" selected>Moonshot Kimi K3.1 Pro (Recommended)</option>
        <option value="openrouter">OpenRouter Multi-Model</option>
        <option value="deepseek">DeepSeek AI</option>
        <option value="anthropic">Anthropic Claude</option>
        <option value="ollama">Local Ollama</option>
        <option value="custom">Custom Endpoint</option>
      </select>
    </div>

    <div class="form-group">
      <label for="base-url-input">API Base URL</label>
      <input type="text" id="base-url-input" class="form-control" value="https://api.moonshot.cn/v1" />
    </div>

    <div class="form-group">
      <label for="model-input">Model Name / ID</label>
      <input type="text" id="model-input" class="form-control" value="kimi-k3.1-pro" />
    </div>

    <div class="form-group">
      <label for="api-key-input">API Key</label>
      <div class="input-with-action">
        <input type="password" id="api-key-input" class="form-control" placeholder="sk-..." />
        <button id="toggle-key-visibility" class="btn-secondary">Show</button>
      </div>
      <small id="env-key-detected-badge" class="badge-subtle"></small>
    </div>

    <div class="actions-bar">
      <button id="test-connection-btn" class="btn-secondary">Test Connection</button>
      <button id="apply-setup-btn" class="btn-primary">Apply Configuration</button>
    </div>

    <div id="test-result-banner" class="status-banner" style="display: none;"></div>
  </div>
</div>
```

---

## 6. Frontend Reactive Controller (`cli/cmdui/ui_assets_js_ide_setup.go`)

Client-side event handlers coordinating presets, testing, and application:

```javascript
// Embedded JavaScript controller for IDE Setup Web UI
(function() {
  const providerSelect = document.getElementById('provider-select');
  const baseUrlInput = document.getElementById('base-url-input');
  const modelInput = document.getElementById('model-input');
  const apiKeyInput = document.getElementById('api-key-input');
  const testBtn = document.getElementById('test-connection-btn');
  const applyBtn = document.getElementById('apply-setup-btn');
  const banner = document.getElementById('test-result-banner');

  let activeTab = 'all';

  providerSelect.addEventListener('change', function(e) {
    const provider = e.target.value;
    updatePresetFields(provider);
  });

  testBtn.addEventListener('click', async function() {
    setLoadingState(testBtn, true);
    try {
      const resp = await fetch('/api/ide/setup/test-connection', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          provider: providerSelect.value,
          baseUrl: baseUrlInput.value,
          apiKey: apiKeyInput.value,
          model: modelInput.value
        })
      });
      const data = await resp.json();
      renderTestResult(data);
    } catch (err) {
      renderTestError(err.message);
    } finally {
      setLoadingState(testBtn, false);
    }
  });

  applyBtn.addEventListener('click', async function() {
    setLoadingState(applyBtn, true);
    try {
      const resp = await fetch('/api/ide/setup/apply', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          target: activeTab,
          provider: providerSelect.value,
          baseUrl: baseUrlInput.value,
          apiKey: apiKeyInput.value,
          model: modelInput.value,
          isBackupEnabled: true
        })
      });
      const data = await resp.json();
      renderApplySuccess(data);
    } catch (err) {
      renderApplyError(err.message);
    } finally {
      setLoadingState(applyBtn, false);
    }
  });
})();
```

---

## 7. Verification Criteria & Acceptance Tests

1. **`ui_server_ide_setup_test.go`**:
   - `GET /api/ide/setup/status` returns HTTP 200 and valid JSON with `isInstalled` flags.
   - `GET /api/ide/setup/presets` returns HTTP 200 with complete preset map containing `kimi`.
   - `POST /api/ide/setup/test-connection` handles valid and invalid payloads cleanly.
   - `POST /api/ide/setup/apply` in dry-run mode returns success without modifying host disk.
2. **Coding Guidelines Gate**:
   - All server handler functions <= 15 lines.
   - Strict positive boolean naming (`isSuccess`, `isInstalled`, `hasConfig`, `isBackupEnabled`).
   - Zero nested ifs.
