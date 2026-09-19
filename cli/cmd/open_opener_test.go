package cmd

import (
	"runtime"
	"strings"
	"testing"
)

func TestHandleHeadlessOpen(t *testing.T) {
	err := handleHeadlessOpen("https://google.com")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBuildOpenerCommand(t *testing.T) {
	cmd := buildOpenerCommand("https://google.com")
	if cmd == nil {
		t.Fatalf("expected non-nil command")
	}

	if runtime.GOOS == "windows" {
		assertWindowsOpener(t, cmd.Path, cmd.Args)
	}
}

func assertWindowsOpener(t *testing.T, path string, args []string) {
	isRundll := strings.Contains(strings.ToLower(path), "rundll32") ||
		strings.Contains(strings.ToLower(args[0]), "rundll32")
	if !isRundll {
		t.Errorf("expected rundll32 on windows, got path=%s args=%v", path, args)
	}
}

func TestIsHeadlessLinux(t *testing.T) {
	isHeadless := isHeadlessLinux()
	if runtime.GOOS != "linux" && isHeadless {
		t.Errorf("expected isHeadlessLinux = false on non-linux OS")
	}
}
