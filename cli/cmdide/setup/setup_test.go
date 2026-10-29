package setup

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIDETargetTypeValidation(t *testing.T) {
	for _, target := range AllIDETargetTypes() {
		assertTargetValid(t, target, true)
	}

	assertTargetValid(t, IDETargetType("invalid"), false)
}

func assertTargetValid(t *testing.T, target IDETargetType, expected bool) {
	t.Helper()
	if target.IsValid() != expected {
		t.Errorf("target %s validity expected %v", target, expected)
	}
}

func TestIDETargetTypeParsing(t *testing.T) {
	verifyTargetParse(t, "cline", TargetCline)
	verifyTargetParse(t, "composer", TargetCline)
	verifyTargetParse(t, "continue", TargetContinue)
	verifyTargetParse(t, "tab", TargetContinue)
	verifyTargetParse(t, "zed", TargetZed)
	verifyTargetParse(t, "all", TargetAll)
}

func verifyTargetParse(t *testing.T, input string, expected IDETargetType) {
	t.Helper()
	parsed, err := ParseIDETargetType(input)
	if err != nil || parsed != expected {
		t.Errorf("input %s expected %s, got %s", input, expected, parsed)
	}
}

func TestIDEProviderTypeParsing(t *testing.T) {
	verifyProviderParse(t, "kimi", ProviderKimi)
	verifyProviderParse(t, "moonshot", ProviderKimi)
	verifyProviderParse(t, "openrouter", ProviderOpenRouter)
	verifyProviderParse(t, "deepseek", ProviderDeepSeek)
	verifyProviderParse(t, "anthropic", ProviderAnthropic)
	verifyProviderParse(t, "ollama", ProviderOllama)
	verifyProviderParse(t, "custom", ProviderCustom)
}

func verifyProviderParse(t *testing.T, input string, expected IDEProviderType) {
	t.Helper()
	parsed, err := ParseIDEProviderType(input)
	if err != nil || parsed != expected {
		t.Errorf("input %s expected %s, got %s", input, expected, parsed)
	}
}

func TestPresetsDefaults(t *testing.T) {
	presets := GetProviderPresets()
	kimi, exists := presets[ProviderKimi]
	if !exists {
		t.Fatalf("kimi preset missing")
	}

	assertKimiPresetFields(t, kimi)
}

func assertKimiPresetFields(t *testing.T, kimi AIProviderPreset) {
	t.Helper()
	if kimi.DefaultBaseURL != "https://api.moonshot.cn/v1" {
		t.Errorf("unexpected kimi base url: %s", kimi.DefaultBaseURL)
	}

	if kimi.DefaultModel != "kimi-k3.1-pro" {
		t.Errorf("unexpected kimi model: %s", kimi.DefaultModel)
	}

	if kimi.ContextLimit != 128000 {
		t.Errorf("unexpected kimi context limit: %d", kimi.ContextLimit)
	}
}

func TestResolveAPIKey(t *testing.T) {
	if explicit := ResolveAPIKey(ProviderKimi, "explicit-key"); explicit != "explicit-key" {
		t.Errorf("expected explicit key, got %s", explicit)
	}

	t.Setenv("KIMI_API_KEY", "env-kimi-key")
	if fromEnv := ResolveAPIKey(ProviderKimi, ""); fromEnv != "env-kimi-key" {
		t.Errorf("expected env key, got %s", fromEnv)
	}
}

func TestApplyPresetDefaults(t *testing.T) {
	opts := IDESetupOptions{Provider: ProviderKimi}
	ApplyPresetDefaults(&opts)

	if opts.BaseURL != "https://api.moonshot.cn/v1" {
		t.Errorf("expected kimi base url, got %s", opts.BaseURL)
	}

	if opts.Model != "kimi-k3.1-pro" {
		t.Errorf("expected kimi model, got %s", opts.Model)
	}
}

func TestPathResolverOS(t *testing.T) {
	resolver := &PathResolver{
		BaseHomeDir:    filepath.Join("C:", "Users", "TestUser"),
		BaseAppDataDir: filepath.Join("C:", "Users", "TestUser", "AppData", "Roaming"),
		TargetOS:       "windows",
	}

	paths, err := resolver.GetAllTargetPaths()
	if err != nil {
		t.Fatalf("GetAllTargetPaths failed: %v", err)
	}

	assertPathsNonEmpty(t, paths)
}

func assertPathsNonEmpty(t *testing.T, paths TargetConfigPaths) {
	t.Helper()
	if paths.ClineSettingsPath == "" || paths.ContinuePath == "" || paths.ZedSettingsPath == "" {
		t.Errorf("all paths should be non-empty: %+v", paths)
	}
}

func TestPathResolverUnix(t *testing.T) {
	resolver := &PathResolver{BaseHomeDir: "/home/testuser", TargetOS: "linux"}
	zedPath, err := resolver.GetZedSettingsPath()
	if err != nil {
		t.Fatalf("GetZedSettingsPath failed: %v", err)
	}

	expected := filepath.Join("/home/testuser", ".config", "zed", "settings.json")
	if zedPath != expected {
		t.Errorf("expected %s, got %s", expected, zedPath)
	}
}

func TestPreFlightBackup(t *testing.T) {
	targetFile := filepath.Join(t.TempDir(), "config.json")
	_ = os.WriteFile(targetFile, []byte(`{"key":"value"}`), 0644)

	res, err := CreatePreFlightBackup(targetFile)
	if err != nil || !res.IsBackupCreated {
		t.Fatalf("backup creation failed: %v", err)
	}

	assertBackupFileContents(t, res.BackupPath, `{"key":"value"}`)
}

func assertBackupFileContents(t *testing.T, path, expected string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil || string(data) != expected {
		t.Errorf("backup data mismatch")
	}
}

func TestDeepMergeMaps(t *testing.T) {
	dest := map[string]interface{}{"existing": "original", "nested": map[string]interface{}{"a": 1}}
	src := map[string]interface{}{"added": "new", "nested": map[string]interface{}{"b": 2}}

	DeepMergeMaps(dest, src)
	assertDeepMergeResults(t, dest)
}

func assertDeepMergeResults(t *testing.T, dest map[string]interface{}) {
	t.Helper()
	if dest["existing"] != "original" || dest["added"] != "new" {
		t.Errorf("top-level merge mismatch: %+v", dest)
	}

	nested := dest["nested"].(map[string]interface{})
	if nested["a"] != 1 || nested["b"] != 2 {
		t.Errorf("nested merge mismatch: %+v", nested)
	}
}

func TestGenerateClineConfig(t *testing.T) {
	tempDir := t.TempDir()
	resolver := &PathResolver{BaseHomeDir: tempDir, BaseAppDataDir: tempDir, TargetOS: "windows"}
	opts := sampleSetupOptions(TargetCline, ProviderKimi)

	if err := GenerateClineConfig(resolver, opts); err != nil {
		t.Fatalf("GenerateClineConfig failed: %v", err)
	}

	verifyClineOutput(t, resolver)
}

func sampleSetupOptions(target IDETargetType, provider IDEProviderType) IDESetupOptions {
	return IDESetupOptions{
		Target:   target,
		Provider: provider,
		BaseURL:  "https://api.moonshot.cn/v1",
		Model:    "kimi-k3.1-pro",
		APIKey:   "test-key",
	}
}

func verifyClineOutput(t *testing.T, resolver *PathResolver) {
	t.Helper()
	settingsPath, _ := resolver.GetClineSettingsPath()
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("failed reading cline settings: %v", err)
	}

	var parsed map[string]interface{}
	_ = json.Unmarshal(data, &parsed)
	mcpServers, hasServers := parsed["mcpServers"].(map[string]interface{})
	if !hasServers || mcpServers["gitmap"] == nil {
		t.Errorf("cline settings missing gitmap mcp server: %s", string(data))
	}
}

func TestGenerateContinueConfig(t *testing.T) {
	tempDir := t.TempDir()
	resolver := &PathResolver{BaseHomeDir: tempDir, TargetOS: "linux"}
	opts := sampleSetupOptions(TargetContinue, ProviderKimi)

	if err := GenerateContinueConfig(resolver, opts); err != nil {
		t.Fatalf("GenerateContinueConfig failed: %v", err)
	}

	verifyContinueOutput(t, resolver)
}

func verifyContinueOutput(t *testing.T, resolver *PathResolver) {
	t.Helper()
	configPath, _ := resolver.GetContinueConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed reading continue config: %v", err)
	}

	var parsed map[string]interface{}
	_ = json.Unmarshal(data, &parsed)
	models, hasModels := parsed["models"].([]interface{})
	if !hasModels || len(models) == 0 {
		t.Fatalf("continue config missing models: %s", string(data))
	}
}

func TestGenerateZedConfig(t *testing.T) {
	tempDir := t.TempDir()
	resolver := &PathResolver{BaseHomeDir: tempDir, BaseAppDataDir: tempDir, TargetOS: "windows"}
	opts := sampleSetupOptions(TargetZed, ProviderKimi)

	if err := GenerateZedConfig(resolver, opts); err != nil {
		t.Fatalf("GenerateZedConfig failed: %v", err)
	}

	verifyZedOutput(t, resolver)
}

func verifyZedOutput(t *testing.T, resolver *PathResolver) {
	t.Helper()
	settingsPath, _ := resolver.GetZedSettingsPath()
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("failed reading zed settings: %v", err)
	}

	var parsed map[string]interface{}
	_ = json.Unmarshal(data, &parsed)
	if parsed["assistant"] == nil || parsed["language_models"] == nil {
		t.Errorf("zed settings missing assistant or language_models")
	}
}

func TestTestConnectionSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	req := ConnectionTestRequest{Provider: ProviderKimi, BaseURL: server.URL, APIKey: "sk-test", Model: "kimi-k3.1-pro"}
	res := TestConnection(req)
	if !res.IsSuccess || res.HTTPStatusCode != http.StatusOK {
		t.Errorf("expected successful connection test: %+v", res)
	}
}

func TestTestConnectionFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	req := ConnectionTestRequest{Provider: ProviderKimi, BaseURL: server.URL, APIKey: "bad-key", Model: "kimi-k3.1-pro"}
	res := TestConnection(req)
	if res.IsSuccess || res.HTTPStatusCode != http.StatusUnauthorized {
		t.Errorf("expected failed connection test: %+v", res)
	}
}

func TestDryRunMode(t *testing.T) {
	resolver := &PathResolver{BaseHomeDir: t.TempDir(), TargetOS: "linux"}
	opts := IDESetupOptions{Target: TargetAll, Provider: ProviderKimi, IsDryRun: true}

	if err := GenerateClineConfig(resolver, opts); err != nil {
		t.Fatalf("dry run cline failed: %v", err)
	}

	path, _ := resolver.GetClineSettingsPath()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("file should not exist in dry run: %s", path)
	}
}

func TestAnthropicProviderGenerators(t *testing.T) {
	tempDir := t.TempDir()
	resolver := &PathResolver{BaseHomeDir: tempDir, BaseAppDataDir: tempDir, TargetOS: "windows"}
	opts := sampleSetupOptions(TargetAll, ProviderAnthropic)

	verifyAnthropicGenerators(t, resolver, opts)
}

func verifyAnthropicGenerators(t *testing.T, resolver *PathResolver, opts IDESetupOptions) {
	t.Helper()
	if err := GenerateClineConfig(resolver, opts); err != nil {
		t.Fatalf("cline anthropic failed: %v", err)
	}

	if err := GenerateContinueConfig(resolver, opts); err != nil {
		t.Fatalf("continue anthropic failed: %v", err)
	}

	if err := GenerateZedConfig(resolver, opts); err != nil {
		t.Fatalf("zed anthropic failed: %v", err)
	}
}

func TestRunSetupDryRun(t *testing.T) {
	opts := IDESetupOptions{Target: TargetAll, Provider: ProviderKimi, IsDryRun: true}
	res, err := RunSetup(opts)
	if err != nil || !res.IsSuccess {
		t.Fatalf("RunSetup dry run failed: %v", err)
	}

	if len(res.TargetResults) != 3 {
		t.Errorf("expected 3 target results, got %d", len(res.TargetResults))
	}
}
