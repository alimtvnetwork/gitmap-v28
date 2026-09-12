package heavy_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func TestExecute_WithGitmapCdAndRelativeCd(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subproject")
	_ = os.MkdirAll(subDir, 0755)

	m := buildGitmapCdMacro(tmpDir)
	if err := macro.Execute(context.Background(), m, macro.ExecOptions{DryRun: false}); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func buildGitmapCdMacro(tmpDir string) *macro.Macro {
	return &macro.Macro{
		Name: "test-gitmap-cd-macro",
		Steps: []macro.MacroStep{
			{StepNum: 1, CommandLine: "cd " + tmpDir},
			{StepNum: 2, CommandLine: "gitmap cd subproject"},
			{StepNum: 3, CommandLine: "cd .."},
			{StepNum: 4, CommandLine: "cd -"},
		},
	}
}
