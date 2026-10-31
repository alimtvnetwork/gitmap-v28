package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAIAutoDefaultNameAndEmail(t *testing.T) {
	name := GetAIAutoDefaultName()
	if len(name) == 0 {
		t.Errorf("expected non-empty default AI auto bot name")
	}

	email := GetAIAutoDefaultEmail()
	if email != DefaultAIAutoBotEmail {
		t.Errorf("expected default email %s, got %s", DefaultAIAutoBotEmail, email)
	}
}

func TestExtractAttrFromJSON(t *testing.T) {
	raw := []byte(`{
		"attributes": {
			"ai_auto.default_agent": "muse",
			"ai_auto.default_name": "custom bot"
		}
	}`)
	agent := extractAttrFromJSON(raw, SettingAIAutoDefaultAgent)
	if agent != "muse" {
		t.Errorf("expected muse, got %s", agent)
	}
}

func TestExtractAttrFromJSON_PayloadEnvelope(t *testing.T) {
	raw := []byte(`{
		"payload": {
			"attributes": {
				"ai_auto.default_agent": "gemini"
			}
		}
	}`)
	agent := extractAttrFromJSON(raw, SettingAIAutoDefaultAgent)
	if agent != "gemini" {
		t.Errorf("expected gemini, got %s", agent)
	}
}

func TestUpdateUISettingsFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "ui_settings.json")
	if err := updateUISettingsFile(tmp, "test_key", "test_val"); err != nil {
		t.Fatalf("unexpected error updating ui settings file: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	val := extractAttrFromJSON(data, "test_key")
	if val != "test_val" {
		t.Errorf("expected test_val, got %s", val)
	}
}

func TestSetAIAutoSetting(t *testing.T) {
	err := SetAIAutoSetting(SettingAIAutoDefaultAgent, "antigravity")
	if err != nil {
		t.Fatalf("unexpected error setting agent: %v", err)
	}
	agent := GetAIAutoDefaultAgent()
	if agent != "antigravity" {
		t.Errorf("expected antigravity, got %s", agent)
	}
}
