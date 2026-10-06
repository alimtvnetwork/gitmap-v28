package cmdide

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseIDEOptions(t *testing.T) {
	args := []string{"--target", "vscode", "--dry-run", "--json", "/some/path"}
	opts, remaining := parseIDEOptions(args)

	if !opts.IsVSCodeTargeted {
		t.Errorf("expected VS Code to be targeted")
	}
	if opts.IsCursorTargeted {
		t.Errorf("expected Cursor not to be targeted")
	}
	if !opts.IsDryRun {
		t.Errorf("expected dry-run to be true")
	}
	if !opts.IsJSON {
		t.Errorf("expected json to be true")
	}
	if len(remaining) != 1 || remaining[0] != "/some/path" {
		t.Errorf("expected remaining args to be ['/some/path'], got %v", remaining)
	}
}

func TestApplyExcludeFilter(t *testing.T) {
	opts := defaultIDEOptions()
	applyExcludeFilter(&opts, "desktop,cursor")

	if !opts.IsVSCodeTargeted {
		t.Errorf("expected VS Code to remain targeted")
	}
	if opts.IsCursorTargeted {
		t.Errorf("expected Cursor to be excluded")
	}
	if opts.IsDesktopTargeted {
		t.Errorf("expected Desktop to be excluded")
	}
	if !opts.IsAntigravityTargeted {
		t.Errorf("expected Antigravity to remain targeted")
	}
}

func TestValidateRepoPath(t *testing.T) {
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("failed to create temp .git dir: %v", err)
	}

	absPath, err := validateRepoPath(tempDir)
	if err != nil {
		t.Errorf("expected valid repo path, got error: %v", err)
	}
	if absPath != tempDir {
		t.Errorf("expected %s, got %s", tempDir, absPath)
	}

	nonRepo := filepath.Join(tempDir, "not-a-repo")
	_ = os.Mkdir(nonRepo, 0755)
	_, nonRepoErr := validateRepoPath(nonRepo)
	if nonRepoErr == nil {
		t.Errorf("expected error for non-git repository path, got nil")
	}
}

func TestDispatchSubcommandHelp(t *testing.T) {
	if err := dispatchSubcommand("help", []string{}); err != nil {
		t.Errorf("expected help to return nil, got %v", err)
	}
}
