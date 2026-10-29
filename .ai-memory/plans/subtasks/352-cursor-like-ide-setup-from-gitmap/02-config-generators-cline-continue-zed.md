# Subtask 02: Specialized Configuration Generators (Cline, Continue.dev, Zed AI) & Safe Merger

> **Subtask ID:** 352-02  
> **Sequence ID:** 352  
> **Status:** PENDING  
> **Module Path:** `cli/cmdide/setup/`  
> **Parent Plan:** `.ai-memory/plans/pending/352-cursor-like-ide-setup-from-gitmap.md`  
> **Specification Reference:** `02-spec/21-app/352-cursor-like-ide-setup-from-gitmap/01-cursor-like-ide-setup-spec.md`  

---

## 1. Objective & Scope

Implement specialized configuration transformers and safe merging engines for the three target IDE ecosystems:
1. **VS Code + Cline Generator (`generator_cline.go`)**: Configures Cline's MCP settings with GitMap's repository indexing tools and sets the active AI model provider to Kimi K3.1 Pro (or selected preset).
2. **VS Code + Continue.dev Generator (`generator_continue.go`)**: Safely updates `~/.continue/config.json` with the primary chat/edit model and tab-autocomplete model without disturbing existing user slash commands or context providers.
3. **Zed AI Editor Generator (`generator_zed.go`)**: Injects provider endpoint definitions and assistant defaults into Zed's `settings.json` while strictly preserving editor themes, font sizes, and user keymaps.
4. **Safe JSON Merger (`merger.go`)**: Non-destructively deep-merges configuration blocks into existing JSON maps and applies atomic file renaming.
5. **Pre-Flight Backup Engine (`backup.go`)**: Enforces mandatory creation of timestamped backups (`<file>.bak-<timestamp>`) before any file modification.

---

## 2. Pre-Flight Backup Engine (`cli/cmdide/setup/backup.go`)

Enforce non-destructive safety: prior to writing to any configuration file, create a verified backup:

```go
package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// BackupResult details the outcome of a pre-flight backup operation.
type BackupResult struct {
	OriginalPath       string `json:"originalPath"`
	BackupPath         string `json:"backupPath"`
	IsBackupCreated    bool   `json:"isBackupCreated"`
}

// CreatePreFlightBackup creates a timestamped copy of the target file if it exists.
func CreatePreFlightBackup(targetPath string) (BackupResult, *apperror.AppError) {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return BackupResult{OriginalPath: targetPath, IsBackupCreated: false}, nil
	}

	timestamp := time.Now().UTC().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.bak-%s", targetPath, timestamp)

	data, err := os.ReadFile(targetPath)
	if err != nil {
		return BackupResult{}, apperror.WrapSimple(err, "CreatePreFlightBackup.ReadFile")
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return BackupResult{}, apperror.WrapSimple(err, "CreatePreFlightBackup.WriteFile")
	}

	return BackupResult{OriginalPath: targetPath, BackupPath: backupPath, IsBackupCreated: true}, nil
}
```

---

## 3. Safe JSON Deep Merger (`cli/cmdide/setup/merger.go`)

Merge configuration trees without clobbering existing user keys:

```go
package setup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// MergeJSONFile reads existing JSON, deep merges updates, and writes back atomically.
func MergeJSONFile(targetPath string, updates map[string]interface{}) *apperror.AppError {
	existingData := make(map[string]interface{})
	if _, err := os.Stat(targetPath); err == nil {
		readData, readErr := os.ReadFile(targetPath)
		if readErr != nil {
			return apperror.WrapSimple(readErr, "MergeJSONFile.ReadFile")
		}
		_ = json.Unmarshal(readData, &existingData)
	}

	deepMergeMaps(existingData, updates)

	outBytes, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "MergeJSONFile.Marshal")
	}

	return atomicWriteFile(targetPath, outBytes)
}

func deepMergeMaps(dest map[string]interface{}, src map[string]interface{}) {
	for key, srcVal := range src {
		destVal, exists := dest[key]
		if !exists {
			dest[key] = srcVal
			continue
		}

		destMap, isDestMap := destVal.(map[string]interface{})
		srcMap, isSrcMap := srcVal.(map[string]interface{})
		if isDestMap && isSrcMap {
			deepMergeMaps(destMap, srcMap)
			continue
		}

		dest[key] = srcVal
	}
}

func atomicWriteFile(targetPath string, data []byte) *apperror.AppError {
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "atomicWriteFile.MkdirAll")
	}

	tempFile := targetPath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return apperror.WrapSimple(err, "atomicWriteFile.WriteTemp")
	}

	if err := os.Rename(tempFile, targetPath); err != nil {
		_ = os.Remove(tempFile)
		return apperror.WrapSimple(err, "atomicWriteFile.Rename")
	}

	return nil
}
```

---

## 4. Cline Configuration Generator (`cli/cmdide/setup/generator_cline.go`)

Injects GitMap MCP tools and provider configuration into Cline storage:

```go
package setup

import (
	"path/filepath"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
)

// GenerateClineConfig writes GitMap MCP tools and provider configuration for Cline.
func GenerateClineConfig(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	settingsPath, err := resolver.GetClineSettingsPath()
	if err != nil {
		return err
	}

	mcpUpdates := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"gitmap": map[string]interface{}{
				"command":  "gitmap",
				"args":     []string{"mcp", "serve"},
				"disabled": false,
				"autoApprove": []string{
					"gitmap_scan_repos",
					"gitmap_search_files",
					"gitmap_get_repo_info",
				},
			},
		},
	}

	if opts.IsBackupEnabled {
		if _, bErr := CreatePreFlightBackup(settingsPath); bErr != nil {
			return bErr
		}
	}

	if opts.IsDryRun {
		return nil
	}

	if mErr := MergeJSONFile(settingsPath, mcpUpdates); mErr != nil {
		return mErr
	}

	return writeClineProviderState(resolver, opts)
}

func writeClineProviderState(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	storageDir, err := resolver.GetClineStorageDir()
	if err != nil {
		return err
	}

	statePath := filepath.Join(storageDir, "settings", "cline_provider_settings.json")
	providerPayload := buildClineProviderPayload(opts)

	return MergeJSONFile(statePath, providerPayload)
}

func buildClineProviderPayload(opts IDESetupOptions) map[string]interface{} {
	if opts.Provider == enums.ProviderAnthropic {
		return map[string]interface{}{
			"apiProvider":      "anthropic",
			"apiKey":           opts.APIKey,
			"anthropicModelId": opts.Model,
		}
	}

	return map[string]interface{}{
		"apiProvider":   "openai-compatible",
		"openAiBaseUrl": opts.BaseURL,
		"openAiApiKey":  opts.APIKey,
		"openAiModelId": opts.Model,
	}
}
```

---

## 5. Continue.dev Configuration Generator (`cli/cmdide/setup/generator_continue.go`)

Configures primary chat and autocomplete models in `~/.continue/config.json`:

```go
package setup

import (
	"encoding/json"
	"os"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
)

// GenerateContinueConfig merges provider configuration into ~/.continue/config.json.
func GenerateContinueConfig(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	targetPath, err := resolver.GetContinueConfigPath()
	if err != nil {
		return err
	}

	if opts.IsBackupEnabled {
		if _, bErr := CreatePreFlightBackup(targetPath); bErr != nil {
			return bErr
		}
	}

	if opts.IsDryRun {
		return nil
	}

	existing := readExistingJSONMap(targetPath)
	newModel := buildContinueModelEntry(opts)
	existing["models"] = upsertContinueModelList(existing["models"], newModel)
	existing["tabAutocompleteModel"] = buildContinueAutocompleteEntry(opts)

	outBytes, mErr := json.MarshalIndent(existing, "", "  ")
	if mErr != nil {
		return apperror.WrapSimple(mErr, "GenerateContinueConfig.Marshal")
	}

	return atomicWriteFile(targetPath, outBytes)
}

func buildContinueModelEntry(opts IDESetupOptions) map[string]interface{} {
	providerType := "openai"
	if opts.Provider == enums.ProviderAnthropic {
		providerType = "anthropic"
	}

	return map[string]interface{}{
		"title":    "Moonshot Kimi K3.1 Pro (GitMap)",
		"provider": providerType,
		"model":    opts.Model,
		"apiBase":  opts.BaseURL,
		"apiKey":   opts.APIKey,
	}
}

func buildContinueAutocompleteEntry(opts IDESetupOptions) map[string]interface{} {
	return map[string]interface{}{
		"title":    "GitMap Autocomplete",
		"provider": "openai",
		"model":    opts.Model,
		"apiBase":  opts.BaseURL,
		"apiKey":   opts.APIKey,
	}
}

func upsertContinueModelList(rawModels interface{}, newModel map[string]interface{}) []interface{} {
	list, isList := rawModels.([]interface{})
	if !isList {
		return []interface{}{newModel}
	}

	filtered := make([]interface{}, 0, len(list)+1)
	filtered = append(filtered, newModel)
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if ok && m["title"] == newModel["title"] {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered
}

func readExistingJSONMap(path string) map[string]interface{} {
	result := make(map[string]interface{})
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &result)
	}

	return result
}
```

---

## 6. Zed AI Configuration Generator (`cli/cmdide/setup/generator_zed.go`)

Configures `language_models` in Zed `settings.json`:

```go
package setup

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
)

// GenerateZedConfig merges provider endpoints and assistant settings into Zed settings.json.
func GenerateZedConfig(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	targetPath, err := resolver.GetZedSettingsPath()
	if err != nil {
		return err
	}

	if opts.IsBackupEnabled {
		if _, bErr := CreatePreFlightBackup(targetPath); bErr != nil {
			return bErr
		}
	}

	if opts.IsDryRun {
		return nil
	}

	zedUpdates := buildZedSettingsPayload(opts)
	return MergeJSONFile(targetPath, zedUpdates)
}

func buildZedSettingsPayload(opts IDESetupOptions) map[string]interface{} {
	if opts.Provider == enums.ProviderAnthropic {
		return map[string]interface{}{
			"assistant": map[string]interface{}{
				"default_model": map[string]interface{}{
					"provider": "anthropic",
					"model":    opts.Model,
				},
				"version": "2",
			},
		}
	}

	return map[string]interface{}{
		"language_models": map[string]interface{}{
			"openai": map[string]interface{}{
				"custom_endpoint": opts.BaseURL,
				"api_url":         opts.BaseURL,
				"available_models": []string{
					opts.Model,
				},
			},
		},
		"assistant": map[string]interface{}{
			"default_model": map[string]interface{}{
				"provider": "openai",
				"model":    opts.Model,
			},
			"version": "2",
		},
	}
}
```

---

## 7. Verification Criteria & Acceptance Tests

1. **`backup_test.go`**:
   - Creates a temporary file and invokes `CreatePreFlightBackup`.
   - Asserts backup file exists, contains exact original byte contents, and returns `IsBackupCreated: true`.
   - Asserts non-existent target returns `IsBackupCreated: false` without error.
2. **`merger_test.go`**:
   - Merges updates into a JSON file containing existing custom themes and keybindings.
   - Verifies existing keys are preserved 100% intact while new AI keys are correctly inserted.
   - Verifies temporary file clean-up upon failure.
3. **`generator_cline_test.go`**:
   - Verifies `mcpServers.gitmap` contains `command: "gitmap"` and required auto-approved tools.
   - Verifies provider state JSON correctly serializes OpenAI-compatible Kimi configuration.
4. **`generator_continue_test.go`**:
   - Verifies `models` array retains previous custom models while prepending Kimi K3.1 Pro.
   - Verifies `tabAutocompleteModel` is correctly populated.
5. **`generator_zed_test.go`**:
   - Verifies `language_models.openai.custom_endpoint` matches provider URL.
   - Verifies `assistant.default_model` references correct provider and model name.
