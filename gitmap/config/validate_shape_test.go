package config

// Tests for the second-tier validation layer (validate_shape.go):
//
//   - ValidateRawShape       -- per-required-key JSON-type checks.
//   - ValidateConfigStruct   -- post-unmarshal shape rules
//                               (non-empty outputDir, etc.).
//
// Plus end-to-end LoadFromFile coverage so a future refactor that
// drops a call site is caught by the test suite, not by users.

// The writeTempConfigBytes helper is shared with validate_test.go.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// TestValidateRawShape_RejectsWrongTypePerKey covers the canonical
// failure mode the shape check exists for: the user wrote a number
// where a string was expected. Asserting the error names BOTH the
// offending key and the expected type, so the message is actionable
// without re-reading the docs.
func TestValidateRawShape_RejectsWrongTypePerKey(t *testing.T) {
	body := []byte(`{"defaultMode":42,"defaultOutput":"terminal","outputDir":"./o"}`)
	err := ValidateRawShape(body)
	if err == nil {
		t.Fatal("expected shape error for numeric defaultMode, got nil")
	}
	for _, want := range []string{"defaultMode", "string", "number"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
}

// TestValidateRawShape_AggregatesAllMismatches asserts every
// per-key violation lands in a single error so the user can fix
// them in one edit cycle (matches ValidateRawConfig's aggregation
// contract).
func TestValidateRawShape_AggregatesAllMismatches(t *testing.T) {
	body := []byte(`{"defaultMode":1,"defaultOutput":true,"outputDir":[]}`)
	err := ValidateRawShape(body)
	if err == nil {
		t.Fatal("expected aggregated shape error, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"defaultMode", "defaultOutput", "outputDir"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
}

// TestValidateRawShape_AcceptsCorrectTypes is the happy path — the
// values are wrong by enum (handled by ValidateConfig downstream)
// but their JSON KIND is right, so the shape layer must pass.
func TestValidateRawShape_AcceptsCorrectTypes(t *testing.T) {
	body := []byte(`{"defaultMode":"x","defaultOutput":"y","outputDir":"z"}`)
	if err := ValidateRawShape(body); err != nil {
		t.Errorf("expected nil for correctly-typed values, got: %v", err)
	}
}

// TestValidateRawShape_IgnoresMalformedJSON delegates the malformed-
// JSON error to ValidateRawConfig (already covered) and returns nil
// here so we don't double-report. Pinned by this test because a
// "return the unmarshal error here too" refactor would superficially
// work but produce a confusing two-line error in LoadFromFile.
func TestValidateRawShape_IgnoresMalformedJSON(t *testing.T) {
	if err := ValidateRawShape([]byte(`{not json`)); err != nil {
		t.Errorf("shape check must defer to raw check on bad JSON, got: %v", err)
	}
}

// TestValidateConfigStruct_RejectsEmptyOutputDir guards against a
// silent-write-to-cwd footgun: an explicit empty outputDir would
// otherwise be accepted by the typed unmarshal.
func TestValidateConfigStruct_RejectsEmptyOutputDir(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.OutputDir = ""
	err := ValidateConfigStruct(cfg)
	if err == nil {
		t.Fatal("expected error for empty outputDir, got nil")
	}
	if !strings.Contains(err.Error(), "outputDir") {
		t.Errorf("error must name outputDir: %v", err)
	}
}

// TestValidateConfigStruct_RejectsWhitespaceOutputDir covers the
// "user typed a space by accident" variant. Treated the same as
// empty so no one ends up with a directory literally named " ".
func TestValidateConfigStruct_RejectsWhitespaceOutputDir(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.OutputDir = "   "
	if err := ValidateConfigStruct(cfg); err == nil {
		t.Fatal("expected error for whitespace-only outputDir, got nil")
	}
}

// TestValidateConfigStruct_RejectsNegativeDashboardRefresh pins the
// "must be >= 0" rule and asserts the error reports the actual bad
// value so the user doesn't have to grep their config to find it.
func TestValidateConfigStruct_RejectsNegativeDashboardRefresh(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.DashboardRefresh = -5
	err := ValidateConfigStruct(cfg)
	if err == nil {
		t.Fatal("expected error for negative dashboardRefresh, got nil")
	}
	if !strings.Contains(err.Error(), "-5") {
		t.Errorf("error must report bad value: %v", err)
	}
}

// TestValidateConfigStruct_AllowsZeroDashboardRefresh pins the
// "zero means default" contract — the validator must not encroach
// on it. Regression guard for a future "make zero invalid" change
// that would break every user with no explicit dashboardRefresh.
func TestValidateConfigStruct_AllowsZeroDashboardRefresh(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.DashboardRefresh = 0
	if err := ValidateConfigStruct(cfg); err != nil {
		t.Errorf("zero dashboardRefresh must be accepted: %v", err)
	}
}

// TestValidateConfigStruct_RejectsEmptyExcludeEntry guards against
// the "every directory matches an empty pattern" footgun. Asserts
// the error names the offending index so the user knows where to
// edit.
func TestValidateConfigStruct_RejectsEmptyExcludeEntry(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.ExcludeDirs = []string{"vendor", "", "node_modules"}
	err := ValidateConfigStruct(cfg)
	if err == nil {
		t.Fatal("expected error for empty excludeDirs entry, got nil")
	}
	if !strings.Contains(err.Error(), "excludeDirs[1]") {
		t.Errorf("error must name the bad index: %v", err)
	}
}

// TestValidateConfigStruct_RejectsIncompleteReleaseTarget covers
// both halves (missing goos, missing goarch) in one test because
// the validation logic is symmetric and aggregating both into a
// single error is part of the contract.
func TestValidateConfigStruct_RejectsIncompleteReleaseTarget(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Release.Targets = []model.ReleaseTarget{
		{GOOS: "linux", GOARCH: ""},
		{GOOS: "", GOARCH: "amd64"},
	}
	err := ValidateConfigStruct(cfg)
	if err == nil {
		t.Fatal("expected error for incomplete release targets, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"release.targets[0].goarch", "release.targets[1].goos"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
}

// TestValidateConfigStruct_AcceptsDefaultConfig is the bedrock
// guarantee: the built-in defaults must pass every shape rule.
// If they don't, every fresh install fails at startup -- catching
// that here is much cheaper than catching it in production.
func TestValidateConfigStruct_AcceptsDefaultConfig(t *testing.T) {
	if err := ValidateConfigStruct(model.DefaultConfig()); err != nil {
		t.Errorf("DefaultConfig must satisfy ValidateConfigStruct: %v", err)
	}
}

// TestLoadFromFile_FailsFastOnTypeMismatch is the end-to-end wiring
// check for the shape layer: a wrong-typed value on disk must
// surface as a shape error from LoadFromFile, NOT as the raw
// json.Unmarshal error (which would mean the validator wasn't called).
func TestLoadFromFile_FailsFastOnTypeMismatch(t *testing.T) {
	path := writeTempConfigBytes(t,
		`{"defaultMode":42,"defaultOutput":"terminal","outputDir":"./o"}`,
	)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for numeric defaultMode, got nil")
	}
	if !strings.Contains(err.Error(), "type mismatch") {
		t.Errorf("expected shape-layer error, got: %v", err)
	}
}

// TestLoadFromFile_FailsFastOnNegativeRefresh is the end-to-end
// wiring check for the struct-level layer.
func TestLoadFromFile_FailsFastOnNegativeRefresh(t *testing.T) {
	path := writeTempConfigBytes(t,
		`{"defaultMode":"https","defaultOutput":"terminal","outputDir":"./o","dashboardRefresh":-1}`,
	)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for negative dashboardRefresh, got nil")
	}
	if !strings.Contains(err.Error(), "dashboardRefresh") {
		t.Errorf("expected struct-layer error naming dashboardRefresh, got: %v", err)
	}
}
