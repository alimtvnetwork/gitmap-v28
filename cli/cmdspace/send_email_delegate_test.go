package cmdspace

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveGoEmailReaderBinEmpty(t *testing.T) {
	empty := t.TempDir()
	t.Setenv("PATH", empty)
	t.Setenv("HOME", empty)

	if got := resolveGoEmailReaderBin(); len(got) > 0 {
		t.Fatalf("expected empty resolution, got %q", got)
	}
}

// stubGoEmailReader installs an executable stub that records its argv.
func stubGoEmailReader(t *testing.T) string {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("skipping shell script stub test on windows")
	}

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if mkErr := os.Mkdir(binDir, 0755); mkErr != nil {
		t.Fatal(mkErr)
	}
	out := filepath.Join(dir, "argv.txt")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"" + out + "\"\n"
	stub := filepath.Join(binDir, "go-email-reader")
	if writeErr := os.WriteFile(stub, []byte(script), 0755); writeErr != nil {
		t.Fatal(writeErr)
	}
	t.Setenv("PATH", binDir)

	return out
}

func TestDelegateSendEmailArgv(t *testing.T) {
	out := stubGoEmailReader(t)

	args := []string{"--to", "a@x.com", "--subject", "Hi", "--body", "Hello", "--html"}
	if err := RunSendEmail(args); err != nil {
		t.Fatal(err)
	}

	raw, readErr := os.ReadFile(out)
	if readErr != nil {
		t.Fatal(readErr)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	want := []string{"send", "--to", "a@x.com", "--subject", "Hi", "--body", "Hello", "--html"}
	if len(lines) != len(want) {
		t.Fatalf("argv wrong: %v", lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestDelegateSendEmailArgvFile(t *testing.T) {
	out := stubGoEmailReader(t)

	dir := t.TempDir()
	yml := writeTempEmailFile(t, dir, "m.yaml", "to: file@x.com\nsubject: FileSubj\n")
	if err := RunSendEmail([]string{"--file", yml}); err != nil {
		t.Fatal(err)
	}

	raw, readErr := os.ReadFile(out)
	if readErr != nil {
		t.Fatal(readErr)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	// WHY subject included (31-Oct-26): delegateSendEmailArgv always
	// passes the merged subject, and the file supplies FileSubj, so
	// the argv carries --subject FileSubj before --file. The old
	// expectation omitted it and failed locally and in full-suite runs.
	want := []string{"send", "--to", "file@x.com", "--subject", "FileSubj", "--file", yml}
	if len(lines) != len(want) {
		t.Fatalf("argv wrong: %v", lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}
