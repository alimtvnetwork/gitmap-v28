package cmdadd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdadd/constants"
)

func TestRunAddHelp(t *testing.T) {
	testRunAddHelpHookInvoked(t)
	testRunAddHelpNilSafe(t)
}

func testRunAddHelpHookInvoked(t *testing.T) {
	origHook := CheckHelpFn
	defer func() { CheckHelpFn = origHook }()

	hasCalledHelp := false
	CheckHelpFn = func(cmd string, args []string) {
		if cmd == constants.CmdAdd {
			hasCalledHelp = true
		}
	}

	_ = RunAdd([]string{"--help"})
	if !hasCalledHelp {
		t.Error("expected CheckHelpFn to be invoked")
	}
}

func testRunAddHelpNilSafe(t *testing.T) {
	origHook := CheckHelpFn
	defer func() { CheckHelpFn = origHook }()

	CheckHelpFn = nil
	checkHelp(constants.CmdAdd, []string{"--help"})
}

func TestRunAddMissingArg(t *testing.T) {
	err := RunAdd([]string{})
	if err == nil {
		t.Fatal("expected error for empty args")
	}

	if !strings.Contains(err.Error(), constants.ErrCodeAddMissingArg) {
		t.Errorf("expected %s in error, got %v", constants.ErrCodeAddMissingArg, err)
	}
}

func TestRunAddUnknownTarget(t *testing.T) {
	err := RunAdd([]string{"unknown-target"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}

	if !strings.Contains(err.Error(), constants.ErrCodeAddUnknownTarget) {
		t.Errorf("expected %s in error, got %v", constants.ErrCodeAddUnknownTarget, err)
	}
}

func TestRunAddCommonAttr(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	err := RunAdd([]string{constants.SubCmdCommonAttr})
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

func TestRunAddCommonIgnore(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	err := RunAdd([]string{constants.SubCmdCommonIgnore})
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
