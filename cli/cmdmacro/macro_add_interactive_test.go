package cmdmacro

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func TestSanitizeRawEscapeCodes(t *testing.T) {
	input := "^[[A^[[Agit status\x1b[B"
	expected := "git status"
	got := sanitizeRawEscapeCodes(input)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestIsMkdirCmd(t *testing.T) {
	cases := []struct {
		cmd      string
		expected bool
	}{
		{"mkdir foo", true},
		{":mkdir bar", true},
		{"gitmap mkdir baz", true},
		{"mkdir", true},
		{":mkdir", true},
		{"ls", false},
		{"git status", false},
	}

	for _, c := range cases {
		if got := isMkdirCmd(c.cmd); got != c.expected {
			t.Errorf("isMkdirCmd(%q) = %v; want %v", c.cmd, got, c.expected)
		}
	}
}

func TestExtractMkdirTarget(t *testing.T) {
	cases := []struct {
		parts    []string
		expected string
	}{
		{[]string{"mkdir", "foo"}, "foo"},
		{[]string{":mkdir", "-p", "bar/baz"}, "bar/baz"},
		{[]string{"gitmap", "mkdir", "myfolder"}, "myfolder"},
	}

	for _, c := range cases {
		if got := extractMkdirTarget(c.parts); got != c.expected {
			t.Errorf("extractMkdirTarget(%v) = %q; want %q", c.parts, got, c.expected)
		}
	}
}

func TestExecuteInteractiveMkdir(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "interactive_created_dir")

	success := executeInteractiveMkdir(target, "mkdir "+target)
	if !success {
		t.Fatalf("expected executeInteractiveMkdir to succeed")
	}

	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		t.Fatalf("expected directory to exist: %s", target)
	}

	tempSub := "//temp/interactive_test_sub_" + filepath.Base(tempDir)
	if !executeInteractiveMkdir(tempSub, "mkdir "+tempSub) {
		t.Fatalf("expected executeInteractiveMkdir with //temp to succeed")
	}

	defer os.RemoveAll(macro.ExpandPathAndEnv(tempSub))
}

func TestIsRemoteSSHSession(t *testing.T) {
	t.Setenv("SSH_CLIENT", "")
	t.Setenv("SSH_TTY", "")
	t.Setenv("SSH_CONNECTION", "")
	if isRemoteSSHSession() {
		t.Fatalf("expected isRemoteSSHSession to be false when env unset")
	}

	t.Setenv("SSH_CLIENT", "192.168.1.50 54321 22")
	if !isRemoteSSHSession() {
		t.Fatalf("expected isRemoteSSHSession to be true when SSH_CLIENT set")
	}
}

func TestHandleZeroPipedSteps_RemoteSSH(t *testing.T) {
	t.Setenv("SSH_CLIENT", "10.0.0.1 12345 22")
	err := handleZeroPipedSteps()
	if err == nil {
		t.Fatalf("expected error when zero piped steps in remote SSH session")
	}
}

func TestHandleZeroInteractiveSteps_RemoteSSH(t *testing.T) {
	t.Setenv("SSH_CONNECTION", "10.0.0.1 12345 10.0.0.2 22")
	err := handleZeroInteractiveSteps("test-macro")
	if err == nil {
		t.Fatalf("expected error when zero interactive steps in remote SSH session")
	}
}

func TestIsCpFileHelperCmd(t *testing.T) {
	cases := []struct {
		cmd      string
		expected bool
	}{
		{"cp ../git-work .", true},
		{":cp foo bar", true},
		{"copy file.txt dest/", true},
		{":copy a b", true},
		{"cpfile src dst", true},
		{"cp", true},
		{"copy", true},
		{"ls", false},
		{"cat foo", false},
	}

	for _, c := range cases {
		if got := isCpFileHelperCmd(c.cmd); got != c.expected {
			t.Errorf("isCpFileHelperCmd(%q) = %v; want %v", c.cmd, got, c.expected)
		}
	}
}
