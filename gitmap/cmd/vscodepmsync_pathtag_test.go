// Package cmd — vscodepmsync_pathtag_test.go: end-to-end tests for
// the v4.37.0 --projects-json + --tag flags on
// `gitmap vscode-pm-sync`.
//
// --projects-json bypasses VS Code discovery entirely, so these
// tests deliberately do NOT swap HOME / XDG_CONFIG_HOME — the
// override path must work even when the resolver would otherwise
// return ErrUserDataMissing. --tag replaces the per-pair detected
// tag set verbatim; the brand "gitmap" is NOT auto-prepended.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// writeOverrideFixture creates a repo dir + a projects.json at an
// arbitrary path (NOT inside any XDG / APPDATA layout) seeded with a
// single entry carrying seedTags. Returns (overridePath, repoPath).
func writeOverrideFixture(t *testing.T, seedTags []string) (string, string) {
	t.Helper()
	tmp := t.TempDir()
	repoDir := filepath.Join(tmp, "demo-repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	overridePath := filepath.Join(tmp, "custom-projects.json")
	seed := []vscodepm.Entry{{
		Name: "demo-repo", RootPath: repoDir,
		Paths: []string{}, Tags: seedTags, Enabled: true,
	}}
	data, err := json.MarshalIndent(seed, "", "\t")
	if err != nil {
		t.Fatalf("marshal seed: %v", err)
	}
	if err := os.WriteFile(overridePath, data, 0o644); err != nil {
		t.Fatalf("write seed: %v", err)
	}

	return overridePath, repoDir
}

// readOverrideFixture unmarshals projects.json at the explicit path
// — loadVSCodePMSyncProjectsJSON in the existing helper goes through
// the resolver, which is exactly what --projects-json bypasses.
func readOverrideFixture(t *testing.T, path string) []vscodepm.Entry {
	t.Helper()
	got, err := vscodepm.ListEntriesAt(path)
	if err != nil {
		t.Fatalf("list at %s: %v", path, err)
	}
	return got
}

// TestVSCodePMSyncProjectsJSONOverrideWritesToOverridePath asserts
// the runner reads AND writes the override path verbatim, never
// touching the XDG-resolved location. We don't swap HOME so any
// accidental fallback would surface as ErrUserDataMissing.
func TestVSCodePMSyncProjectsJSONOverrideWritesToOverridePath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	overridePath, _ := writeOverrideFixture(t, []string{"user-only"})

	runVSCodePMSync([]string{"--projects-json", overridePath})

	got := readOverrideFixture(t, overridePath)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	// UNION default — user tag preserved, brand added by detector.
	if !containsTag(got[0].Tags, "user-only") {
		t.Errorf("override path did not preserve user tag: %v", got[0].Tags)
	}
	if !containsTag(got[0].Tags, "gitmap") {
		t.Errorf("override path missing brand tag: %v", got[0].Tags)
	}
}

// TestVSCodePMSyncTagOverrideReplacesDetectedSet asserts --tag
// substitutes the per-pair detected tags with the user-supplied list
// verbatim. Combined with --mode replace this also drops every
// pre-existing on-disk tag, including the brand — proving the
// "user owns the full set" contract from constants_cli.go:191-193.
func TestVSCodePMSyncTagOverrideReplacesDetectedSet(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	overridePath, _ := writeOverrideFixture(t, []string{"old-tag", "gitmap"})

	runVSCodePMSync([]string{
		"--projects-json", overridePath,
		"--mode", "replace",
		"--tag", "ruby,python",
		"--tag", "scratch",
	})

	got := readOverrideFixture(t, overridePath)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	for _, want := range []string{"ruby", "python", "scratch"} {
		if !containsTag(got[0].Tags, want) {
			t.Errorf("--tag override missing %q in: %v", want, got[0].Tags)
		}
	}
	for _, gone := range []string{"old-tag", "gitmap"} {
		if containsTag(got[0].Tags, gone) {
			t.Errorf("replace+--tag did not drop %q: %v", gone, got[0].Tags)
		}
	}
}

// TestVSCodePMSyncTagOverrideUnionKeepsExisting asserts the default
// --mode union still merges --tag values with what is already on
// disk (user/CI may want to bulk-add a tag without replacing).
// Brand tag survives only because it was already present on disk —
// the override list itself does NOT auto-prepend it.
func TestVSCodePMSyncTagOverrideUnionKeepsExisting(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	overridePath, _ := writeOverrideFixture(t, []string{"keep-me", "gitmap"})

	runVSCodePMSync([]string{
		"--projects-json", overridePath,
		"--tag", "added",
	})

	got := readOverrideFixture(t, overridePath)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	for _, want := range []string{"keep-me", "gitmap", "added"} {
		if !containsTag(got[0].Tags, want) {
			t.Errorf("union with --tag missing %q in: %v", want, got[0].Tags)
		}
	}
}
