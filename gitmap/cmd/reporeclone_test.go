package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSplitRepoRecloneArgs covers every spelling of the -y flag and
// ensures unknown tokens pass through as positionals.
func TestSplitRepoRecloneArgs(t *testing.T) {
	cases := []struct {
		name       string
		in         []string
		wantYes    bool
		wantPosLen int
		wantFirst  string
	}{
		{"empty", nil, false, 0, ""},
		{"single -y", []string{"-y"}, true, 0, ""},
		{"long --yes", []string{"--yes"}, true, 0, ""},
		{"path passthrough", []string{"./repo"}, false, 1, "./repo"},
		{"yes plus path", []string{"-y", "./repo"}, true, 1, "./repo"},
		{"manifest path passthrough", []string{".gitmap/output/gitmap.json"}, false, 1, ".gitmap/output/gitmap.json"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotYes, gotPos := splitRepoRecloneArgs(tc.in)
			if gotYes != tc.wantYes {
				t.Fatalf("yes: got %v want %v", gotYes, tc.wantYes)
			}
			if len(gotPos) != tc.wantPosLen {
				t.Fatalf("positionals: got %d want %d (%#v)", len(gotPos), tc.wantPosLen, gotPos)
			}
			if tc.wantPosLen > 0 && gotPos[0] != tc.wantFirst {
				t.Fatalf("first positional: got %q want %q", gotPos[0], tc.wantFirst)
			}
		})
	}
}

// TestResolveRepoRecloneTarget verifies the shape-detection that
// gates the overlay: a non-git directory MUST NOT be claimed, and a
// real `.git` directory MUST be claimed.
func TestResolveRepoRecloneTarget(t *testing.T) {
	tmp := t.TempDir()
	notGit := filepath.Join(tmp, "plain")
	if err := os.MkdirAll(notGit, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRepo := filepath.Join(tmp, "real")
	if err := os.MkdirAll(filepath.Join(gitRepo, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	if _, ok := resolveRepoRecloneTarget([]string{notGit}); ok {
		t.Fatal("non-git dir must not be claimed")
	}
	got, ok := resolveRepoRecloneTarget([]string{gitRepo})
	if !ok {
		t.Fatal("git dir must be claimed")
	}
	if abs, _ := filepath.Abs(gitRepo); got != abs {
		t.Fatalf("target: got %q want %q", got, abs)
	}
	if _, ok := resolveRepoRecloneTarget([]string{notGit, gitRepo}); ok {
		t.Fatal("multi-positional must not be claimed (manifest pipeline owns it)")
	}
	if _, ok := resolveRepoRecloneTarget([]string{"/nonexistent/path/xyz123"}); ok {
		t.Fatal("missing path must not be claimed")
	}
}

// TestIsGitRepoDir is a small belt-and-braces guard so a future
// refactor of the helper can't silently change the trigger shape.
func TestIsGitRepoDirHelper(t *testing.T) {
	tmp := t.TempDir()
	if isGitRepoDir(tmp) {
		t.Fatal("empty dir is not a git repo")
	}
	if err := os.MkdirAll(filepath.Join(tmp, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !isGitRepoDir(tmp) {
		t.Fatal("dir with .git/ must be detected")
	}
}

// TestTryRunRepoRecloneFallthrough is the explicit non-regression
// guard for the v6.5.0 overlay: any arg shape that the manifest
// pipeline owns MUST cause tryRunRepoReclone to return false
// WITHOUT touching disk. Documented in the spec as the
// "manifest behavior is unchanged" promise — this test makes that
// promise executable.
func TestTryRunRepoRecloneFallthrough(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(originalDir) }()
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	cases := [][]string{
		{".gitmap/output/gitmap.json"},               // manifest path, no .git
		{"--manifest", ".gitmap/output/gitmap.json"}, // flag form
		{"file1.json", "file2.json"},                 // multi-positional
		{"/definitely/not/a/real/path/xyz123"},       // missing path
	}
	for _, args := range cases {
		if tryRunRepoReclone(args) {
			t.Fatalf("manifest-shaped args must fall through: %#v", args)
		}
	}
}
