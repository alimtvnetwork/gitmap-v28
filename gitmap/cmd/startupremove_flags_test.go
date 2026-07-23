package cmd

// Cross-platform tests for parseStartupRemoveFlags. Validates that
// --backend, --dry-run, and --output parse independently, in any
// order, and that the positional name comes through cleanly. These
// tests are parser-only — they don't exercise startup.RemoveWith
// Options, so they pass on every OS without touching the registry
// or filesystem.

import (
	"testing"
)

// TestParseStartupRemoveFlags_PositionalOnly is the baseline:
// no flags, just a name. Should return name + zero dryRun + empty
// backend (= unspecified / dual-backend fallback on Windows).
func TestParseStartupRemoveFlags_PositionalOnly(t *testing.T) {
	cfg, err := parseStartupRemoveFlags([]string{"foo"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.name != "foo" || cfg.dryRun || cfg.backend != "" {
		t.Errorf("got name=%q dryRun=%v backend=%q, want foo/false/\"\"",
			cfg.name, cfg.dryRun, cfg.backend)
	}
}

// TestParseStartupRemoveFlags_Backend confirms --backend=registry
// is captured into the returned struct. ParseBackend is responsible
// for validating the value — this test only proves the flag wiring
// routes the value through.
func TestParseStartupRemoveFlags_Backend(t *testing.T) {
	cfg, err := parseStartupRemoveFlags(
		[]string{"--backend=registry", "foo"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.backend != "registry" {
		t.Errorf("backend = %q, want registry", cfg.backend)
	}
}

// TestParseStartupRemoveFlags_BackendStartupFolder confirms the
// other valid backend value also routes through. Together with
// _Backend above, this proves the flag accepts both Windows
// backends without per-value special-casing.
func TestParseStartupRemoveFlags_BackendStartupFolder(t *testing.T) {
	cfg, err := parseStartupRemoveFlags(
		[]string{"--backend=startup-folder", "foo"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.backend != "startup-folder" {
		t.Errorf("backend = %q, want startup-folder", cfg.backend)
	}
}

// TestParseStartupRemoveFlags_BackendAndDryRun confirms the two
// flags compose without interfering. Order of flags on the command
// line should not matter — flag.Parse handles arbitrary ordering
// of named flags, but positional args must come last.
func TestParseStartupRemoveFlags_BackendAndDryRun(t *testing.T) {
	cfg, err := parseStartupRemoveFlags(
		[]string{"--dry-run", "--backend=startup-folder", "myapp"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.name != "myapp" || !cfg.dryRun || cfg.backend != "startup-folder" {
		t.Errorf("got name=%q dryRun=%v backend=%q, want myapp/true/startup-folder",
			cfg.name, cfg.dryRun, cfg.backend)
	}
}

// TestParseStartupRemoveFlags_MissingName confirms a missing
// positional name surfaces as an error so the caller exits 2 with
// the usage message rather than silently no-op'ing.
func TestParseStartupRemoveFlags_MissingName(t *testing.T) {
	if _, err := parseStartupRemoveFlags([]string{"--dry-run"}); err == nil {
		t.Fatal("expected error for missing positional name, got nil")
	}
}

// TestParseStartupRemoveFlags_OutputAndIndent confirms the new
// --output and --json-indent flags route into the struct so the
// JSON status emitter can read them downstream.
func TestParseStartupRemoveFlags_OutputAndIndent(t *testing.T) {
	cfg, err := parseStartupRemoveFlags(
		[]string{"--output=json", "--json-indent=0", "foo"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.output != "json" || cfg.jsonIndent != 0 {
		t.Errorf("got output=%q indent=%d, want json/0", cfg.output, cfg.jsonIndent)
	}
}
