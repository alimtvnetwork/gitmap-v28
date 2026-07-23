// Package cmd — vscodepmsync_test.go: regression coverage for
// `gitmap vscode-pm-sync` (alias `vpm`).
//
// The runner itself is glue around vscodepm.Sync, which has its own
// extensive test suite. These tests focus on the glue:
//
//  1. Pair construction skips entries whose rootPath no longer exists.
//  2. Pair construction always re-detects tags via DetectTagsCustom,
//     so the brand "gitmap" tag lands on every surviving entry.
//  3. Existing user-added tags survive (Sync's mergePairs UNION).
//
// We write a real projects.json on disk under a temp HOME so the
// path resolver picks it up without mocks. The tests are skipped on
// GOOS=windows because the Linux/macOS XDG path layout is the
// happy-path covered here; vscodepm/path_test.go covers Windows.
package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

// TestVSCodePMSyncSkipsMissingRootPaths asserts that buildVSCodePMSyncPairs
// drops entries whose rootPath no longer exists on disk and KEEPS entries
// whose rootPath is a real directory.
func TestVSCodePMSyncSkipsMissingRootPaths(t *testing.T) {
	tmp := t.TempDir()
	realDir := filepath.Join(tmp, "real-repo")
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatalf("mkdir real: %v", err)
	}

	entries := []vscodepm.Entry{
		{Name: "real-repo", RootPath: realDir, Tags: []string{"user"}},
		{Name: "ghost", RootPath: filepath.Join(tmp, "does-not-exist"), Tags: []string{"user"}},
		{Name: "blank", RootPath: "", Tags: []string{"user"}},
	}

	pairs, skipped := buildVSCodePMSyncPairs(entries, vscodePMSyncOpts{})

	if got, want := len(pairs), 1; got != want {
		t.Fatalf("pairs: got %d, want %d", got, want)
	}
	if got, want := skipped, 2; got != want {
		t.Fatalf("skipped: got %d, want %d", got, want)
	}
	if pairs[0].RootPath != realDir {
		t.Errorf("kept wrong entry: got %q, want %q", pairs[0].RootPath, realDir)
	}
}

// TestVSCodePMSyncPairsCarryGitmapBrandTag asserts every produced
// pair gets the "gitmap" auto-tag prepended by DetectTagsCustom, so
// Sync's UNION will guarantee branding regardless of pre-existing tags.
func TestVSCodePMSyncPairsCarryGitmapBrandTag(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	entries := []vscodepm.Entry{{Name: "x", RootPath: tmp, Tags: nil}}
	pairs, _ := buildVSCodePMSyncPairs(entries, vscodePMSyncOpts{})

	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if !containsTag(pairs[0].Tags, "gitmap") {
		t.Errorf("brand tag missing from pair: %v", pairs[0].Tags)
	}
}

// TestVSCodePMSyncEndToEndPreservesUserTags writes a projects.json
// containing a hand-crafted "user" tag, runs runVSCodePMSync, and
// asserts the on-disk file still contains the user tag PLUS the
// freshly-detected "gitmap" brand tag.
func TestVSCodePMSyncEndToEndPreservesUserTags(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	repoDir, restore := setupVSCodePMSyncFixture(t)
	defer restore()

	runVSCodePMSync(nil)

	got := loadVSCodePMSyncProjectsJSON(t)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry on disk, got %d", len(got))
	}
	if got[0].RootPath != repoDir {
		t.Fatalf("rootPath mutated: got %q, want %q", got[0].RootPath, repoDir)
	}
	tags := append([]string{}, got[0].Tags...)
	sort.Strings(tags)
	if !containsTag(tags, "user") {
		t.Errorf("user tag stripped: %v", tags)
	}
	if !containsTag(tags, "gitmap") {
		t.Errorf("brand tag missing: %v", tags)
	}
}

// TestVSCodePMSyncDryRunDoesNotMutate ensures `--dry-run` leaves the
// file byte-identical even when re-tagging would change something.
func TestVSCodePMSyncDryRunDoesNotMutate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	_, restore := setupVSCodePMSyncFixture(t)
	defer restore()

	path, err := vscodepm.ProjectsJSONPath()
	if err != nil {
		t.Fatalf("resolve path: %v", err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}

	runVSCodePMSync([]string{"--dry-run"})

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Errorf("dry-run mutated projects.json")
	}
}
