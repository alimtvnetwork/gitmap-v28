package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdadd/constants"
)

func TestAddRunMissingArg(t *testing.T) {
	err := Run([]string{})
	if err == nil {
		t.Fatal("expected error for empty args")
	}

	if !strings.Contains(err.Error(), constants.ErrCodeAddMissingArg) {
		t.Errorf("expected %s in error, got %v", constants.ErrCodeAddMissingArg, err)
	}
}

func TestAddRunUnknownTarget(t *testing.T) {
	err := Run([]string{"unknown-target"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}

	if !strings.Contains(err.Error(), constants.ErrCodeAddUnknownTarget) {
		t.Errorf("expected %s in error, got %v", constants.ErrCodeAddUnknownTarget, err)
	}
}

func TestAddRunCommonAttr(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	err := Run([]string{constants.SubCmdCommonAttr})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, readErr := os.ReadFile(filepath.Join(dir, constants.FileGitAttributes))
	if readErr != nil {
		t.Fatalf("failed reading .gitattributes: %v", readErr)
	}

	if string(data) != constants.TemplateCommonAttr {
		t.Errorf("template mismatch, got %q", string(data))
	}
}

func TestAddRunCommonIgnore(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	err := Run([]string{constants.SubCmdCommonIgnore})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, readErr := os.ReadFile(filepath.Join(dir, constants.FileGitIgnore))
	if readErr != nil {
		t.Fatalf("failed reading .gitignore: %v", readErr)
	}

	if string(data) != constants.TemplateCommonIgnore {
		t.Errorf("template mismatch, got %q", string(data))
	}
}
