//go:build e2e

package e2e_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
)

func TestPromptInject_E2E(t *testing.T) {
	tempProject := t.TempDir()

	promptText := "Analyze repository structure and optimize imports."
	title := "E2E Prompt Test"

	res := cmdagy.InjectAgyPrompt(tempProject, promptText, title, "user_inject", false)
	if !res.IsSuccess && res.Mode == cmdagy.AgyInjectionModeNone {
		t.Fatalf("expected non-none injection mode, got %s", res.Mode)
	}

	stagedPath := filepath.Join(tempProject, ".ai-memory", "prompts", "active-prompt.md")
	if _, err := os.Stat(stagedPath); err == nil {
		content, _ := os.ReadFile(stagedPath)
		if len(content) == 0 {
			t.Errorf("expected staged content at %s", stagedPath)
		}
	}
}

func TestPromptInject_SkipInjectFlag(t *testing.T) {
	tempProject := t.TempDir()

	res := cmdagy.InjectAgyPrompt(tempProject, "Some prompt", "Title", "user_inject", true)
	if res.Mode != cmdagy.AgyInjectionModeNone {
		t.Fatalf("expected ModeNone on isSkipInject=true, got %+v", res)
	}
}

func TestRunAgyPromptInject_CLI(t *testing.T) {
	tempProject := t.TempDir()
	promptFile := filepath.Join(tempProject, "test-prompt.md")
	err := os.WriteFile(promptFile, []byte("Prompt file content"), 0644)
	if err != nil {
		t.Fatalf("failed to write test prompt file: %v", err)
	}

	err = cmdagy.RunAgyPromptInject(nil, []string{promptFile, tempProject})
	if err != nil {
		t.Fatalf("RunAgyPromptInject failed: %v", err)
	}
}
