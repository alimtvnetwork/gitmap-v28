package cmd

import (
	"bytes"
	"io"
	"os"
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
		emitIdentityRows("", "https://example.com/owner/gitmap-v27", "main", "abc123")
	})
	for _, want := range []string{"https://example.com/owner/gitmap-v27", "main", "abc123"} {
		if !strings.Contains(out, want) {
			t.Fatalf("footer missing %q in %q", want, out)
		}
	}
}

func captureStdoutForTest(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}
