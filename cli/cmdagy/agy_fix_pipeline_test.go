// Package cmdagy — agy_fix_pipeline_test.go tests the fix-pipeline prompt assembly and payload formatting.
package cmdagy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssembleFixPipelinePayload_BothNonEmpty(t *testing.T) {
	logs := "ERROR: step 1 failed with exit code 1"
	prompt := "# Fix Prompt Instructions\nN = 200"
	got := AssembleFixPipelinePayload(logs, prompt)
	expected := logs + "\n\n" + prompt
	if got != expected {
		t.Fatalf("expected payload with two-line gap, got:\n%q", got)
	}
}

func TestAssembleFixPipelinePayload_TrimsSurroundingWhitespace(t *testing.T) {
	logs := "  \nERROR: compile error\n\n  "
	prompt := "\n  # Fix Prompt Instructions  \n"
	got := AssembleFixPipelinePayload(logs, prompt)
	expected := "ERROR: compile error\n\n# Fix Prompt Instructions"
	if got != expected {
		t.Fatalf("expected whitespace trimmed payload, got:\n%q", got)
	}
}

func TestAssembleFixPipelinePayload_EmptyLogsOrPrompt(t *testing.T) {
	promptOnly := AssembleFixPipelinePayload("", "# Only Prompt")
	if promptOnly != "# Only Prompt" {
		t.Fatalf("expected prompt only, got %q", promptOnly)
	}

	logsOnly := AssembleFixPipelinePayload("Only Logs", "")
	if logsOnly != "Only Logs" {
		t.Fatalf("expected logs only, got %q", logsOnly)
	}
}

func TestResolveTargetRepoArg(t *testing.T) {
	if resolveTargetRepoArg([]string{"alimtvnetwork/gitmap-v28"}) != "alimtvnetwork/gitmap-v28" {
		t.Fatalf("expected target repo arg to resolve")
	}

	if resolveTargetRepoArg([]string{"--detailed"}) != "" {
		t.Fatalf("expected flag to not be resolved as repo")
	}

	if resolveTargetRepoArg([]string{}) != "" {
		t.Fatalf("expected empty args to resolve to empty string")
	}
}

func TestSelectPromptPath(t *testing.T) {
	withRelease := selectPromptPath(false)
	if !strings.Contains(withRelease, "04-ci-cd-fix-with-release.md") {
		t.Fatalf("expected 04-ci-cd-fix-with-release.md, got %s", withRelease)
	}

	noRelease := selectPromptPath(true)
	if !strings.Contains(noRelease, "01-ci-cd-fix.md") {
		t.Fatalf("expected 01-ci-cd-fix.md, got %s", noRelease)
	}
}

func TestLoadCicdFixPrompt_CustomPath(t *testing.T) {
	tempCustom := filepath.Join(t.TempDir(), "custom_prompt.md")
	_ = os.WriteFile(tempCustom, []byte("# Custom Fix Prompt"), 0644)
	content, source := loadCicdFixPrompt(tempCustom, false)
	if content != "# Custom Fix Prompt" || source != tempCustom {
		t.Fatalf("expected custom prompt loaded, got content=%q source=%q", content, source)
	}
}
