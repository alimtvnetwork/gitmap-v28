package mapper

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/scanner"
)

// TestBuildRecordsWithRoot_PinsBaseToProvidedRoot verifies that passing a
// non-empty relRoot rewrites RelativePath against that root regardless of
// what the scanner originally computed. This is the cwd-stable contract
// for the `--relative-root` flag.
func TestBuildRecordsWithRoot_PinsBaseToProvidedRoot(t *testing.T) {
	root := absRoot(t)
	repos := []scanner.RepoInfo{
		{
			AbsolutePath: filepath.Join(root, "team-a", "repo-1"),
			RelativePath: "should-be-overridden",
		},
	}
	records := BuildRecordsWithRoot(repos, "https", "", root)
	if len(records) != 1 {
		t.Fatalf("want 1 record, got %d", len(records))
	}
	want := filepath.Join("team-a", "repo-1")
	if records[0].RelativePath != want {
		t.Errorf("RelativePath = %q, want %q", records[0].RelativePath, want)
	}
}

// TestBuildRecordsWithRoot_EmptyRootKeepsScannerValue confirms that the
// thin BuildRecords wrapper (relRoot == "") is a no-op rewrite — the
// scanner's RelativePath flows through verbatim.
func TestBuildRecordsWithRoot_EmptyRootKeepsScannerValue(t *testing.T) {
	repos := []scanner.RepoInfo{
		{
			AbsolutePath: filepath.Join(absRoot(t), "x"),
			RelativePath: "kept-as-is",
		},
	}
	records := BuildRecordsWithRoot(repos, "https", "", "")
	if records[0].RelativePath != "kept-as-is" {
		t.Errorf("RelativePath = %q, want kept-as-is", records[0].RelativePath)
	}
}

// TestBuildRecordsWithRoot_OutOfTreeFallsBack ensures repos that live
// outside relRoot do NOT silently produce "../"-prefixed paths in the
// output. The mapper falls back to the scanner-computed RelativePath
// for that single row instead of corrupting downstream clone scripts.
func TestBuildRecordsWithRoot_OutOfTreeFallsBack(t *testing.T) {
	root := filepath.Join(absRoot(t), "inside")
	outside := filepath.Join(absRoot(t), "outside", "repo")
	repos := []scanner.RepoInfo{
		{
			AbsolutePath: outside,
			RelativePath: "scanner-fallback",
		},
	}
	records := BuildRecordsWithRoot(repos, "https", "", root)
	if records[0].RelativePath != "scanner-fallback" {
		t.Errorf("expected fallback to scanner value, got %q", records[0].RelativePath)
	}
}

// absRoot returns a platform-appropriate absolute root for synthesizing
// repo paths in tests without touching the filesystem.
func absRoot(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		return `C:\repos`
	}
	if !strings.HasPrefix("/tmp", "/") {
		t.Fatal("unexpected: /tmp is not absolute")
	}

	return "/tmp/gitmap-relroot-test"
}
