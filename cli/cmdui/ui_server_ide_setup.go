package cmdui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// EditorStatus encapsulates detection data for an IDE.
type EditorStatus struct {
	EditorID    string            `json:"editorId"`
	DisplayName string            `json:"displayName"`
	IsInstalled bool              `json:"isInstalled"`
	BinaryPath  string            `json:"binaryPath"`
	ConfigPath  string            `json:"configPath"`
	Extensions  []ExtensionStatus `json:"extensions"`
}

// ExtensionStatus details configuration of an AI extension.
type ExtensionStatus struct {
	ExtensionID    string `json:"extensionId"`
	DisplayName    string `json:"displayName"`
	IsInstalled    bool   `json:"isInstalled"`
	IsConfigured   bool   `json:"isConfigured"`
	ActiveProvider string `json:"activeProvider"`
	ActiveModel    string `json:"activeModel"`
	ConfigPath     string `json:"configPath"`
}

// IDEItemInfo models brief status info for map-based status consumers.
type IDEItemInfo struct {
	IsInstalled    bool   `json:"isInstalled"`
	ConfigPath     string `json:"configPath"`
	HasConfig      bool   `json:"hasConfig"`
	ActiveProvider string `json:"activeProvider"`
	ActiveModel    string `json:"activeModel,omitempty"`
}

// IDESetupStatusResponse is the envelope returned by GET /api/ide/setup/status.
type IDESetupStatusResponse struct {
	IsSuccess       bool                   `json:"isSuccess"`
	DetectedEditors []EditorStatus         `json:"detectedEditors"`
	Data            map[string]IDEItemInfo `json:"data,omitempty"`
	ErrorMessage    string                 `json:"errorMessage,omitempty"`
}

// ProviderPreset defines known API endpoints, models, and defaults.
type ProviderPreset struct {
	PresetID        string   `json:"presetId"`
	DisplayName     string   `json:"displayName"`
	DefaultBaseURL  string   `json:"defaultBaseUrl"`
	DefaultModel    string   `json:"defaultModel"`
	SupportedModels []string `json:"supportedModels"`
	RequiresAPIKey  bool     `json:"requiresApiKey"`
	HelpURL         string   `json:"helpUrl"`
}

// IDESetupPresetsResponse is the envelope returned by GET /api/ide/setup/presets.
type IDESetupPresetsResponse struct {
	IsSuccess    bool             `json:"isSuccess"`
	Presets      []ProviderPreset `json:"presets"`
	ErrorMessage string           `json:"errorMessage,omitempty"`
}

// IDETestConnectionRequest defines the payload for testing an endpoint.
type IDETestConnectionRequest struct {
	Provider  string `json:"provider"`
	PresetID  string `json:"presetId"`
	BaseURL   string `json:"baseUrl"`
	APIKey    string `json:"apiKey"`
	Model     string `json:"model"`
	ModelName string `json:"modelName"`
}

// IDETestConnectionData provides detailed latency and probe diagnostics.
type IDETestConnectionData struct {
	IsSuccess      bool   `json:"isSuccess"`
	LatencyMs      int64  `json:"latencyMs"`
	Model          string `json:"model,omitempty"`
	HTTPStatusCode int    `json:"httpStatusCode"`
	ModelEcho      string `json:"modelEcho,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	ErrorCode      string `json:"errorCode,omitempty"`
}

// IDETestConnectionResponse is the envelope returned by POST /api/ide/setup/test-connection.
type IDETestConnectionResponse struct {
	IsSuccess      bool                   `json:"isSuccess"`
	HTTPStatusCode int                    `json:"httpStatusCode,omitempty"`
	LatencyMs      int64                  `json:"latencyMs,omitempty"`
	ModelEcho      string                 `json:"modelEcho,omitempty"`
	ErrorMessage   string                 `json:"errorMessage,omitempty"`
	Data           *IDETestConnectionData `json:"data,omitempty"`
}

// IDESetupApplyRequest specifies the configuration target and credentials.
type IDESetupApplyRequest struct {
	Target          string `json:"target"`
	TargetEditor    string `json:"targetEditor"`
	TargetExtension string `json:"targetExtension"`
	Provider        string `json:"provider"`
	PresetID        string `json:"presetId"`
	BaseURL         string `json:"baseUrl"`
	APIKey          string `json:"apiKey"`
	Model           string `json:"model"`
	ModelName       string `json:"modelName"`
	IsDryRun        bool   `json:"isDryRun"`
	IsBackupEnabled bool   `json:"isBackupEnabled"`
}

// IDESetupApplyData holds outcomes of applying setup configurations.
type IDESetupApplyData struct {
	IsSuccess         bool     `json:"isSuccess"`
	ConfiguredTargets []string `json:"configuredTargets"`
	BackupPaths       []string `json:"backupPaths"`
	IsDryRun          bool     `json:"isDryRun"`
	DiffSummary       string   `json:"diffSummary,omitempty"`
}

// IDESetupApplyResponse is the envelope returned by POST /api/ide/setup/apply.
type IDESetupApplyResponse struct {
	IsSuccess         bool               `json:"isSuccess"`
	TargetFile        string             `json:"targetFile,omitempty"`
	BackupFile        string             `json:"backupFile,omitempty"`
	IsDryRun          bool               `json:"isDryRun"`
	DiffSummary       string             `json:"diffSummary,omitempty"`
	ConfiguredTargets []string           `json:"configuredTargets,omitempty"`
	BackupPaths       []string           `json:"backupPaths,omitempty"`
	ErrorMessage      string             `json:"errorMessage,omitempty"`
	Data              *IDESetupApplyData `json:"data,omitempty"`
}

var (
	ideConfigOverridesMu      sync.RWMutex
	customConfigPathOverrides = make(map[string]string)
	ideHTTPClient             = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func setTargetConfigOverride(target, path string) {
	ideConfigOverridesMu.Lock()
	defer ideConfigOverridesMu.Unlock()

	customConfigPathOverrides[target] = path
}

func clearTargetConfigOverrides() {
	ideConfigOverridesMu.Lock()
	defer ideConfigOverridesMu.Unlock()

	customConfigPathOverrides = make(map[string]string)
}

func getTargetConfigPath(target string) string {
	ideConfigOverridesMu.RLock()
	override, hasOverride := customConfigPathOverrides[target]
	ideConfigOverridesMu.RUnlock()

	if hasOverride && override != "" {
		return override
	}

	return resolveDefaultConfigPath(target)
}

func handleIDESetupStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)

		return
	}

	resp := collectIDESetupStatus()

	_ = json.NewEncoder(w).Encode(resp)
}

func collectIDESetupStatus() IDESetupStatusResponse {
	editors := detectHostEditors()

	return IDESetupStatusResponse{
		IsSuccess:       true,
		DetectedEditors: editors,
		Data:            buildStatusDataMap(editors),
	}
}

func detectHostEditors() []EditorStatus {
	return []EditorStatus{
		detectVSCode(),
		detectCursor(),
		detectZed(),
	}
}

func detectVSCode() EditorStatus {
	home, _ := os.UserHomeDir()
	appData := os.Getenv("APPDATA")
	binPath := findExecutable("code", getVSCodeFallbackPaths())
	cfgPath := resolveVSCodeConfigPath(home, appData)

	return EditorStatus{
		EditorID:    "vscode",
		DisplayName: "Visual Studio Code",
		IsInstalled: binPath != "" || isDirExisting(filepath.Dir(cfgPath)),
		BinaryPath:  binPath,
		ConfigPath:  cfgPath,
		Extensions:  buildVSCodeExtensions(home, appData),
	}
}

func getVSCodeFallbackPaths() []string {
	return []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Microsoft VS Code", "Code.exe"),
		"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
		"/usr/bin/code",
	}
}

func resolveVSCodeConfigPath(home, appData string) string {
	if runtime.GOOS == "windows" && appData != "" {
		return filepath.Join(appData, "Code", "User", "settings.json")
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")
	}

	return filepath.Join(home, ".config", "Code", "User", "settings.json")
}

func buildVSCodeExtensions(home, appData string) []ExtensionStatus {
	return []ExtensionStatus{
		detectClineExtension(home, appData),
		detectContinueExtension(home),
	}
}

func detectClineExtension(home, appData string) ExtensionStatus {
	cfgPath := resolveClineConfigPath(home, appData)
	hasCfg := isFileExisting(cfgPath)

	return ExtensionStatus{
		ExtensionID:  "cline",
		DisplayName:  "Cline (Cursor Composer)",
		IsInstalled:  hasCfg || isDirExisting(filepath.Dir(cfgPath)),
		IsConfigured: hasCfg,
		ConfigPath:   cfgPath,
	}
}

func resolveClineConfigPath(home, appData string) string {
	if runtime.GOOS == "windows" && appData != "" {
		return filepath.Join(appData, "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")
	}

	return filepath.Join(home, ".config", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")
}

func detectContinueExtension(home string) ExtensionStatus {
	cfgPath := filepath.Join(home, ".continue", "config.json")
	hasCfg := isFileExisting(cfgPath)

	return ExtensionStatus{
		ExtensionID:  "continue",
		DisplayName:  "Continue.dev (Cursor Tab & Autocomplete)",
		IsInstalled:  hasCfg || isDirExisting(filepath.Dir(cfgPath)),
		IsConfigured: hasCfg,
		ConfigPath:   cfgPath,
	}
}

func detectCursor() EditorStatus {
	home, _ := os.UserHomeDir()
	appData := os.Getenv("APPDATA")
	binPath := findExecutable("cursor", getCursorFallbackPaths())
	cfgPath := resolveCursorConfigPath(home, appData)

	return EditorStatus{
		EditorID:    "cursor",
		DisplayName: "Cursor",
		IsInstalled: binPath != "" || isDirExisting(filepath.Dir(cfgPath)),
		BinaryPath:  binPath,
		ConfigPath:  cfgPath,
		Extensions:  []ExtensionStatus{},
	}
}

func getCursorFallbackPaths() []string {
	return []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "cursor", "Cursor.exe"),
		"/Applications/Cursor.app/Contents/Resources/app/bin/cursor",
	}
}

func resolveCursorConfigPath(home, appData string) string {
	if runtime.GOOS == "windows" && appData != "" {
		return filepath.Join(appData, "Cursor", "User", "settings.json")
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "Cursor", "User", "settings.json")
	}

	return filepath.Join(home, ".config", "Cursor", "User", "settings.json")
}

func detectZed() EditorStatus {
	home, _ := os.UserHomeDir()
	appData := os.Getenv("APPDATA")
	binPath := findExecutable("zed", getZedFallbackPaths(home))
	cfgPath := resolveZedConfigPath(home, appData)

	return EditorStatus{
		EditorID:    "zed",
		DisplayName: "Zed",
		IsInstalled: binPath != "" || isDirExisting(filepath.Dir(cfgPath)),
		BinaryPath:  binPath,
		ConfigPath:  cfgPath,
		Extensions:  buildZedExtensions(cfgPath),
	}
}

func getZedFallbackPaths(home string) []string {
	return []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Zed", "zed.exe"),
		filepath.Join(home, ".local", "bin", "zed"),
		"/usr/local/bin/zed",
	}
}

func resolveZedConfigPath(home, appData string) string {
	if runtime.GOOS == "windows" && appData != "" {
		return filepath.Join(appData, "Zed", "settings.json")
	}

	return filepath.Join(home, ".config", "zed", "settings.json")
}

func buildZedExtensions(cfgPath string) []ExtensionStatus {
	hasCfg := isFileExisting(cfgPath)

	return []ExtensionStatus{
		{
			ExtensionID:  "zed-native",
			DisplayName:  "Zed AI Native Assistant",
			IsInstalled:  hasCfg || isDirExisting(filepath.Dir(cfgPath)),
			IsConfigured: hasCfg,
			ConfigPath:   cfgPath,
		},
	}
}

func findExecutable(name string, fallbacks []string) string {
	path, err := exec.LookPath(name)
	if err == nil && path != "" {
		return path
	}

	for _, fb := range fallbacks {
		if isFileExisting(fb) {
			return fb
		}
	}

	return ""
}

func isFileExisting(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func isDirExisting(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func buildStatusDataMap(editors []EditorStatus) map[string]IDEItemInfo {
	data := make(map[string]IDEItemInfo)
	for _, ed := range editors {
		populateEditorDataMap(data, ed)
	}

	return data
}

func populateEditorDataMap(data map[string]IDEItemInfo, ed EditorStatus) {
	data[ed.EditorID] = IDEItemInfo{
		IsInstalled: ed.IsInstalled,
		ConfigPath:  ed.ConfigPath,
		HasConfig:   isFileExisting(ed.ConfigPath),
	}
	for _, ext := range ed.Extensions {
		data[ext.ExtensionID] = IDEItemInfo{
			IsInstalled:    ext.IsInstalled,
			ConfigPath:     ext.ConfigPath,
			HasConfig:      ext.IsConfigured,
			ActiveProvider: ext.ActiveProvider,
			ActiveModel:    ext.ActiveModel,
		}
	}
}

func handleIDESetupPresets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)

		return
	}

	resp := IDESetupPresetsResponse{
		IsSuccess: true,
		Presets:   getDefaultProviderPresets(),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func getDefaultProviderPresets() []ProviderPreset {
	return []ProviderPreset{
		buildKimiPreset(),
		buildOpenRouterPreset(),
		buildDeepSeekPreset(),
		buildAnthropicPreset(),
		buildOllamaPreset(),
		buildCustomPreset(),
	}
}

func buildKimiPreset() ProviderPreset {
	return ProviderPreset{
		PresetID:        "kimi",
		DisplayName:     "Kimi K3.1 Pro (Moonshot AI)",
		DefaultBaseURL:  "https://api.moonshot.cn/v1",
		DefaultModel:    "kimi-k3.1-pro",
		SupportedModels: []string{"kimi-k3.1-pro", "moonshot-v1-32k", "moonshot-v1-128k"},
		RequiresAPIKey:  true,
		HelpURL:         "https://platform.moonshot.cn/docs",
	}
}

func buildOpenRouterPreset() ProviderPreset {
	return ProviderPreset{
		PresetID:        "openrouter",
		DisplayName:     "OpenRouter Universal",
		DefaultBaseURL:  "https://openrouter.ai/api/v1",
		DefaultModel:    "anthropic/claude-3.5-sonnet",
		SupportedModels: []string{"anthropic/claude-3.5-sonnet", "deepseek/deepseek-coder"},
		RequiresAPIKey:  true,
		HelpURL:         "https://openrouter.ai/keys",
	}
}

func buildDeepSeekPreset() ProviderPreset {
	return ProviderPreset{
		PresetID:        "deepseek",
		DisplayName:     "DeepSeek Coder / V3",
		DefaultBaseURL:  "https://api.deepseek.com/v1",
		DefaultModel:    "deepseek-coder",
		SupportedModels: []string{"deepseek-coder", "deepseek-chat", "deepseek-reasoner"},
		RequiresAPIKey:  true,
		HelpURL:         "https://platform.deepseek.com/api_keys",
	}
}

func buildAnthropicPreset() ProviderPreset {
	return ProviderPreset{
		PresetID:        "anthropic",
		DisplayName:     "Anthropic Claude Official",
		DefaultBaseURL:  "https://api.anthropic.com/v1",
		DefaultModel:    "claude-3-5-sonnet-20241022",
		SupportedModels: []string{"claude-3-5-sonnet-20241022", "claude-3-5-haiku-20241022"},
		RequiresAPIKey:  true,
		HelpURL:         "https://console.anthropic.com/settings/keys",
	}
}

func buildOllamaPreset() ProviderPreset {
	return ProviderPreset{
		PresetID:        "ollama",
		DisplayName:     "Ollama (Local Offline)",
		DefaultBaseURL:  "http://localhost:11434/v1",
		DefaultModel:    "qwen2.5-coder:7b",
		SupportedModels: []string{"qwen2.5-coder:7b", "deepseek-coder-v2:16b", "llama3.1:8b"},
		RequiresAPIKey:  false,
		HelpURL:         "https://ollama.com/library",
	}
}

func buildCustomPreset() ProviderPreset {
	return ProviderPreset{
		PresetID:        "custom",
		DisplayName:     "Custom OpenAI-Compatible API",
		DefaultBaseURL:  "",
		DefaultModel:    "",
		SupportedModels: []string{},
		RequiresAPIKey:  true,
		HelpURL:         "",
	}
}

func handleIDESetupTestConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)

		return
	}

	var req IDETestConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeTestConnectionError(w, http.StatusBadRequest, "invalid request body: "+err.Error())

		return
	}

	normalizeTestRequest(&req)
	resp := executeTestConnection(req)

	_ = json.NewEncoder(w).Encode(resp)
}

func normalizeTestRequest(req *IDETestConnectionRequest) {
	if req.Provider == "" {
		req.Provider = req.PresetID
	}
	if req.Model == "" {
		req.Model = req.ModelName
	}
	if req.BaseURL == "" {
		req.BaseURL = resolvePresetBaseURL(req.Provider)
	}
}

func resolvePresetBaseURL(provider string) string {
	for _, p := range getDefaultProviderPresets() {
		if strings.EqualFold(p.PresetID, provider) {
			return p.DefaultBaseURL
		}
	}

	return ""
}

func validateTestConnectionReq(req IDETestConnectionRequest) error {
	if req.BaseURL == "" {
		return apperror.NewSimple("test_connection", "baseURL is required")
	}
	if req.APIKey == "" && !strings.EqualFold(req.Provider, "ollama") {
		return apperror.NewSimple("test_connection", "apiKey is required for provider "+req.Provider)
	}

	return nil
}

func executeTestConnection(req IDETestConnectionRequest) IDETestConnectionResponse {
	valErr := validateTestConnectionReq(req)
	if valErr != nil {
		return buildTestFailureResponse(http.StatusBadRequest, 0, valErr.Error(), "VALIDATION_ERROR")
	}

	probeURL := strings.TrimRight(req.BaseURL, "/") + "/models"
	httpReq, err := http.NewRequest(http.MethodGet, probeURL, nil)
	if err != nil {
		return buildTestFailureResponse(http.StatusBadRequest, 0, err.Error(), "INVALID_REQUEST")
	}

	applyAuthHeaders(httpReq, req.Provider, req.APIKey)

	return performProbe(httpReq, req.Model)
}

func performProbe(httpReq *http.Request, model string) IDETestConnectionResponse {
	start := time.Now()
	httpResp, err := ideHTTPClient.Do(httpReq)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return buildTestFailureResponse(http.StatusServiceUnavailable, latency, err.Error(), "NETWORK_ERROR")
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 400 {
		msg := fmt.Sprintf("endpoint returned HTTP %d", httpResp.StatusCode)

		return buildTestFailureResponse(httpResp.StatusCode, latency, msg, "HTTP_ERROR")
	}

	return buildTestSuccessResponse(httpResp.StatusCode, latency, model)
}

func applyAuthHeaders(req *http.Request, provider, apiKey string) {
	if apiKey == "" {
		return
	}

	if strings.EqualFold(provider, "anthropic") {
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")

		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
}

func buildTestFailureResponse(statusCode int, latency int64, msg, code string) IDETestConnectionResponse {
	data := &IDETestConnectionData{
		IsSuccess:      false,
		LatencyMs:      latency,
		HTTPStatusCode: statusCode,
		ErrorMessage:   msg,
		ErrorCode:      code,
	}

	return IDETestConnectionResponse{
		IsSuccess:      false,
		HTTPStatusCode: statusCode,
		LatencyMs:      latency,
		ErrorMessage:   msg,
		Data:           data,
	}
}

func buildTestSuccessResponse(statusCode int, latency int64, model string) IDETestConnectionResponse {
	echo := "GitMap AI connection verified successfully."
	data := &IDETestConnectionData{
		IsSuccess:      true,
		LatencyMs:      latency,
		Model:          model,
		HTTPStatusCode: statusCode,
		ModelEcho:      echo,
	}

	return IDETestConnectionResponse{
		IsSuccess:      true,
		HTTPStatusCode: statusCode,
		LatencyMs:      latency,
		ModelEcho:      echo,
		Data:           data,
	}
}

func writeTestConnectionError(w http.ResponseWriter, code int, msg string) {
	resp := buildTestFailureResponse(code, 0, msg, "INVALID_REQUEST")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}

func handleIDESetupApply(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)

		return
	}

	var req IDESetupApplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeApplyError(w, http.StatusBadRequest, "invalid request body: "+err.Error())

		return
	}

	normalizeApplyRequest(&req)
	resp := executeApplySetup(req)

	_ = json.NewEncoder(w).Encode(resp)
}

func normalizeApplyRequest(req *IDESetupApplyRequest) {
	if req.Target == "" {
		req.Target = resolveTargetName(req.TargetEditor, req.TargetExtension)
	}
	if req.Provider == "" {
		req.Provider = req.PresetID
	}
	if req.Model == "" {
		req.Model = req.ModelName
	}
	if req.BaseURL == "" {
		req.BaseURL = resolvePresetBaseURL(req.Provider)
	}
}

func resolveTargetName(editor, ext string) string {
	if ext != "" {
		return strings.ToLower(ext)
	}
	if editor != "" {
		return strings.ToLower(editor)
	}

	return "all"
}

func executeApplySetup(req IDESetupApplyRequest) IDESetupApplyResponse {
	if req.IsDryRun {
		return buildDryRunApplyResponse(req)
	}

	return applyLiveSetup(req)
}

func buildDryRunApplyResponse(req IDESetupApplyRequest) IDESetupApplyResponse {
	targets := resolveApplyTargets(req.Target)
	summary := fmt.Sprintf("Dry-run validation successful. Would configure %d target(s): %s with provider '%s' and model '%s'.",
		len(targets), strings.Join(targets, ", "), req.Provider, req.Model)

	data := &IDESetupApplyData{
		IsSuccess:         true,
		ConfiguredTargets: targets,
		BackupPaths:       []string{},
		IsDryRun:          true,
		DiffSummary:       summary,
	}

	return IDESetupApplyResponse{
		IsSuccess:         true,
		IsDryRun:          true,
		DiffSummary:       summary,
		ConfiguredTargets: targets,
		Data:              data,
	}
}

func resolveApplyTargets(target string) []string {
	clean := strings.ToLower(strings.TrimSpace(target))
	switch clean {
	case "cline":
		return []string{"cline"}
	case "continue":
		return []string{"continue"}
	case "zed", "zed-native":
		return []string{"zed"}
	case "vscode":
		return []string{"cline", "continue"}
	default:
		return []string{"cline", "continue", "zed"}
	}
}

func resolveDefaultConfigPath(target string) string {
	home, _ := os.UserHomeDir()
	appData := os.Getenv("APPDATA")

	switch target {
	case "cline":
		return resolveClineConfigPath(home, appData)
	case "continue":
		return filepath.Join(home, ".continue", "config.json")
	case "zed":
		return resolveZedConfigPath(home, appData)
	default:
		return ""
	}
}

func applyLiveSetup(req IDESetupApplyRequest) IDESetupApplyResponse {
	targets := resolveApplyTargets(req.Target)
	configured, backups := processTargets(targets, req)

	return buildApplySuccessResponse(configured, backups, req)
}

func processTargets(targets []string, req IDESetupApplyRequest) ([]string, []string) {
	var configured []string
	var backups []string

	for _, t := range targets {
		destPath := getTargetConfigPath(t)
		bak := performOptionalBackup(destPath, req.IsBackupEnabled)
		if bak != "" {
			backups = append(backups, bak)
		}

		_, err := configureTargetFile(t, req)
		if err == nil {
			configured = append(configured, t)
		}
	}

	return configured, backups
}

func performOptionalBackup(destPath string, isBackupEnabled bool) string {
	if !isBackupEnabled {
		return ""
	}

	bak, err := backupConfigFile(destPath)
	if err != nil {
		return ""
	}

	return bak
}

func backupConfigFile(path string) (string, error) {
	if !isFileExisting(path) {
		return "", nil
	}

	backupPath := fmt.Sprintf("%s.bak.%d", path, time.Now().UnixNano())
	data, err := os.ReadFile(path)
	if err != nil {
		return "", apperror.WrapSimple(err, "read_backup_source")
	}

	err = os.WriteFile(backupPath, data, 0644)
	if err != nil {
		return "", apperror.WrapSimple(err, "write_backup_dest")
	}

	return backupPath, nil
}

func loadOrCreateConfig(path string) map[string]any {
	var cfg map[string]any
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &cfg)
	}
	if cfg == nil {
		cfg = make(map[string]any)
	}

	return cfg
}

func configureTargetFile(target string, req IDESetupApplyRequest) (string, error) {
	path := getTargetConfigPath(target)
	if path == "" {
		return "", apperror.NewSimple("configure_target", "unknown target: "+target)
	}

	cfg := loadOrCreateConfig(path)
	mergeTargetConfig(target, cfg, req)

	return writeConfigFile(path, cfg)
}

func writeConfigFile(path string, cfg map[string]any) (string, error) {
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", apperror.WrapSimple(err, "marshal_config")
	}

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	writeErr := os.WriteFile(path, out, 0644)
	if writeErr != nil {
		return "", apperror.WrapSimple(writeErr, "write_config")
	}

	return path, nil
}

func mergeTargetConfig(target string, cfg map[string]any, req IDESetupApplyRequest) {
	switch target {
	case "cline":
		mergeClineConfig(cfg, req)
	case "continue":
		mergeContinueConfig(cfg, req)
	case "zed":
		mergeZedConfig(cfg, req)
	}
}

func mergeClineConfig(cfg map[string]any, req IDESetupApplyRequest) {
	cfg["apiProvider"] = "openai-compatible"
	cfg["openAiBaseUrl"] = req.BaseURL
	cfg["openAiApiKey"] = req.APIKey
	cfg["openAiModelId"] = req.Model
}

func mergeContinueConfig(cfg map[string]any, req IDESetupApplyRequest) {
	modelEntry := map[string]any{
		"title":    req.Provider + " (" + req.Model + ")",
		"provider": "openai",
		"model":    req.Model,
		"apiKey":   req.APIKey,
		"apiBase":  req.BaseURL,
	}
	cfg["models"] = []any{modelEntry}
	cfg["tabAutocompleteModel"] = modelEntry
}

func mergeZedConfig(cfg map[string]any, req IDESetupApplyRequest) {
	cfg["assistant"] = map[string]any{
		"default_model": map[string]any{
			"provider": "openai",
			"model":    req.Model,
		},
		"version": "2",
	}
	cfg["language_models"] = map[string]any{
		"openai": map[string]any{
			"api_url": req.BaseURL,
		},
	}
}

func buildApplySuccessResponse(configured, backups []string, req IDESetupApplyRequest) IDESetupApplyResponse {
	summary := fmt.Sprintf("Successfully configured %d IDE target(s): %s with provider '%s'.",
		len(configured), strings.Join(configured, ", "), req.Provider)

	data := &IDESetupApplyData{
		IsSuccess:         true,
		ConfiguredTargets: configured,
		BackupPaths:       backups,
		IsDryRun:          false,
		DiffSummary:       summary,
	}

	var firstTargetFile string
	if len(configured) > 0 {
		firstTargetFile = getTargetConfigPath(configured[0])
	}

	return IDESetupApplyResponse{
		IsSuccess:         true,
		TargetFile:        firstTargetFile,
		DiffSummary:       summary,
		ConfiguredTargets: configured,
		BackupPaths:       backups,
		Data:              data,
	}
}

func writeApplyError(w http.ResponseWriter, code int, msg string) {
	resp := IDESetupApplyResponse{
		IsSuccess:    false,
		ErrorMessage: msg,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}
