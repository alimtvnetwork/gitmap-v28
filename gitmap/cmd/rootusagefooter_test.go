package cmd

import (
	"strings"
	"testing"
)

func TestCaptureGitRejectsEmptyDir(t *testing.T) {
	if got := captureGit("", "rev-parse", "--show-toplevel"); got != "" {
		t.Fatalf("captureGit empty dir = %q, want empty", got)
	}
}

func TestEmitIdentityRowsUsesBuildOverridesWithoutDir(t *testing.T) {
	out := captureStdoutForTest(t, func() {
		emitIdentityRows(IdentityRowParams{
			Dir:            "",
			RepoOverride:   "https://example.com/owner/gitmap-v28",
			BranchOverride: "main",
			ShaOverride:    "abc123",
		})
	})
	for _, want := range []string{"https://example.com/owner/gitmap-v28", "main", "abc123"} {
		if !strings.Contains(out, want) {
			t.Fatalf("footer missing %q in %q", want, out)
		}
	}
}
