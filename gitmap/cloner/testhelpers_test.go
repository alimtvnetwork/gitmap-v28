package cloner

// Test-only helpers shared by concurrent_test.go (and any future tests
// in this package that need to spin up a real local git repo on disk).
// Lives in its own file so the helpers don't accidentally leak into the
// production binary and so adding a second consumer doesn't churn the
// big test file.

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// runGit invokes `git` synchronously and fails the test on a non-zero
// exit. When workdir is non-empty, the command runs inside that
// directory (equivalent to `git -C workdir ...`); when empty, the
// caller is expected to put the target path in args (e.g.
// `git init -b main /tmp/foo`). Combined stdout+stderr is attached to
// any failure so CI logs explain WHY git rejected the call.
func runGit(t *testing.T, workdir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	if workdir != "" {
		cmd.Dir = workdir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v (workdir=%q) failed: %v\n--- output ---\n%s",
			args, workdir, err, string(out))
	}
}

// writeFile writes content to path, creating the parent directory tree
// as needed. Mirrors the convenience signature used in other gitmap
// test packages so call sites stay symmetric. 0o644 keeps the bytes
// readable by the user running `go test`; 0o755 on the parents is the
// usual test-fixture default.
func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdirall %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
