package cmdagy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveInjectProjectTarget(t *testing.T) {
	injectProjectFlag = ""
	cwd, _ := filepath.Abs(".")
	res, err := resolveInjectProjectTarget([]string{"test-prompt"})
	if err != nil || res != cwd {
		t.Errorf("expected default to cwd %q, got %q (err: %v)", cwd, res, err)
	}

	injectProjectFlag = ".."
	parent, _ := filepath.Abs("..")
	res, err = resolveInjectProjectTarget([]string{"test-prompt"})
	if err != nil || res != parent {
		t.Errorf("expected flag override %q, got %q (err: %v)", parent, res, err)
	}
	injectProjectFlag = ""
}

func TestLoadInjectPromptPayload(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "sample-prompt.md")
	expectedText := "# Sample Prompt\nPlease test this feature."
	if err := os.WriteFile(tmpFile, []byte(expectedText), 0644); err != nil {
		t.Fatalf("failed to write temp prompt: %v", err)
	}

	content, title, err := loadInjectPromptPayload(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error loading file: %v", err)
	}
	if content != expectedText {
		t.Errorf("expected content %q, got %q", expectedText, content)
	}
	if title != "sample-prompt" {
		t.Errorf("expected title 'sample-prompt', got %q", title)
	}
}

func TestResolveInjectConvID(t *testing.T) {
	injectConvFlag = "custom-conv-12345"
	got := resolveInjectConvID(".")
	if got != "custom-conv-12345" {
		t.Errorf("expected explicit flag 'custom-conv-12345', got %q", got)
	}
	injectConvFlag = ""
}
