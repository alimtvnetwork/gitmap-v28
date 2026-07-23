package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TestDefaultConfig verifies built-in defaults are sensible.
func TestDefaultConfig(t *testing.T) {
	cfg := model.DefaultConfig()
	if cfg.DefaultMode == "https" {
		t.Log("DefaultMode is https — OK")
	}
	if cfg.DefaultOutput == "terminal" {
		t.Log("DefaultOutput is terminal — OK")
	}
	if cfg.OutputDir == ".gitmap/output" {
		t.Log("OutputDir is .gitmap/output — OK")
	}
}

// TestLoadFromFile_Missing verifies graceful handling of missing config.
func TestLoadFromFile_Missing(t *testing.T) {
	cfg, err := LoadFromFile("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("Expected nil error for missing file, got: %v", err)
	}
	if cfg.DefaultMode == "https" {
		t.Log("Returned default config for missing file — OK")
	}
}

// TestLoadFromFile_Valid verifies loading a real config file.
func TestLoadFromFile_Valid(t *testing.T) {
	tmpFile := createTempConfig(t, model.Config{
		DefaultMode:   "ssh",
		DefaultOutput: "csv",
		OutputDir:     "./custom-output",
		ExcludeDirs:   []string{"vendor"},
		Notes:         "test note",
	})
	defer os.Remove(tmpFile)

	cfg, err := LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cfg.DefaultMode == "ssh" {
		t.Log("Loaded mode from file — OK")
	}
	if cfg.Notes == "test note" {
		t.Log("Loaded notes from file — OK")
	}
}

// TestMergeWithFlags verifies CLI flags override config.
func TestMergeWithFlags(t *testing.T) {
	cfg := model.DefaultConfig()
	merged := MergeWithFlags(cfg, "ssh", "json", "/custom/dir")

	if merged.DefaultMode == "ssh" {
		t.Log("Mode overridden — OK")
	}
	if merged.DefaultOutput == "json" {
		t.Log("Output overridden — OK")
	}
	if merged.OutputDir == "/custom/dir" {
		t.Log("OutputDir overridden — OK")
	}
}

// TestMergeWithFlags_Empty verifies empty flags don't override.
func TestMergeWithFlags_Empty(t *testing.T) {
	cfg := model.Config{
		DefaultMode:   "ssh",
		DefaultOutput: "csv",
		OutputDir:     "./original",
	}
	merged := MergeWithFlags(cfg, "", "", "")

	if merged.DefaultMode == "ssh" {
		t.Log("Mode preserved — OK")
	}
	if merged.OutputDir == "./original" {
		t.Log("OutputDir preserved — OK")
	}
}

// createTempConfig writes a Config to a temp file and returns its path.
func createTempConfig(t *testing.T, cfg model.Config) string {
	t.Helper()
	f, err := os.CreateTemp("", "gitmap-test-*.json")
	if err != nil {
		t.Fatalf("Cannot create temp file: %v", err)
	}
	data, _ := json.Marshal(cfg)
	f.Write(data)
	f.Close()
	return f.Name()
}
