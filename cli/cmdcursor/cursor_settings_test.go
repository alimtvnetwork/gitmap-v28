package cmdcursor

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetWindowsSystemDriveRoot(t *testing.T) {
	driveRoot := getWindowsSystemDriveRoot()
	if driveRoot == "" {
		t.Fatal("expected non-empty drive root")
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(driveRoot, "\\") && !strings.HasSuffix(driveRoot, "/") {
		t.Fatalf("expected drive root to end with path separator, got %s", driveRoot)
	}
}

func TestGetCursorSettingsPath(t *testing.T) {
	path, err := getCursorSettingsPath()
	if err != nil {
		t.Fatalf("unexpected error resolving cursor settings path: %v", err)
	}
	if path == "" {
		t.Fatal("expected non-empty cursor settings path")
	}
	if !strings.HasSuffix(path, "settings.json") {
		t.Fatalf("expected path to end with settings.json, got %s", path)
	}
}

func TestGetCursorProjectsJSONPath(t *testing.T) {
	path, err := GetCursorProjectsJSONPath()
	if err != nil {
		t.Fatalf("unexpected error resolving cursor projects path: %v", err)
	}
	if path == "" {
		t.Fatal("expected non-empty cursor projects path")
	}
	if !strings.HasSuffix(path, "projects.json") {
		t.Fatalf("expected path to end with projects.json, got %s", path)
	}
}
