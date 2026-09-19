package cmdssh

import (
	"os"
	"strings"
	"testing"
)

func TestExpandUniversalPath_Tilde(t *testing.T) {
	home, _ := os.UserHomeDir()
	res := ExpandUniversalPath("~/myfolder", "linux")
	if !strings.HasPrefix(res, strings.ReplaceAll(home, "\\", "/")) {
		t.Errorf("expected path to start with home %q, got %q", home, res)
	}
}

func TestExpandUniversalPath_WinMacros(t *testing.T) {
	res := ExpandUniversalPath("%win%\\System32", "windows")
	if !strings.Contains(strings.ToLower(res), "windows\\system32") {
		t.Errorf("expected C:\\Windows\\System32, got %q", res)
	}

	driveRes := ExpandUniversalPath("%win-drive%\\data", "windows")
	if !strings.HasPrefix(driveRes, "C:\\") {
		t.Errorf("expected C:\\data, got %q", driveRes)
	}
}

func TestExpandUniversalPath_TempAndAppData(t *testing.T) {
	tempRes := ExpandUniversalPath("%temp%/file.txt", "linux")
	if strings.Contains(tempRes, "%temp%") {
		t.Errorf("temp macro was not expanded: %q", tempRes)
	}
}
