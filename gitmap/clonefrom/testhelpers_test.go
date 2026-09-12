package clonefrom

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// writeFile is a tiny os.WriteFile wrapper with a fixed permission.
func writeFile(path, body string) error {
	return os.WriteFile(path, []byte(body), 0o644)
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not on PATH: %v", err)
	}
}

func makeBareRepo(t *testing.T) string {
	t.Helper()
	work := t.TempDir()
	bare := filepath.Join(t.TempDir(), "src.git")

	runGit(t, work, "init", "-q")
	runGit(t, work, "config", "user.email", "t@e")
	runGit(t, work, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(work, "README"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	runGit(t, work, "add", ".")
	runGit(t, work, "commit", "-q", "-m", "init")
	runGit(t, work, "clone", "--bare", "-q", work, bare)

	return bare
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, string(out))
	}
}
