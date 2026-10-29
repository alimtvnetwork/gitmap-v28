package setup

import (
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// GenerateClineConfig writes GitMap MCP server tools and provider configuration for Cline.
func GenerateClineConfig(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	settingsPath, err := resolver.GetClineSettingsPath()
	if err != nil {
		return err
	}

	backupIfEnabled(settingsPath, opts.IsBackupEnabled)
	if opts.IsDryRun {
		return nil
	}

	return applyClineSettings(resolver, settingsPath, opts)
}

func backupIfEnabled(path string, isEnabled bool) {
	if isEnabled {
		_, _ = CreatePreFlightBackup(path)
	}
}

func applyClineSettings(resolver *PathResolver, path string, opts IDESetupOptions) *apperror.AppError {
	if mErr := MergeJSONFile(path, buildClineMCPUpdates()); mErr != nil {
		return mErr
	}

	return writeClineProviderState(resolver, opts)
}

func buildClineMCPUpdates() map[string]interface{} {
	return map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"gitmap": buildGitMapMCPEntry(),
		},
	}
}

func buildGitMapMCPEntry() map[string]interface{} {
	return map[string]interface{}{
		"command":     "gitmap",
		"args":        []string{"mcp", "serve"},
		"disabled":    false,
		"autoApprove": []string{"gitmap_scan_repos", "gitmap_search_files", "gitmap_get_repo_info"},
	}
}

func writeClineProviderState(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	storageDir, err := resolver.GetClineStorageDir()
	if err != nil {
		return err
	}

	statePath := filepath.Join(storageDir, "settings", "cline_provider_settings.json")
	backupIfEnabled(statePath, opts.IsBackupEnabled)
	payload := buildClineProviderPayload(opts)

	return MergeJSONFile(statePath, payload)
}

func buildClineProviderPayload(opts IDESetupOptions) map[string]interface{} {
	if opts.Provider.IsAnthropic() {
		return buildClineAnthropicPayload(opts)
	}

	return buildClineOpenAIPayload(opts)
}

func buildClineAnthropicPayload(opts IDESetupOptions) map[string]interface{} {
	return map[string]interface{}{
		"apiProvider":      "anthropic",
		"apiKey":           opts.APIKey,
		"anthropicModelId": opts.Model,
	}
}

func buildClineOpenAIPayload(opts IDESetupOptions) map[string]interface{} {
	return map[string]interface{}{
		"apiProvider":   "openai-compatible",
		"openAiBaseUrl": opts.BaseURL,
		"openAiApiKey":  opts.APIKey,
		"openAiModelId": opts.Model,
		"openAiCustomModelInfo": map[string]interface{}{
			"maxTokens":           8192,
			"contextWindow":       131072,
			"supportsImages":      true,
			"supportsPromptCache": true,
		},
	}
}

// ConfigureCline executes Cline setup and returns an IDETargetResult.
func ConfigureCline(resolver *PathResolver, opts IDESetupOptions) IDETargetResult {
	path, pathErr := resolver.GetClineSettingsPath()
	if pathErr != nil {
		return IDETargetResult{Target: TargetCline, IsConfigured: false, ErrorMessage: pathErr.Error()}
	}

	backupRes, _ := CreatePreFlightBackup(path)
	opts.IsBackupEnabled = false
	if err := GenerateClineConfig(resolver, opts); err != nil {
		return IDETargetResult{Target: TargetCline, ConfigPath: path, IsConfigured: false, ErrorMessage: err.Error()}
	}

	return buildClineSuccessResult(path, backupRes)
}

func buildClineSuccessResult(path string, backup BackupResult) IDETargetResult {
	return IDETargetResult{
		Target:          TargetCline,
		ConfigPath:      path,
		IsConfigured:    true,
		IsBackupCreated: backup.IsBackupCreated,
		BackupPath:      backup.BackupPath,
	}
}
