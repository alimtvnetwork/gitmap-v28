# Subtask 01: Core IDE Setup Engine, Target Enums, and Provider Presets

> **Subtask ID:** 352-01  
> **Sequence ID:** 352  
> **Status:** PENDING  
> **Module Path:** `cli/cmdide/setup/`, `cli/cmdide/enums/`  
> **Parent Plan:** `.ai-memory/plans/pending/352-cursor-like-ide-setup-from-gitmap.md`  
> **Specification Reference:** `02-spec/21-app/352-cursor-like-ide-setup-from-gitmap/01-cursor-like-ide-setup-spec.md`  

---

## 1. Objective & Scope

Establish the core domain foundations for the Cursor-like IDE setup engine in GitMap:
1. Define strongly typed enums for target IDEs (`IDETargetType`) and provider presets (`IDEProviderType`).
2. Construct the provider presets catalog, pre-configuring **Kimi K3.1 Pro (Moonshot AI)** as the recommended default alongside OpenRouter, DeepSeek, Anthropic Claude, Local Ollama, and Custom endpoints.
3. Architect the cross-platform path resolver (`PathResolver`) dynamically mapping configuration locations across Windows, macOS, and Linux without hardcoded path separators.
4. Implement the connection test engine (`ConnectionTester`) validating API keys and measuring network latency before configuration persistence.

---

## 2. Strongly Typed Enums (`cli/cmdide/enums/`)

### 2.1 Provider Preset Enum (`ide_provider_enum.go`)

```go
package enums

import (
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IDEProviderType represents an LLM provider preset.
type IDEProviderType string

const (
	ProviderKimi       IDEProviderType = "kimi"
	ProviderOpenRouter IDEProviderType = "openrouter"
	ProviderDeepSeek   IDEProviderType = "deepseek"
	ProviderAnthropic  IDEProviderType = "anthropic"
	ProviderOllama     IDEProviderType = "ollama"
	ProviderCustom     IDEProviderType = "custom"
)

// IsValid reports whether the provider variant is known.
func (p IDEProviderType) IsValid() bool {
	switch p {
	case ProviderKimi, ProviderOpenRouter, ProviderDeepSeek, ProviderAnthropic, ProviderOllama, ProviderCustom:
		return true
	default:
		return false
	}
}

// String returns the string representation of the provider.
func (p IDEProviderType) String() string {
	return string(p)
}

// ParseIDEProviderType validates and returns an IDEProviderType.
func ParseIDEProviderType(raw string) (IDEProviderType, *apperror.AppError) {
	clean := strings.ToLower(strings.TrimSpace(raw))
	switch clean {
	case "kimi", "moonshot":
		return ProviderKimi, nil
	case "openrouter", "open-router":
		return ProviderOpenRouter, nil
	case "deepseek", "deep-seek":
		return ProviderDeepSeek, nil
	case "anthropic", "claude":
		return ProviderAnthropic, nil
	case "ollama", "local":
		return ProviderOllama, nil
	case "custom":
		return ProviderCustom, nil
	default:
		return "", apperror.NewValidation("ParseIDEProviderType", "E1001", "unsupported provider: "+raw)
	}
}
```

### 2.2 Target Enum Extensions (`ide_target_enum.go`)

Extend `IDETargetType` to support the dedicated Cursor-like setup targets:
```go
const (
	TargetCline    IDETargetType = "cline"
	TargetContinue IDETargetType = "continue"
	TargetZed      IDETargetType = "zed"
	TargetAllAI    IDETargetType = "all"
)
```

---

## 3. Core Models & Option Structs (`cli/cmdide/setup/types.go`)

Enforce strict positive boolean naming (`is*`, `has*`):

```go
package setup

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"

// IDESetupOptions controls command execution and behavioral flags.
type IDESetupOptions struct {
	Target                enums.IDETargetType   `json:"target"`
	Provider              enums.IDEProviderType `json:"provider"`
	Model                 string                `json:"model,omitempty"`
	BaseURL               string                `json:"baseUrl,omitempty"`
	APIKey                string                `json:"apiKey,omitempty"`
	IsDryRun              bool                  `json:"isDryRun"`
	IsQuiet               bool                  `json:"isQuiet"`
	IsJSON                bool                  `json:"isJSON"`
	IsForceCreate         bool                  `json:"isForceCreate"`
	IsTestEnabled         bool                  `json:"isTestEnabled"`
	IsBackupEnabled       bool                  `json:"isBackupEnabled"`
}

// ProviderPresetConfig defines default network and parameter metadata for a provider.
type ProviderPresetConfig struct {
	Provider              enums.IDEProviderType `json:"provider"`
	DisplayName           string                `json:"displayName"`
	DefaultBaseURL        string                `json:"defaultBaseUrl"`
	DefaultModel          string                `json:"defaultModel"`
	FallbackModels        []string              `json:"fallbackModels"`
	ContextLimit          int                   `json:"contextLimit"`
	AuthHeaderScheme      string                `json:"authHeaderScheme"`
	EnvVarKeys            []string              `json:"envVarKeys"`
	HasStreamingSupport   bool                  `json:"hasStreamingSupport"`
}

// TestConnectionRequest models the input for live connectivity testing.
type TestConnectionRequest struct {
	Provider              enums.IDEProviderType `json:"provider"`
	BaseURL               string                `json:"baseUrl"`
	APIKey                string                `json:"apiKey"`
	Model                 string                `json:"model"`
}

// TestConnectionResult models the output of live connectivity testing.
type TestConnectionResult struct {
	IsSuccess             bool                  `json:"isSuccess"`
	LatencyMs             int64                 `json:"latencyMs"`
	Model                 string                `json:"model"`
	ErrorMessage          string                `json:"errorMessage,omitempty"`
	HTTPStatusCode        int                   `json:"httpStatusCode"`
}
```

---

## 4. Provider Presets Registry (`cli/cmdide/setup/presets.go`)

Implement deterministic preset catalog and environment cascade resolver:

```go
package setup

import (
	"os"
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
)

// GetProviderPresets returns all production preset definitions.
func GetProviderPresets() map[enums.IDEProviderType]ProviderPresetConfig {
	return map[enums.IDEProviderType]ProviderPresetConfig{
		enums.ProviderKimi: {
			Provider:            enums.ProviderKimi,
			DisplayName:         "Moonshot Kimi K3.1 Pro (Recommended)",
			DefaultBaseURL:      "https://api.moonshot.cn/v1",
			DefaultModel:        "kimi-k3.1-pro",
			FallbackModels:      []string{"moonshot-v1-128k", "moonshot-v1-32k"},
			ContextLimit:        128000,
			AuthHeaderScheme:    "Bearer",
			EnvVarKeys:          []string{"KIMI_API_KEY", "MOONSHOT_API_KEY"},
			HasStreamingSupport: true,
		},
		enums.ProviderOpenRouter: {
			Provider:            enums.ProviderOpenRouter,
			DisplayName:         "OpenRouter Multi-Model",
			DefaultBaseURL:      "https://openrouter.ai/api/v1",
			DefaultModel:        "moonshotai/kimi-k3.1-pro",
			FallbackModels:      []string{"anthropic/claude-3.7-sonnet", "deepseek/deepseek-chat"},
			ContextLimit:        128000,
			AuthHeaderScheme:    "Bearer",
			EnvVarKeys:          []string{"OPENROUTER_API_KEY"},
			HasStreamingSupport: true,
		},
		enums.ProviderDeepSeek: {
			Provider:            enums.ProviderDeepSeek,
			DisplayName:         "DeepSeek AI",
			DefaultBaseURL:      "https://api.deepseek.com/v1",
			DefaultModel:        "deepseek-chat",
			FallbackModels:      []string{"deepseek-reasoner"},
			ContextLimit:        64000,
			AuthHeaderScheme:    "Bearer",
			EnvVarKeys:          []string{"DEEPSEEK_API_KEY"},
			HasStreamingSupport: true,
		},
		enums.ProviderAnthropic: {
			Provider:            enums.ProviderAnthropic,
			DisplayName:         "Anthropic Claude",
			DefaultBaseURL:      "https://api.anthropic.com/v1",
			DefaultModel:        "claude-3-7-sonnet-20250219",
			FallbackModels:      []string{"claude-3-5-sonnet-latest"},
			ContextLimit:        200000,
			AuthHeaderScheme:    "x-api-key",
			EnvVarKeys:          []string{"ANTHROPIC_API_KEY"},
			HasStreamingSupport: true,
		},
		enums.ProviderOllama: {
			Provider:            enums.ProviderOllama,
			DisplayName:         "Local Ollama Server",
			DefaultBaseURL:      "http://localhost:11434/v1",
			DefaultModel:        "deepseek-r1",
			FallbackModels:      []string{"qwen2.5-coder:32b"},
			ContextLimit:        32000,
			AuthHeaderScheme:    "Bearer",
			EnvVarKeys:          []string{"OLLAMA_API_KEY"},
			HasStreamingSupport: true,
		},
	}
}

// ResolveAPIKey determines the active API key through environment cascades.
func ResolveAPIKey(provider enums.IDEProviderType, explicitKey string) string {
	if strings.TrimSpace(explicitKey) != "" {
		return strings.TrimSpace(explicitKey)
	}

	preset, exists := GetProviderPresets()[provider]
	if exists {
		for _, envKey := range preset.EnvVarKeys {
			if val := os.Getenv(envKey); val != "" {
				return val
			}
		}
	}

	if val := os.Getenv("GITMAP_AI_API_KEY"); val != "" {
		return val
	}

	return os.Getenv("OPENAI_API_KEY")
}
```

---

## 5. Cross-Platform Path Resolver (`cli/cmdide/setup/path_resolver.go`)

Construct dynamic filesystem paths using `os.UserHomeDir`, `os.Getenv`, and `filepath.Join` without hardcoded separators:

```go
package setup

import (
	"os"
	"path/filepath"
	"runtime"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PathResolver exposes cross-platform configuration discovery.
type PathResolver struct{}

// NewPathResolver creates a platform-aware PathResolver.
func NewPathResolver() *PathResolver {
	return &PathResolver{}
}

// GetClineSettingsPath returns the absolute path to cline_mcp_settings.json.
func (r *PathResolver) GetClineSettingsPath() (string, *apperror.AppError) {
	storageDir, err := r.GetClineStorageDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(storageDir, "settings", "cline_mcp_settings.json"), nil
}

// GetClineStorageDir locates Cline global storage directory across OS platforms.
func (r *PathResolver) GetClineStorageDir() (string, *apperror.AppError) {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", apperror.NewSystem("GetClineStorageDir", "E2001", "APPDATA environment variable is empty")
		}
		return filepath.Join(appData, "Code", "User", "globalStorage", "saoudrizwan.claude-dev"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "GetClineStorageDir.UserHomeDir")
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "saoudrizwan.claude-dev"), nil
	}

	return filepath.Join(home, ".config", "Code", "User", "globalStorage", "saoudrizwan.claude-dev"), nil
}

// GetContinueConfigPath locates ~/.continue/config.json.
func (r *PathResolver) GetContinueConfigPath() (string, *apperror.AppError) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "GetContinueConfigPath.UserHomeDir")
	}

	return filepath.Join(home, ".continue", "config.json"), nil
}

// GetZedSettingsPath locates Zed settings.json across OS platforms.
func (r *PathResolver) GetZedSettingsPath() (string, *apperror.AppError) {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", apperror.NewSystem("GetZedSettingsPath", "E2002", "APPDATA environment variable is empty")
		}
		return filepath.Join(appData, "Zed", "settings.json"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "GetZedSettingsPath.UserHomeDir")
	}

	return filepath.Join(home, ".config", "zed", "settings.json"), nil
}
```

---

## 6. Connection Testing Engine (`cli/cmdide/setup/connection_tester.go`)

Performs a lightweight HTTP ping against the provider base URL to confirm API key validity and measure response latency:

```go
package setup

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
)

// TestConnection executes an authorization verification ping to the provider.
func TestConnection(req TestConnectionRequest) TestConnectionResult {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	endpointURL := buildTestURL(req.BaseURL, req.Provider)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return TestConnectionResult{IsSuccess: false, ErrorMessage: err.Error()}
	}

	applyAuthHeaders(httpReq, req.Provider, req.APIKey)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return TestConnectionResult{IsSuccess: false, ErrorMessage: err.Error()}
	}
	defer resp.Body.Close()

	isSuccess := resp.StatusCode >= 200 && resp.StatusCode < 300
	return TestConnectionResult{
		IsSuccess:      isSuccess,
		LatencyMs:      time.Since(start).Milliseconds(),
		Model:          req.Model,
		HTTPStatusCode: resp.StatusCode,
		ErrorMessage:   formatErrorSummary(resp.StatusCode),
	}
}

func buildTestURL(baseURL string, provider enums.IDEProviderType) string {
	clean := strings.TrimRight(baseURL, "/")
	if provider == enums.ProviderAnthropic {
		return clean + "/models"
	}
	return clean + "/models"
}

func applyAuthHeaders(req *http.Request, provider enums.IDEProviderType, key string) {
	if provider == enums.ProviderAnthropic {
		req.Header.Set("x-api-key", key)
		req.Header.Set("anthropic-version", "2023-06-01")
		return
	}
	req.Header.Set("Authorization", "Bearer "+key)
}

func formatErrorSummary(code int) string {
	if code >= 200 && code < 300 {
		return ""
	}
	return fmt.Sprintf("provider returned status code %d", code)
}
```

---

## 7. Verification Criteria & Acceptance Tests

1. **`ide_provider_enum_test.go`**:
   - `ParseIDEProviderType` cleanly maps `"kimi"`, `"moonshot"`, `"openrouter"`, `"deepseek"`, `"anthropic"`, `"ollama"`, and `"custom"`.
   - Rejects unknown variants with strongly typed `*apperror.AppError`.
2. **`presets_test.go`**:
   - Validates that Kimi K3.1 Pro is populated with `https://api.moonshot.cn/v1` and context limit 128,000.
   - Tests environment variable cascade fallback (`KIMI_API_KEY` -> `GITMAP_AI_API_KEY` -> `OPENAI_API_KEY`).
3. **`path_resolver_test.go`**:
   - Simulates Windows, macOS, and Linux paths using mocked environment variables.
   - Asserts that all returned paths use correct system separators and never return empty strings on valid environments.
4. **`connection_tester_test.go`**:
   - Uses `httptest.NewServer` to mock successful 200 responses and failed 401 unauthorized responses.
   - Verifies `isSuccess`, `latencyMs`, and `httpStatusCode` values.
5. **Coding Guidelines Gate**:
   - All functions <= 15 lines of executable code.
   - 100% positive booleans (`is*`, `has*`).
   - Zero nested ifs.
