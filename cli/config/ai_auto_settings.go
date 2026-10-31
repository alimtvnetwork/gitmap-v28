package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const (
	SettingAIAutoDefaultName  = "ai_auto.default_name"
	SettingAIAutoDefaultEmail = "ai_auto.default_email"
	SettingAIAutoDefaultAgent = "ai_auto.default_agent"

	DefaultAIAutoBotName  = "gitmap automation bot - runs autonomously"
	DefaultAIAutoBotEmail = "info@rgitmap.com"
)

type AIAutoSettings struct {
	DefaultName  string `json:"default_name"`
	DefaultEmail string `json:"default_email"`
	DefaultAgent string `json:"default_agent"`
}

func GetAIAutoDefaultName() string {
	if val := readUISettingsAttr(SettingAIAutoDefaultName); len(val) > 0 {
		return val
	}
	if val := readGlobalVar(SettingAIAutoDefaultName); len(val) > 0 {
		return val
	}
	if val := readDBSetting(SettingAIAutoDefaultName); len(val) > 0 {
		return val
	}

	return DefaultAIAutoBotName
}

func GetAIAutoDefaultEmail() string {
	if val := readUISettingsAttr(SettingAIAutoDefaultEmail); len(val) > 0 {
		return val
	}
	if val := readGlobalVar(SettingAIAutoDefaultEmail); len(val) > 0 {
		return val
	}
	if val := readDBSetting(SettingAIAutoDefaultEmail); len(val) > 0 {
		return val
	}

	return DefaultAIAutoBotEmail
}

func GetAIAutoDefaultAgent() string {
	if val := readUISettingsAttr(SettingAIAutoDefaultAgent); len(val) > 0 {
		return val
	}
	if val := readGlobalVar(SettingAIAutoDefaultAgent); len(val) > 0 {
		return val
	}

	return readDBSetting(SettingAIAutoDefaultAgent)
}

func readGlobalVar(key string) string {
	val, hasVal := GetVariable("global", key)
	if !hasVal {
		return ""
	}

	return strings.TrimSpace(val)
}

func readDBSetting(key string) string {
	dbConn, errOpen := store.OpenDefault()
	if errOpen != nil {
		return ""
	}
	defer dbConn.Close()

	return strings.TrimSpace(dbConn.GetSetting(key))
}

func readUISettingsAttr(key string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	path := filepath.Join(home, ".gitmap", "ui_settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return extractAttrFromJSON(data, key)
}

func extractAttrFromJSON(data []byte, key string) string {
	var payload struct {
		Payload struct {
			Attributes map[string]string `json:"attributes"`
		} `json:"payload"`
		Attributes map[string]string `json:"attributes"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return ""
	}
	if len(payload.Payload.Attributes[key]) > 0 {
		return strings.TrimSpace(payload.Payload.Attributes[key])
	}

	return strings.TrimSpace(payload.Attributes[key])
}

func SetAIAutoSetting(key, val string) error {
	trimmed := strings.TrimSpace(val)
	_ = SetVariable("global", key, trimmed)
	saveDBSetting(key, trimmed)

	return writeUISettingsAttr(key, trimmed)
}

func saveDBSetting(key, val string) {
	dbConn, errOpen := store.OpenDefault()
	if errOpen == nil {
		defer dbConn.Close()
		_ = dbConn.SetSetting(key, val)
	}
}

func writeUISettingsAttr(key, val string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".gitmap", "ui_settings.json")

	return updateUISettingsFile(path, key, val)
}

func updateUISettingsFile(path, key, val string) error {
	var root map[string]any
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &root)
	}
	if root == nil {
		root = make(map[string]any)
	}
	applyAttrToRoot(root, key, val)
	out, mErr := json.MarshalIndent(root, "", "  ")
	if mErr != nil {
		return mErr
	}
	_ = os.MkdirAll(filepath.Dir(path), 0755)

	return os.WriteFile(path, out, 0644)
}

func applyAttrToRoot(root map[string]any, key, val string) {
	if p, ok := root["payload"].(map[string]any); ok {
		attrs, _ := p["attributes"].(map[string]any)
		if attrs == nil {
			attrs = make(map[string]any)
			p["attributes"] = attrs
		}
		attrs[key] = val

		return
	}
	attrs, _ := root["attributes"].(map[string]any)
	if attrs == nil {
		attrs = make(map[string]any)
		root["attributes"] = attrs
	}
	attrs[key] = val
}
