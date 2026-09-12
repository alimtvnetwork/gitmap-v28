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

	testSaveAndLoadMacro(t, m)
	assertMacroListed(t, "test-macro")
	testDryRunAndCleanup(t, m)
}

func testSaveAndLoadMacro(t *testing.T, m *Macro) {
	if err := SaveMacro(m); err != nil {
		t.Fatalf("SaveMacro failed: %v", err)
	}

	loaded, err := LoadMacro(m.Name)
	if err != nil || loaded.Name != m.Name || len(loaded.Steps) != len(m.Steps) {
		t.Fatalf("LoadMacro failed: %v, loaded: %+v", err, loaded)
	}
}

func assertMacroListed(t *testing.T, name string) {
	list, err := ListMacros()
	if err != nil {
		t.Fatalf("ListMacros failed: %v", err)
	}

	for _, item := range list {
		if item.Name == name {
			return
		}
	}

	t.Fatalf("macro %q not found in list", name)
}

func testDryRunAndCleanup(t *testing.T, m *Macro) {
	if err := Execute(context.Background(), m, ExecOptions{DryRun: true}); err != nil {
		t.Fatalf("Dry run execution failed: %v", err)
	}

	if err := DeleteMacro(m.Name); err != nil {
		t.Fatalf("DeleteMacro failed: %v", err)
	}

	if err := DeleteMacro(m.Name); err != nil {
		t.Fatalf("Idempotent delete failed: %v", err)
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

	if err := Execute(context.Background(), m, ExecOptions{DryRun: false}); err != nil {
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

	opts := ExecOptions{JSON: true, FilePath: outFile}
	runAndAssertReportFile(t, m, opts, outFile)
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

	opts := ExecOptions{YAML: true, FilePath: outFile}
	runAndAssertReportFile(t, m, opts, outFile)
}

func runAndAssertReportFile(t *testing.T, m *Macro, opts ExecOptions, outFile string) {
	if err := Execute(context.Background(), m, opts); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	content, err := os.ReadFile(outFile)
	if err != nil || len(content) == 0 {
		t.Fatalf("Report file not written: %v", err)
	}
}

func TestResolveTargetDir_WorkingDirPrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "step_workdir")
	_ = os.MkdirAll(subDir, 0755)

	if got := resolveTargetDir(tmpDir, subDir); got != subDir {
		t.Errorf("resolveTargetDir() = %q, want %q", got, subDir)
	}

	nonExistent := filepath.Join(tmpDir, "does-not-exist")
	if fallback := resolveTargetDir(tmpDir, nonExistent); fallback != tmpDir {
		t.Errorf("resolveTargetDir() fallback = %q, want %q", fallback, tmpDir)
	}
}

func TestExecute_StepWorkingDirPrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "custom_target_dir")
	_ = os.MkdirAll(subDir, 0755)

	m := &Macro{
		Name: "test-workingdir-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "echo step-ok", WorkingDir: subDir},
		},
	}

	if err := Execute(context.Background(), m, ExecOptions{DryRun: false}); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}
