package cmdcursor

import (
	"os"
	"path/filepath"
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

func TestRequireCursorSettingsFileMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "settings.json")
	err := requireCursorSettingsFile(missing)
	if err == nil {
		t.Fatal("expected error for missing settings file, got nil")
	}

	if !strings.Contains(err.Error(), missing) {
		t.Errorf("expected error to mention missing path %s, got %v", missing, err)
	}

	if !strings.Contains(err.Error(), "gitmap cursor settings apply") {
		t.Errorf("expected error to suggest remediation, got %v", err)
	}
}

func TestRequireCursorSettingsFilePresent(t *testing.T) {
	present := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(present, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create temp settings file: %v", err)
	}

	if err := requireCursorSettingsFile(present); err != nil {
		t.Fatalf("expected nil for present settings file, got %v", err)
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
