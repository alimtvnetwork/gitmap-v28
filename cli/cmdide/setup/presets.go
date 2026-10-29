package setup

import (
	"os"
	"strings"
)

// GetProviderPresets returns all production preset definitions.
func GetProviderPresets() map[IDEProviderType]AIProviderPreset {
	return map[IDEProviderType]AIProviderPreset{
		ProviderKimi:       kimiPreset(),
		ProviderOpenRouter: openRouterPreset(),
		ProviderDeepSeek:   deepSeekPreset(),
		ProviderAnthropic:  anthropicPreset(),
		ProviderOllama:     ollamaPreset(),
		ProviderCustom:     customPreset(),
	}
}

func kimiPreset() AIProviderPreset {
	return AIProviderPreset{
		Provider:            ProviderKimi,
		DisplayName:         "Moonshot Kimi K3.1 Pro (Recommended)",
		DefaultBaseURL:      "https://api.moonshot.cn/v1",
		DefaultModel:        "kimi-k3.1-pro",
		FallbackModels:      []string{"moonshot-v1-128k", "moonshot-v1-32k"},
		ContextLimit:        128000,
		AuthHeaderScheme:    "Bearer",
		EnvVarKeys:          []string{"KIMI_API_KEY", "MOONSHOT_API_KEY"},
		HasStreamingSupport: true,
	}
}

func openRouterPreset() AIProviderPreset {
	return AIProviderPreset{
		Provider:            ProviderOpenRouter,
		DisplayName:         "OpenRouter Multi-Model",
		DefaultBaseURL:      "https://openrouter.ai/api/v1",
		DefaultModel:        "moonshotai/kimi-k3.1-pro",
		FallbackModels:      []string{"anthropic/claude-3.7-sonnet", "deepseek/deepseek-chat"},
		ContextLimit:        128000,
		AuthHeaderScheme:    "Bearer",
		EnvVarKeys:          []string{"OPENROUTER_API_KEY"},
		HasStreamingSupport: true,
	}
}

func deepSeekPreset() AIProviderPreset {
	return AIProviderPreset{
		Provider:            ProviderDeepSeek,
		DisplayName:         "DeepSeek AI",
		DefaultBaseURL:      "https://api.deepseek.com/v1",
		DefaultModel:        "deepseek-chat",
		FallbackModels:      []string{"deepseek-reasoner"},
		ContextLimit:        64000,
		AuthHeaderScheme:    "Bearer",
		EnvVarKeys:          []string{"DEEPSEEK_API_KEY"},
		HasStreamingSupport: true,
	}
}

func anthropicPreset() AIProviderPreset {
	return AIProviderPreset{
		Provider:            ProviderAnthropic,
		DisplayName:         "Anthropic Claude",
		DefaultBaseURL:      "https://api.anthropic.com/v1",
		DefaultModel:        "claude-3-7-sonnet-20250219",
		FallbackModels:      []string{"claude-3-5-sonnet-latest", "claude-3-5-haiku-latest"},
		ContextLimit:        200000,
		AuthHeaderScheme:    "x-api-key",
		EnvVarKeys:          []string{"ANTHROPIC_API_KEY"},
		HasStreamingSupport: true,
	}
}

func ollamaPreset() AIProviderPreset {
	return AIProviderPreset{
		Provider:            ProviderOllama,
		DisplayName:         "Local Ollama Server",
		DefaultBaseURL:      "http://localhost:11434/v1",
		DefaultModel:        "deepseek-r1",
		FallbackModels:      []string{"qwen2.5-coder:32b", "llama3.3:70b"},
		ContextLimit:        32000,
		AuthHeaderScheme:    "Bearer",
		EnvVarKeys:          []string{"OLLAMA_API_KEY"},
		HasStreamingSupport: true,
	}
}

func customPreset() AIProviderPreset {
	return AIProviderPreset{
		Provider:            ProviderCustom,
		DisplayName:         "Custom Endpoint",
		DefaultBaseURL:      "",
		DefaultModel:        "",
		FallbackModels:      []string{},
		ContextLimit:        32000,
		AuthHeaderScheme:    "Bearer",
		EnvVarKeys:          []string{"OPENAI_API_KEY"},
		HasStreamingSupport: true,
	}
}

// GetPreset returns a preset by provider type if registered.
func GetPreset(provider IDEProviderType) (AIProviderPreset, bool) {
	preset, isFound := GetProviderPresets()[provider]

	return preset, isFound
}

// ResolveAPIKey determines the active API key through environment cascades.
func ResolveAPIKey(provider IDEProviderType, explicitKey string) string {
	if strings.TrimSpace(explicitKey) != "" {
		return strings.TrimSpace(explicitKey)
	}

	if key := resolvePresetEnvKey(provider); key != "" {
		return key
	}

	if val := os.Getenv("GITMAP_AI_API_KEY"); val != "" {
		return val
	}

	return os.Getenv("OPENAI_API_KEY")
}

func resolvePresetEnvKey(provider IDEProviderType) string {
	preset, isFound := GetPreset(provider)
	if !isFound {
		return ""
	}

	for _, envKey := range preset.EnvVarKeys {
		if val := os.Getenv(envKey); val != "" {
			return val
		}
	}

	return ""
}

// ResolveBaseURL resolves the effective base URL for a provider.
func ResolveBaseURL(provider IDEProviderType, explicitURL string) string {
	if strings.TrimSpace(explicitURL) != "" {
		return strings.TrimSpace(explicitURL)
	}

	preset, isFound := GetPreset(provider)
	if isFound {
		return preset.DefaultBaseURL
	}

	return ""
}

// ResolveModel resolves the effective model identifier for a provider.
func ResolveModel(provider IDEProviderType, explicitModel string) string {
	if strings.TrimSpace(explicitModel) != "" {
		return strings.TrimSpace(explicitModel)
	}

	preset, isFound := GetPreset(provider)
	if isFound {
		return preset.DefaultModel
	}

	return ""
}

// ApplyPresetDefaults enriches options with preset defaults when unset.
func ApplyPresetDefaults(opts *IDESetupOptions) {
	if opts == nil {
		return
	}

	if opts.Provider == "" {
		opts.Provider = ProviderKimi
	}

	applyEndpointDefaults(opts)
	applyAuthDefaults(opts)
}

func applyEndpointDefaults(opts *IDESetupOptions) {
	if opts.BaseURL == "" {
		opts.BaseURL = ResolveBaseURL(opts.Provider, "")
	}

	if opts.Model == "" {
		opts.Model = ResolveModel(opts.Provider, "")
	}
}

func applyAuthDefaults(opts *IDESetupOptions) {
	if opts.APIKey == "" {
		opts.APIKey = ResolveAPIKey(opts.Provider, "")
	}
}
