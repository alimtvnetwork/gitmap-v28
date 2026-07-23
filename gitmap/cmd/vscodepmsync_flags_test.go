// Package cmd — vscodepmsync_flags_test.go: pure unit tests for
// parseVSCodePMSyncFlags + the tagListValue flag.Value adapter.
//
// These tests do NOT touch the filesystem and do NOT require any
// VS Code state — they exercise the parsing layer in isolation so
// regressions surface immediately on every CI run, regardless of
// GOOS or whether projects.json exists.
package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// TestTagListValueRepeatAndCommaList asserts the documented contract
// from constants_cli.go: `--tag a --tag b,c` produces {a,b,c}.
// Regression guard for the v4.37.0 wiring of --tag, which was
// documented in v4.36.0 but inert until this release.
func TestTagListValueRepeatAndCommaList(t *testing.T) {
	v := newTagListValue()
	for _, raw := range []string{"a", "b,c", " d , e ", "a"} {
		if err := v.Set(raw); err != nil {
			t.Fatalf("Set(%q) error: %v", raw, err)
		}
	}

	want := []string{"a", "b", "c", "d", "e"}
	if len(v.values) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(v.values), len(want), v.values)
	}
	for i := range want {
		if v.values[i] != want[i] {
			t.Errorf("values[%d] = %q, want %q", i, v.values[i], want[i])
		}
	}
	if !v.wasSet {
		t.Error("wasSet = false after Set, want true")
	}
}

// TestTagListValueEmptyEntriesDropped asserts trailing/empty
// comma-separated parts are silently dropped (they are never
// useful and leaking them as "" tags would corrupt projects.json).
func TestTagListValueEmptyEntriesDropped(t *testing.T) {
	v := newTagListValue()
	if err := v.Set(",,a,,b,"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if len(v.values) != 2 || v.values[0] != "a" || v.values[1] != "b" {
		t.Fatalf("values = %v, want [a b]", v.values)
	}
}

// TestParseVSCodePMSyncFlagsDefaults asserts the no-args case
// matches the v4.36.0 baseline: dry-run off, mode=union, no override
// path, no tag override.
func TestParseVSCodePMSyncFlagsDefaults(t *testing.T) {
	opts, err := parseVSCodePMSyncFlags(nil)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if opts.DryRun {
		t.Error("DryRun = true, want false")
	}
	if opts.Mode != vscodepm.MergeModeUnion {
		t.Errorf("Mode = %v, want union", opts.Mode)
	}
	if opts.ProjectsJSON != "" {
		t.Errorf("ProjectsJSON = %q, want \"\"", opts.ProjectsJSON)
	}
	if opts.HasTagOverride {
		t.Error("HasTagOverride = true, want false")
	}
}

// TestParseVSCodePMSyncFlagsAllSet asserts every flag round-trips
// into the opts struct. Uses every documented form:
//
//   - --dry-run            (bool)
//   - --mode replace       (string)
//   - --projects-json /tmp (string)
//   - --tag a --tag b,c    (custom flag.Value, repeat + comma)
func TestParseVSCodePMSyncFlagsAllSet(t *testing.T) {
	args := []string{
		"--dry-run",
		"--mode", "replace",
		"--projects-json", "/tmp/projects.json",
		"--tag", "a",
		"--tag", "b,c",
	}
	opts, err := parseVSCodePMSyncFlags(args)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !opts.DryRun {
		t.Error("DryRun = false, want true")
	}
	if opts.Mode != vscodepm.MergeModeReplace {
		t.Errorf("Mode = %v, want replace", opts.Mode)
	}
	if opts.ProjectsJSON != "/tmp/projects.json" {
		t.Errorf("ProjectsJSON = %q", opts.ProjectsJSON)
	}
	if !opts.HasTagOverride {
		t.Fatal("HasTagOverride = false, want true")
	}

	want := []string{"a", "b", "c"}
	if len(opts.TagOverride) != len(want) {
		t.Fatalf("TagOverride = %v, want %v", opts.TagOverride, want)
	}
	for i := range want {
		if opts.TagOverride[i] != want[i] {
			t.Errorf("TagOverride[%d] = %q, want %q", i, opts.TagOverride[i], want[i])
		}
	}
}

// TestParseVSCodePMSyncFlagsBadModeFailsLoud asserts the
// zero-swallow contract: an unknown --mode literal returns a
// non-nil error rather than silently defaulting to union.
func TestParseVSCodePMSyncFlagsBadModeFailsLoud(t *testing.T) {
	_, err := parseVSCodePMSyncFlags([]string{"--mode", "bogus"})
	if err == nil {
		t.Fatal("err = nil, want non-nil for unknown --mode value")
	}
}
