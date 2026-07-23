package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestShouldRunBatch_ExplicitFlagsWin asserts that the explicit
// --csv / --all flags short-circuit the implicit cwd check, so the
// implicit logic can never accidentally OVERRIDE a user-requested
// mode.
func TestShouldRunBatch_ExplicitFlagsWin(t *testing.T) {
	cases := []struct {
		name  string
		flags CloneNextFlags
	}{
		{"csv path set", CloneNextFlags{CSVPath: "/some/file.csv"}},
		{"all flag set", CloneNextFlags{All: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Pass an empty cwd to prove the implicit check is bypassed.
			if !shouldRunBatch(tc.flags, "") {
				t.Fatalf("shouldRunBatch(%+v, \"\") = false, want true", tc.flags)
			}
		})
	}
}

// TestShouldRunBatch_ImplicitTrigger_FiresOnScanRoot creates a fixture
// directory that is NOT a git repo but contains one git subdirectory,
// then asserts the dispatcher recognizes it as a scan root.
func TestShouldRunBatch_ImplicitTrigger_FiresOnScanRoot(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo-a")
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if !shouldRunBatch(CloneNextFlags{}, root) {
		t.Fatalf("expected implicit batch trigger to fire on scan root %s", root)
	}
}

// TestShouldRunBatch_ImplicitTrigger_SkipsInsideRepo asserts the
// dispatcher does NOT promote a single-repo invocation to batch mode
// when the user is sitting inside a real git repo. This is the
// regression we're guarding: clobbering the single-repo path here
// would silently change the cn semantics.
func TestShouldRunBatch_ImplicitTrigger_SkipsInsideRepo(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if shouldRunBatch(CloneNextFlags{}, repo) {
		t.Fatalf("expected single-repo path on git cwd %s, got batch", repo)
	}
}

// TestShouldRunBatch_ImplicitTrigger_SkipsEmptyDir asserts the
// dispatcher falls through to the single-repo path (which then prints
// a clean "no remote" error) when cwd is neither a repo nor a scan
// root. We don't want to surprise users in random directories.
func TestShouldRunBatch_ImplicitTrigger_SkipsEmptyDir(t *testing.T) {
	dir := t.TempDir()

	if shouldRunBatch(CloneNextFlags{}, dir) {
		t.Fatalf("expected single-repo path on empty dir %s, got batch", dir)
	}
}

// TestShouldRunBatch_ImplicitTrigger_IgnoresNonRepoSubdirs verifies
// that plain directories one level down don't trip the trigger — only
// directories with their own .git entry count.
func TestShouldRunBatch_ImplicitTrigger_IgnoresNonRepoSubdirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "not-a-repo"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if shouldRunBatch(CloneNextFlags{}, root) {
		t.Fatalf("expected single-repo path on non-repo subdirs, got batch")
	}
}
