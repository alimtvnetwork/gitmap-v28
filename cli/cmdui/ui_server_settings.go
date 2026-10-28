package cmdui

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
)

var customSettingsPath string

func getSettingsFilePath() string {
	if customSettingsPath != "" {
		return customSettingsPath
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".gitmap", "ui_settings.json")
}

func loadSettings() SettingsData {
	def := SettingsData{
		Theme:           "dark",
		DefaultRemote:   "origin",
		ClusterPort:     49152,
		AutoDeployKey:   true,
		GraphicsMode:    "high",
		AutoOpenBrowser: true,
		BrowserMode:     "default",
		CommitInLayout:  "split",
		PullDirection:   "pull-left",
		PRReplayMode:    "merges",
		Attributes:      make(map[string]string),
	}
	path := getSettingsFilePath()
	if path == "" {
		return def
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return def
	}

	payload, _, extractErr := jsonx.ExtractPayload(data)
	if extractErr != nil {
		payload = data
	}

	var loaded SettingsData
	err = json.Unmarshal(payload, &loaded)
	if err != nil {
		return def
	}

	if loaded.Attributes == nil {
		loaded.Attributes = make(map[string]string)
	}

	return loaded
}

func saveSettingsData(s SettingsData) error {
	path := getSettingsFilePath()
	if path == "" {
		return nil
	}

	if s.Attributes == nil {
		s.Attributes = make(map[string]string)
	}

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	envelope := jsonx.NewEnvelope(
		jsonx.TypeUISettings,
		path,
		"gitmap ui settings",
		"1.0",
		s,
	)
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ImportSettingsFromFile loads settings from a file (supporting JSON envelope) and persists them to ui_settings.json.
func ImportSettingsFromFile(filePath string) error {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	payload, _, err := jsonx.ExtractPayload(raw)
	if err == nil && len(payload) > 0 {
		raw = payload
	}

	var loaded SettingsData
	err = json.Unmarshal(raw, &loaded)
	if err != nil {
		return err
	}

	return saveSettingsData(loaded)
}

func mergeSettings(current, req SettingsData) SettingsData {
	if req.Theme != "" {
		current.Theme = req.Theme
	}

	if req.DefaultRemote != "" {
		current.DefaultRemote = req.DefaultRemote
	}

	if req.ClusterPort > 0 {
		current.ClusterPort = req.ClusterPort
	}

	current.AutoDeployKey = req.AutoDeployKey

	if req.GraphicsMode != "" {
		current.GraphicsMode = req.GraphicsMode
	}

	current.AutoOpenBrowser = req.AutoOpenBrowser

	if req.BrowserMode != "" {
		current.BrowserMode = req.BrowserMode
		_ = os.Setenv("GITMAP_BROWSER_MODE", req.BrowserMode)
	}

	if req.CommitInLayout != "" {
		current.CommitInLayout = req.CommitInLayout
	}

	if req.PullDirection != "" {
		current.PullDirection = req.PullDirection
	}

	if req.PRReplayMode != "" {
		current.PRReplayMode = req.PRReplayMode
	}

	if current.Attributes == nil {
		current.Attributes = make(map[string]string)
	}

	if req.Attributes != nil {
		for k, v := range req.Attributes {
			current.Attributes[k] = v
		}
	}

	return current
}

func handleAPISettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)

	if r.Method == http.MethodGet {
		_ = json.NewEncoder(w).Encode(loadSettings())
		return
	}

	var req SettingsData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	merged := mergeSettings(loadSettings(), req)
	saveErr := saveSettingsData(merged)
	if saveErr != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": saveErr.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
}
