package config

// Tests for the schema-validation layer added on top of LoadFromFile.
// Two surfaces under test:
//
//   - ValidateRawConfig: missing-required-key detection on raw bytes.
//   - ValidateConfig:    invalid enum detection on the typed struct.
//
// Plus end-to-end coverage through LoadFromFile to confirm the
// validators are actually wired into the config-loading hot path
// (regression guard: a future refactor that drops the call sites
// would silently re-introduce the "broken config limps along" bug
// these checks were added to prevent).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TestValidateRawConfig_MissingKeysAggregated verifies that EVERY
// missing required key shows up in a single error — not just the
// first one — so users can fix the file in one edit cycle.
func TestValidateRawConfig_MissingKeysAggregated(t *testing.T) {
	// Empty object → all three required keys are missing.
	err := ValidateRawConfig([]byte(`{}`))
	if err == nil {
		t.Fatal("expected error for empty config object, got nil")
	}
	for _, want := range []string{"defaultMode", "defaultOutput", "outputDir"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error missing key %q: %v", want, err)
		}
	}
}

// TestValidateRawConfig_AllPresentNoError is the happy path — every
// required key present, no error. Values themselves are not
// inspected at this layer (that's ValidateConfig's job).
func TestValidateRawConfig_AllPresentNoError(t *testing.T) {
	raw := []byte(`{"defaultMode":"https","defaultOutput":"terminal","outputDir":"./out"}`)
	if err := ValidateRawConfig(raw); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestValidateRawConfig_InvalidJSON surfaces a wrapped JSON parse
// error rather than panicking or returning a misleading
// "missing keys" message.
func TestValidateRawConfig_InvalidJSON(t *testing.T) {
	err := ValidateRawConfig([]byte(`{not valid json`))
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("expected wrapped 'invalid JSON' message, got: %v", err)
	}
}

// TestValidateConfig_RejectsEmptyEnumValue verifies that an explicit
// `"defaultMode": ""` fails — silently falling back to a default
// would be the exact bug this validator exists to prevent.
func TestValidateConfig_RejectsEmptyEnumValue(t *testing.T) {
	cfg := model.Config{
		DefaultMode:   "",
		DefaultOutput: "terminal",
		OutputDir:     "./out",
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for empty defaultMode, got nil")
	}
	if !strings.Contains(err.Error(), "defaultMode") {
		t.Errorf("error must name the offending key: %v", err)
	}
}

// TestValidateConfig_RejectsTypoEnumValue verifies typos in enum
// values fail with a message that lists the accepted set, so the
// fix is obvious from reading the error alone.
func TestValidateConfig_RejectsTypoEnumValue(t *testing.T) {
	cfg := model.Config{
		DefaultMode:   "htps",    // typo
		DefaultOutput: "trminal", // typo
		OutputDir:     "./out",
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for typo'd enums, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"defaultMode", "defaultOutput", "https", "terminal"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
}

// TestValidateConfig_AcceptsAllValidEnums sweeps the cartesian
// product of valid mode × output values to guard against a future
// refactor accidentally narrowing the accepted set.
func TestValidateConfig_AcceptsAllValidEnums(t *testing.T) {
	for _, mode := range validModeValues {
		for _, output := range validOutputValues {
			cfg := model.Config{
				DefaultMode: mode, DefaultOutput: output, OutputDir: "./out",
			}
			if err := ValidateConfig(cfg); err != nil {
				t.Errorf("rejected valid combo mode=%s output=%s: %v", mode, output, err)
			}
		}
	}
}

// TestLoadFromFile_FailsFastOnMissingKey is the end-to-end wiring
// check: a partial config file on disk must propagate the
// validation error out of LoadFromFile, not be silently patched
// with defaults. Regression guard for the original "limps along"
// bug this whole feature was added to prevent.
func TestLoadFromFile_FailsFastOnMissingKey(t *testing.T) {
	path := writeTempConfigBytes(t, `{"defaultMode":"https"}`)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for partial config, got nil")
	}
	if !strings.Contains(err.Error(), "defaultOutput") {
		t.Errorf("expected error to name missing 'defaultOutput', got: %v", err)
	}
}

// TestLoadFromFile_FailsFastOnInvalidEnum is the second wiring
// check: present-but-invalid values must also bubble out.
func TestLoadFromFile_FailsFastOnInvalidEnum(t *testing.T) {
	path := writeTempConfigBytes(t,
		`{"defaultMode":"sftp","defaultOutput":"terminal","outputDir":"./o"}`,
	)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for invalid enum, got nil")
	}
	if !strings.Contains(err.Error(), "sftp") {
		t.Errorf("expected error to name bad value 'sftp', got: %v", err)
	}
}

// writeTempConfigBytes drops `body` into a temp .json file and
// returns its path. Used instead of createTempConfig (the existing
// test helper that goes through json.Marshal) because we need to
// emit DELIBERATELY-broken JSON shapes that json.Marshal can't.
func writeTempConfigBytes(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	return path
}
