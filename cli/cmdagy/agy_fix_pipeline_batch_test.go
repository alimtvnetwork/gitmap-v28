package cmdagy

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAgyFixArgs_BatchFlagsCounts(t *testing.T) {
	args := []string{"fix", "--all", "--projects", "5", "--limit", "2"}
	opts := parseAgyFixArgs(args)
	if !opts.IsAll || opts.ProjectsCount != 5 || opts.Limit != 2 {
		t.Fatalf("unexpected batch counts: %+v", opts)
	}
}

func TestParseAgyFixArgs_BatchToggles(t *testing.T) {
	args := []string{"fix", "--reset-batch", "--no-inject"}
	opts := parseAgyFixArgs(args)
	if !opts.IsResetBatch || !opts.IsNoInject {
		t.Fatalf("unexpected toggles: %+v", opts)
	}
}

func testSaveAndReloadCursor(t *testing.T, cursorPath string) {
	initial := LoadPipelineFixBatchCursor(cursorPath, 3)
	initial.LastIndex = 3
	initial.ProcessedRepos = []string{"repo1", "repo2", "repo3"}
	if err := SavePipelineFixBatchCursor(cursorPath, initial); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	reloaded := LoadPipelineFixBatchCursor(cursorPath, 3)
	if reloaded.LastIndex != 3 || len(reloaded.ProcessedRepos) != 3 {
		t.Fatalf("unexpected reloaded cursor: %+v", reloaded)
	}
}

func TestBatchCursor_SaveAndReset(t *testing.T) {
	cursorPath := filepath.Join(t.TempDir(), "batch-cursor.json")
	testSaveAndReloadCursor(t, cursorPath)
	if err := ResetPipelineFixBatchCursor(cursorPath); err != nil {
		t.Fatalf("unexpected reset error: %v", err)
	}

	reset := LoadPipelineFixBatchCursor(cursorPath, 3)
	if reset.LastIndex != 0 {
		t.Fatalf("expected LastIndex = 0 after reset, got %d", reset.LastIndex)
	}
}

func TestInjectAgyFixTask_SkipInject(t *testing.T) {
	ok, msg := InjectAgyFixTask(".", "some/path", true)
	if ok || !strings.Contains(msg, "skipped by flag") {
		t.Fatalf("expected skipped by flag, got ok=%v msg=%s", ok, msg)
	}
}

func TestToAbsPath(t *testing.T) {
	rel := ".ai-memory/temp/active-agy-pipeline-fix-prompt.txt"
	abs := toAbsPath(rel)
	if !filepath.IsAbs(abs) {
		t.Fatalf("expected absolute path, got %s", abs)
	}
}
