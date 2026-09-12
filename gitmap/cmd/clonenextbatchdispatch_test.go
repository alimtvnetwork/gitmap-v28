package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsBatchRunnable_ExplicitFlagsWin asserts that the explicit
// --csv / --all flags short-circuit the implicit cwd check, so the
// implicit logic can never accidentally OVERRIDE a user-requested
// mode.
func TestIsBatchRunnable_ExplicitFlagsWin(t *testing.T) {
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
			isNonRunBatch := !isBatchRunnable(tc.flags, "")
			if isNonRunBatch {
				t.Fatalf("isBatchRunnable(%+v, \"\") = false, want true", tc.flags)
			}
		})
	}
}

// TestIsBatchRunnable_ImplicitTrigger_FiresOnScanRoot creates a fixture
// directory that is NOT a git repo but contains one git subdirectory,
// then asserts the dispatcher recognizes it as a scan root.
func TestIsBatchRunnable_ImplicitTrigger_FiresOnScanRoot(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo-a")
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	isNonRunBatch := !isBatchRunnable(CloneNextFlags{}, root)
	if isNonRunBatch {
		t.Fatalf("expected implicit batch trigger to fire on scan root %s", root)
	}
}

// TestIsBatchRunnable_ImplicitTrigger_SkipsInsideRepo asserts the
// dispatcher does NOT promote a single-repo invocation to batch mode
// when the user is sitting inside a real git repo. This is the
// regression we're guarding: clobbering the single-repo path here
// would silently change the cn semantics.
func TestIsBatchRunnable_ImplicitTrigger_SkipsInsideRepo(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if isBatchRunnable(CloneNextFlags{}, repo) {
		t.Fatalf("expected single-repo path on git cwd %s, got batch", repo)
	}
}

// TestIsBatchRunnable_ImplicitTrigger_SkipsEmptyDir asserts the
// dispatcher falls through to the single-repo path (which then prints
// a clean "no remote" error) when cwd is neither a repo nor a scan
// root. We don't want to surprise users in random directories.
func TestIsBatchRunnable_ImplicitTrigger_SkipsEmptyDir(t *testing.T) {
	dir := t.TempDir()

	if isBatchRunnable(CloneNextFlags{}, dir) {
		t.Fatalf("expected single-repo path on empty dir %s, got batch", dir)
	}
}

// TestIsBatchRunnable_ImplicitTrigger_IgnoresNonRepoSubdirs verifies
// that plain directories one level down don't trip the trigger — only
// directories with their own .git entry count.
func TestIsBatchRunnable_ImplicitTrigger_IgnoresNonRepoSubdirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "not-a-repo"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if isBatchRunnable(CloneNextFlags{}, root) {
		t.Fatalf("expected single-repo path on non-repo subdirs, got batch")
	}
}
