package e2e

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

// readBlobAtHead returns the contents of `rel` at refs/heads/main as
// recorded in the git object store (NOT the working tree). Used by
// edge-case tests that need to assert which side of a force-merge won.
func readBlobAtHead(t *testing.T, r *Repo, rel string) string {
	t.Helper()
	cmd := exec.Command("git", "show", "HEAD:"+rel)
	cmd.Dir = r.Path
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git show HEAD:%s: %v", rel, err)
	}
	// Git appends a trailing newline only if the blob has one; preserve.
	return string(out)
}

// isLinux reports the current GOOS without exposing the runtime
// import to test files that don't need it.
func isLinux() bool { return strings.EqualFold(runtime.GOOS, "linux") }
