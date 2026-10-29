package setup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// GenerateContinueConfig merges provider and autocomplete configuration into ~/.continue/config.json.
func GenerateContinueConfig(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	targetPath, err := resolver.GetContinueConfigPath()
	if err != nil {
		return err
	}

	backupIfEnabled(targetPath, opts.IsBackupEnabled)

	if opts.IsDryRun {
		return nil
	}

	return writeContinueConfig(targetPath, opts)
}

func writeContinueConfig(targetPath string, opts IDESetupOptions) *apperror.AppError {
	existing := readExistingJSON(targetPath)
	newModel := buildContinueModelEntry(opts)
	existing["models"] = upsertContinueModelList(existing["models"], newModel)
	existing["tabAutocompleteModel"] = buildContinueAutocompleteEntry(opts)
	existing["allowAnonymousTelemetry"] = false

	outBytes, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return apperror.WrapExecution(err, "GenerateContinueConfig.Marshal")
	}

	return AtomicWriteFile(targetPath, outBytes)
}

func resolveContinueModelTitle(provider IDEProviderType, model string) string {
	if provider.IsKimi() {
		return "Moonshot Kimi K3.1 Pro (GitMap)"
	}

	return fmt.Sprintf("GitMap (%s)", model)
}

func resolveContinueProviderName(provider IDEProviderType) string {
	if provider.IsAnthropic() {
		return "anthropic"
	}

	return "openai"
}

func buildContinueModelEntry(opts IDESetupOptions) map[string]interface{} {
	return map[string]interface{}{
		"title":         resolveContinueModelTitle(opts.Provider, opts.Model),
		"provider":      resolveContinueProviderName(opts.Provider),
		"model":         opts.Model,
		"apiBase":       opts.BaseURL,
		"apiKey":        opts.APIKey,
		"contextLength": resolveContinueContextLimit(opts.Provider),
		"roles":         []string{"chat", "edit"},
	}
}

func resolveContinueContextLimit(provider IDEProviderType) int {
	preset, isFound := GetPreset(provider)
	if isFound && preset.ContextLimit > 0 {
		return preset.ContextLimit
	}

	return 128000
}

func buildContinueAutocompleteEntry(opts IDESetupOptions) map[string]interface{} {
	return map[string]interface{}{
		"title":    "GitMap Autocomplete",
		"provider": resolveContinueProviderName(opts.Provider),
		"model":    opts.Model,
		"apiBase":  opts.BaseURL,
		"apiKey":   opts.APIKey,
	}
}

func upsertContinueModelList(rawModels any, newModel map[string]interface{}) []any {
	list, isList := rawModels.([]any)
	if !isList {
		return []any{newModel}
	}

	return filterContinueModels(list, newModel)
}

func filterContinueModels(list []any, newModel map[string]interface{}) []any {
	filtered := make([]any, 0, len(list)+1)
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

func readExistingJSON(path string) map[string]interface{} {
	result := make(map[string]interface{})
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &result)
	}

	return result
}

// ConfigureContinue executes Continue setup and returns an IDETargetResult.
func ConfigureContinue(resolver *PathResolver, opts IDESetupOptions) IDETargetResult {
	path, pathErr := resolver.GetContinueConfigPath()
	if pathErr != nil {
		return IDETargetResult{Target: TargetContinue, IsConfigured: false, ErrorMessage: pathErr.Error()}
	}

	backupRes, _ := CreatePreFlightBackup(path)
	opts.IsBackupEnabled = false
	if err := GenerateContinueConfig(resolver, opts); err != nil {
		return IDETargetResult{Target: TargetContinue, ConfigPath: path, IsConfigured: false, ErrorMessage: err.Error()}
	}

	return buildContinueSuccessResult(path, backupRes)
}

func buildContinueSuccessResult(path string, backup BackupResult) IDETargetResult {
	return IDETargetResult{
		Target:          TargetContinue,
		ConfigPath:      path,
		IsConfigured:    true,
		IsBackupCreated: backup.IsBackupCreated,
		BackupPath:      backup.BackupPath,
	}
}
