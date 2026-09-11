package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
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
