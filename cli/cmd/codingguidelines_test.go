package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPatchCGArithmeticIncrements(t *testing.T) {
	t.Parallel()

	in := "((WROTE_NEW++))\n((COPIED++))\n((count + 1))\n"
	got := patchCGArithmeticIncrements(in)
	want := "((WROTE_NEW+=1))\n((COPIED+=1))\n((count + 1))\n"
	isMismatch := got != want
	if isMismatch {
		t.Fatalf("patched script mismatch:\nwant %q\n got %q", want, got)
	}
}

func assertCGNotes(t *testing.T, stderr string) {
	t.Helper()
	for _, want := range []string{"Note: --no-commit set", "Note: --no-push set"} {
		hasNote := strings.Contains(stderr, want)
		if !hasNote {
			t.Fatalf("stderr missing %q: %q", want, stderr)
		}
	}
}

func TestCommitCodingGuidelinesNoCommitNoPushPrintsBothNotes(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	err := CommitCodingGuidelines(CGCommitOpts{IsSkipCommit: true, IsSkipPush: true, Stdout: &stdout, Stderr: &stderr})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	assertCGNotes(t, stderr.String())
}

func TestPatchCGWindowsScriptFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.ps1")
	sample := "Write-Warning \"failed: $oldFile: msg\"\nWrite-Warning \"failed: $destPath: msg\"\nWrite-Warning \"failed: $targetVersionFile: msg\"\n"
	_ = os.WriteFile(path, []byte(sample), 0644)

	err := patchCGWindowsScriptFile(path)
	if err != nil {
		t.Fatalf("patchCGWindowsScriptFile failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		t.Fatalf("expected UTF-8 BOM prefix, got: %x", data[:min(len(data), 3)])
	}

	out := string(data)
	if strings.Contains(out, "$oldFile:") || strings.Contains(out, "$destPath:") || strings.Contains(out, "$targetVersionFile:") {
		t.Fatalf("unpatched syntax remains: %s", out)
	}

	if !strings.Contains(out, "${oldFile}:") || !strings.Contains(out, "${destPath}:") || !strings.Contains(out, "${targetVersionFile}:") {
		t.Fatalf("expected patched variables with braces, got: %s", out)
	}
}
