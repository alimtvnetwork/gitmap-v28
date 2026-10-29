package setup

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IDETargetType represents a supported IDE or editor target for AI configuration.
type IDETargetType string

const (
	TargetCline    IDETargetType = "cline"
	TargetContinue IDETargetType = "continue"
	TargetZed      IDETargetType = "zed"
	TargetAll      IDETargetType = "all"
	TargetAllAI    IDETargetType = "all"
)

var targetAliasMap = map[string]IDETargetType{
	"cline": TargetCline, "composer": TargetCline,
	"continue": TargetContinue, "tab": TargetContinue,
	"zed": TargetZed,
	"all": TargetAll, "allai": TargetAll, "all-ai": TargetAll, "": TargetAll,
}

// IsValid reports whether the target variant is known.
func (t IDETargetType) IsValid() bool {
	switch t {
	case TargetCline, TargetContinue, TargetZed, TargetAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the target.
func (t IDETargetType) String() string {
	return string(t)
}

// IsCline reports whether the target is TargetCline.
func (t IDETargetType) IsCline() bool {
	return t == TargetCline
}

// IsContinue reports whether the target is TargetContinue.
func (t IDETargetType) IsContinue() bool {
	return t == TargetContinue
}

// IsZed reports whether the target is TargetZed.
func (t IDETargetType) IsZed() bool {
	return t == TargetZed
}

// IsAll reports whether the target is TargetAll.
func (t IDETargetType) IsAll() bool {
	return t == TargetAll || t == TargetAllAI
}

// ParseIDETargetType validates and returns an IDETargetType.
func ParseIDETargetType(raw string) (IDETargetType, *apperror.AppError) {
	clean := strings.ToLower(strings.TrimSpace(raw))
	if target, exists := targetAliasMap[clean]; exists {
		return target, nil
	}

	return "", apperror.NewValidation("ParseIDETargetType", "E1000", "invalid ide target: "+raw)
}

// AllIDETargetTypes returns all valid targets.
func AllIDETargetTypes() []IDETargetType {
	return []IDETargetType{
		TargetCline,
		TargetContinue,
		TargetZed,
		TargetAll,
	}
}

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

var providerAliasMap = map[string]IDEProviderType{
	"kimi": ProviderKimi, "moonshot": ProviderKimi, "moonshotai": ProviderKimi,
	"openrouter": ProviderOpenRouter, "open-router": ProviderOpenRouter,
	"deepseek": ProviderDeepSeek, "deep-seek": ProviderDeepSeek,
	"anthropic": ProviderAnthropic, "claude": ProviderAnthropic,
	"ollama": ProviderOllama, "local": ProviderOllama,
	"custom": ProviderCustom,
}

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

// IsKimi reports whether the provider is ProviderKimi.
func (p IDEProviderType) IsKimi() bool {
	return p == ProviderKimi
}

// IsOpenRouter reports whether the provider is ProviderOpenRouter.
func (p IDEProviderType) IsOpenRouter() bool {
	return p == ProviderOpenRouter
}

// IsDeepSeek reports whether the provider is ProviderDeepSeek.
func (p IDEProviderType) IsDeepSeek() bool {
	return p == ProviderDeepSeek
}

// IsAnthropic reports whether the provider is ProviderAnthropic.
func (p IDEProviderType) IsAnthropic() bool {
	return p == ProviderAnthropic
}

// IsOllama reports whether the provider is ProviderOllama.
func (p IDEProviderType) IsOllama() bool {
	return p == ProviderOllama
}

// IsCustom reports whether the provider is ProviderCustom.
func (p IDEProviderType) IsCustom() bool {
	return p == ProviderCustom
}

// ParseIDEProviderType validates and returns an IDEProviderType.
func ParseIDEProviderType(raw string) (IDEProviderType, *apperror.AppError) {
	clean := strings.ToLower(strings.TrimSpace(raw))
	if provider, exists := providerAliasMap[clean]; exists {
		return provider, nil
	}

	return "", apperror.NewValidation("ParseIDEProviderType", "E1001", "unsupported provider: "+raw)
}

// AllIDEProviderTypes returns all valid providers.
func AllIDEProviderTypes() []IDEProviderType {
	return []IDEProviderType{
		ProviderKimi,
		ProviderOpenRouter,
		ProviderDeepSeek,
		ProviderAnthropic,
		ProviderOllama,
		ProviderCustom,
	}
}

// IDESetupOptions controls command execution and behavioral flags.
type IDESetupOptions struct {
	Target          IDETargetType   `json:"target"`
	Provider        IDEProviderType `json:"provider"`
	Model           string          `json:"model,omitempty"`
	BaseURL         string          `json:"baseUrl,omitempty"`
	APIKey          string          `json:"apiKey,omitempty"`
	IsDryRun        bool            `json:"isDryRun"`
	IsQuiet         bool            `json:"isQuiet"`
	IsJSON          bool            `json:"isJSON"`
	IsForceCreate   bool            `json:"isForceCreate"`
	IsTestEnabled   bool            `json:"isTestEnabled"`
	IsBackupEnabled bool            `json:"isBackupEnabled"`
}

// AIProviderPreset defines metadata and default configuration for an AI model provider.
type AIProviderPreset struct {
	Provider            IDEProviderType `json:"provider"`
	DisplayName         string          `json:"displayName"`
	DefaultBaseURL      string          `json:"defaultBaseUrl"`
	DefaultModel        string          `json:"defaultModel"`
	FallbackModels      []string        `json:"fallbackModels"`
	ContextLimit        int             `json:"contextLimit"`
	AuthHeaderScheme    string          `json:"authHeaderScheme"`
	EnvVarKeys          []string        `json:"envVarKeys"`
	HasStreamingSupport bool            `json:"hasStreamingSupport"`
}

// ProviderPresetConfig aliases AIProviderPreset for backward compatibility.
type ProviderPresetConfig = AIProviderPreset

// TargetConfigPaths details the resolved file paths for all supported IDE targets.
type TargetConfigPaths struct {
	ClineSettingsPath string `json:"clineSettingsPath"`
	ClineStorageDir   string `json:"clineStorageDir"`
	ContinuePath      string `json:"continuePath"`
	ZedSettingsPath   string `json:"zedSettingsPath"`
}

// IDETargetResult records the configuration outcome for a single editor target.
type IDETargetResult struct {
	Target          IDETargetType `json:"target"`
	ConfigPath      string        `json:"configPath"`
	IsConfigured    bool          `json:"isConfigured"`
	IsBackupCreated bool          `json:"isBackupCreated"`
	BackupPath      string        `json:"backupPath,omitempty"`
	ErrorMessage    string        `json:"errorMessage,omitempty"`
}

// ConnectionTestRequest specifies parameters to test endpoint reachability.
type ConnectionTestRequest struct {
	Provider IDEProviderType `json:"provider"`
	BaseURL  string          `json:"baseUrl"`
	APIKey   string          `json:"apiKey"`
	Model    string          `json:"model"`
}

// TestConnectionRequest aliases ConnectionTestRequest.
type TestConnectionRequest = ConnectionTestRequest

// ConnectionTestResult details the outcome of an endpoint connection test.
type ConnectionTestResult struct {
	IsSuccess      bool   `json:"isSuccess"`
	LatencyMs      int64  `json:"latencyMs"`
	Model          string `json:"model"`
	HTTPStatusCode int    `json:"httpStatusCode"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}

// TestConnectionResult aliases ConnectionTestResult.
type TestConnectionResult = ConnectionTestResult

// IDESetupResult aggregates configuration outcomes across selected targets.
type IDESetupResult struct {
	Target         IDETargetType              `json:"target"`
	Provider       IDEProviderType            `json:"provider"`
	IsSuccess      bool                       `json:"isSuccess"`
	TargetResults  map[string]IDETargetResult `json:"targetResults"`
	ConnectionTest ConnectionTestResult       `json:"connectionTest,omitempty"`
	BackupPaths    []string                   `json:"backupPaths,omitempty"`
	SummaryMessage string                     `json:"summaryMessage"`
}

// BackupResult details the outcome of a pre-flight backup operation.
type BackupResult struct {
	OriginalPath    string `json:"originalPath"`
	BackupPath      string `json:"backupPath,omitempty"`
	IsBackupCreated bool   `json:"isBackupCreated"`
}
