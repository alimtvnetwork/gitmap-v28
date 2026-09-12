package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsEditorTarget(t *testing.T) {
	if !isEditorTarget("code") {
		t.Errorf("expected isEditorTarget(code) = true")
	}

	if !isEditorTarget("vscode") {
		t.Errorf("expected isEditorTarget(vscode) = true")
	}

	if isEditorTarget("notepad") {
		t.Errorf("expected isEditorTarget(notepad) = false")
	}
}

func TestIsWebTarget(t *testing.T) {
	if !isWebTarget("https://github.com") {
		t.Errorf("expected isWebTarget(https://...) = true")
	}

	if !isWebTarget("hotmail.com") {
		t.Errorf("expected isWebTarget(hotmail.com) = true")
	}

	if isWebTarget("nonexistent_directory_without_dot") {
		t.Errorf("expected isWebTarget = false")
	}
}

func TestNormalizeWebUrl(t *testing.T) {
	if got := normalizeWebUrl("hotmail.com"); got != "https://hotmail.com" {
		t.Errorf("expected https://hotmail.com, got %s", got)
	}

	if got := normalizeWebUrl("https://example.com"); got != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", got)
	}
}

func TestResolveOpenTarget(t *testing.T) {
	target, err := resolveOpenTarget([]string{"."})
	if err != nil {
		t.Fatalf("resolveOpenTarget failed: %v", err)
	}

	cwd, _ := os.Getwd()
	if !strings.EqualFold(filepath.Clean(target), filepath.Clean(cwd)) {
		t.Errorf("expected %s, got %s", cwd, target)
	}
}
