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

func TestLongFooterBlockHasRequiredFields(t *testing.T) {
	out := captureStdoutForTest(t, func() {
		printGitmapIdentityBlockLong()
	})

	required := []string{
		"gitmap binary",
		"● Name:",
		"gitmap",
		"● Git URL:",
		"● Version:",
		"● Commit SHA:",
		"● Database:",
		"● Installed path:",
	}

	for _, want := range required {
		if !strings.Contains(out, want) {
			t.Errorf("long footer missing field %q in output:\n%s", want, out)
		}
	}
}

func TestShortFooterBlockHasRequiredFields(t *testing.T) {
	out := captureStdoutForTest(t, func() {
		printGitmapIdentityBlockShort()
	})

	if !strings.Contains(out, "● Version:") {
		t.Errorf("short footer missing Version in output:\n%s", out)
	}

	if !strings.Contains(out, "● Commit SHA:") {
		t.Errorf("short footer missing Commit SHA in output:\n%s", out)
	}

	if strings.Contains(out, "● Installed path:") {
		t.Errorf("short footer should not contain Installed path:\n%s", out)
	}
}
