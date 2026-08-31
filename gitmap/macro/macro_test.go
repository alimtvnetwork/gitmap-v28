package macro

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestMacroSaveLoadListDelete(t *testing.T) {
	m := &Macro{
		Name: "test-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "echo hello", TimeoutSeconds: 10},
			{StepNum: 2, CommandLine: "echo world", TimeoutSeconds: 10},
		},
	}

	if err := SaveMacro(m); err != nil {
		t.Fatalf("SaveMacro failed: %v", err)
	}

	loaded, err := LoadMacro("test-macro")
	if err != nil {
		t.Fatalf("LoadMacro failed: %v", err)
	}
	if loaded.Name != "test-macro" || len(loaded.Steps) != 2 {
		t.Fatalf("Loaded macro mismatch: %+v", loaded)
	}

	list, err := ListMacros()
	if err != nil {
		t.Fatalf("ListMacros failed: %v", err)
	}
	found := false
	for _, item := range list {
		if item.Name == "test-macro" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("test-macro not found in list: %+v", list)
	}

	// Dry run execute
	err = Execute(context.Background(), loaded, ExecOptions{DryRun: true})
	if err != nil {
		t.Fatalf("Dry run execution failed: %v", err)
	}

	if err := DeleteMacro("test-macro"); err != nil {
		t.Fatalf("DeleteMacro failed: %v", err)
	}
}

func TestExecute_WithCdAndEnvExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("TEST_EXEC_DIR", tmpDir)
	defer os.Unsetenv("TEST_EXEC_DIR")

	m := &Macro{
		Name: "test-cd-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "cd %TEST_EXEC_DIR%"},
			{StepNum: 2, CommandLine: "echo active"},
		},
	}

	err := Execute(context.Background(), m, ExecOptions{DryRun: false})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestExecute_WithGitmapCdAndRelativeCd(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subproject")
	_ = os.MkdirAll(subDir, 0755)

	m := &Macro{
		Name: "test-gitmap-cd-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "cd " + tmpDir},
			{StepNum: 2, CommandLine: "gitmap cd subproject"},
			{StepNum: 3, CommandLine: "cd .."},
			{StepNum: 4, CommandLine: "cd -"},
		},
	}

	err := Execute(context.Background(), m, ExecOptions{DryRun: false})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestExecute_WithJSONAndFileReport(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "report.json")

	m := &Macro{
		Name: "test-json-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "echo json-step"},
		},
	}

	err := Execute(context.Background(), m, ExecOptions{
		JSON:     true,
		FilePath: outFile,
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	content, readErr := os.ReadFile(outFile)
	if readErr != nil || len(content) == 0 {
		t.Fatalf("Report file not written: %v", readErr)
	}
}

func TestExecute_WithYAMLAndFileReport(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "report.yaml")

	m := &Macro{
		Name: "test-yaml-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "echo yaml-step"},
		},
	}

	err := Execute(context.Background(), m, ExecOptions{
		YAML:     true,
		FilePath: outFile,
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	content, readErr := os.ReadFile(outFile)
	if readErr != nil || len(content) == 0 {
		t.Fatalf("Report file not written: %v", readErr)
	}
}
