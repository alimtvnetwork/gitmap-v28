package setup

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// GenerateZedConfig merges provider endpoints and assistant settings into Zed settings.json.
func GenerateZedConfig(resolver *PathResolver, opts IDESetupOptions) *apperror.AppError {
	targetPath, err := resolver.GetZedSettingsPath()
	if err != nil {
		return err
	}

	backupIfEnabled(targetPath, opts.IsBackupEnabled)

	if opts.IsDryRun {
		return nil
	}

	return MergeJSONFile(targetPath, buildZedSettingsPayload(opts))
}

func buildZedSettingsPayload(opts IDESetupOptions) map[string]interface{} {
	if opts.Provider.IsAnthropic() {
		return map[string]interface{}{
			"assistant": buildZedAssistantModel("anthropic", opts.Model),
		}
	}

	return buildZedOpenAIPayload(opts)
}

func buildZedOpenAIPayload(opts IDESetupOptions) map[string]interface{} {
	return map[string]interface{}{
		"language_models": buildZedLanguageModels(opts),
		"assistant":       buildZedAssistantModel("openai", opts.Model),
	}
}

func buildZedLanguageModels(opts IDESetupOptions) map[string]interface{} {
	return map[string]interface{}{
		"openai": map[string]interface{}{
			"custom_endpoint": opts.BaseURL,
			"api_url":         opts.BaseURL,
			"available_models": []string{
				opts.Model,
			},
		},
	}
}

func buildZedAssistantModel(provider string, model string) map[string]interface{} {
	return map[string]interface{}{
		"default_model": map[string]interface{}{
			"provider": provider,
			"model":    model,
		},
		"version": "2",
	}
}

// ConfigureZed executes Zed setup and returns an IDETargetResult.
func ConfigureZed(resolver *PathResolver, opts IDESetupOptions) IDETargetResult {
	path, pathErr := resolver.GetZedSettingsPath()
	if pathErr != nil {
		return IDETargetResult{Target: TargetZed, IsConfigured: false, ErrorMessage: pathErr.Error()}
	}

	backupRes, _ := CreatePreFlightBackup(path)
	opts.IsBackupEnabled = false
	if err := GenerateZedConfig(resolver, opts); err != nil {
		return IDETargetResult{Target: TargetZed, ConfigPath: path, IsConfigured: false, ErrorMessage: err.Error()}
	}

	return buildZedSuccessResult(path, backupRes)
}

func buildZedSuccessResult(path string, backup BackupResult) IDETargetResult {
	return IDETargetResult{
		Target:          TargetZed,
		ConfigPath:      path,
		IsConfigured:    true,
		IsBackupCreated: backup.IsBackupCreated,
		BackupPath:      backup.BackupPath,
	}
}

// RunSetup executes configuration across selected target IDEs.
func RunSetup(opts IDESetupOptions) (*IDESetupResult, error) {
	ApplyPresetDefaults(&opts)
	res := initSetupResult(opts)

	if err := runTestIfRequested(opts, res); err != nil {
		return res, err
	}

	executeTargetSetups(NewPathResolver(), opts, res)
	res.SummaryMessage = formatSummaryMessage(*res)

	return res, nil
}

func initSetupResult(opts IDESetupOptions) *IDESetupResult {
	return &IDESetupResult{
		Target:        opts.Target,
		Provider:      opts.Provider,
		IsSuccess:     true,
		TargetResults: make(map[string]IDETargetResult),
	}
}

func runTestIfRequested(opts IDESetupOptions, res *IDESetupResult) error {
	if !opts.IsTestEnabled {
		return nil
	}

	testRes, err := RunConnectionTest(opts)
	if testRes != nil {
		res.ConnectionTest = *testRes
	}

	return evaluateTestResult(opts, res, err)
}

func evaluateTestResult(opts IDESetupOptions, res *IDESetupResult, err error) error {
	if err != nil && !opts.IsForceCreate {
		res.IsSuccess = false
		res.SummaryMessage = "Connection test failed: " + err.Error()

		return err
	}

	return nil
}

func executeTargetSetups(resolver *PathResolver, opts IDESetupOptions, res *IDESetupResult) {
	if opts.Target.IsCline() || opts.Target.IsAll() {
		recordTargetResult(res, TargetCline, ConfigureCline(resolver, opts))
	}

	if opts.Target.IsContinue() || opts.Target.IsAll() {
		recordTargetResult(res, TargetContinue, ConfigureContinue(resolver, opts))
	}

	if opts.Target.IsZed() || opts.Target.IsAll() {
		recordTargetResult(res, TargetZed, ConfigureZed(resolver, opts))
	}
}

func recordTargetResult(res *IDESetupResult, key IDETargetType, targetRes IDETargetResult) {
	res.TargetResults[string(key)] = targetRes
	if !targetRes.IsConfigured {
		res.IsSuccess = false
	}

	if targetRes.IsBackupCreated && targetRes.BackupPath != "" {
		res.BackupPaths = append(res.BackupPaths, targetRes.BackupPath)
	}
}

func formatSummaryMessage(res IDESetupResult) string {
	if !res.IsSuccess {
		return "Setup completed with errors"
	}

	return fmt.Sprintf("Successfully configured %d target(s) for provider %s", len(res.TargetResults), res.Provider)
}
