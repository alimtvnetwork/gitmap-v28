// Package cmd — macro_add_helpers_test.go: unit tests for interactive macro builder helpers.
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/uipref"
)

func TestMacroPromptPwdToggle(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)
	uipref.ResetMacroPwdOverride()
	defer uipref.ResetMacroPwdOverride()

	if !handlePwdCmd("pwd off") {
		t.Fatal("expected handlePwdCmd to handle 'pwd off'")
	}

	if uipref.IsMacroPwdVisible() {
		t.Fatal("expected PWD display to be false after 'pwd off'")
	}

	if !handlePwdCmd("pwd on") {
		t.Fatal("expected handlePwdCmd to handle 'pwd on'")
	}

	if !uipref.IsMacroPwdVisible() {
		t.Fatal("expected PWD display to be true after 'pwd on'")
	}
}

func TestMacroInteractiveLs(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	_ = os.MkdirAll(subDir, 0755)

	sampleFile := filepath.Join(tempDir, "sample.txt")
	_ = os.WriteFile(sampleFile, []byte("hello gitmap macro builder"), 0644)

	state := newInteractiveState()
	if !handleLsOrDir("ls "+tempDir, state) {
		t.Fatal("expected handleLsOrDir to return true")
	}

	if state.lastInspectedCmd != "ls "+tempDir {
		t.Fatalf("expected lastInspectedCmd to be saved, got %q", state.lastInspectedCmd)
	}

	if err := executeInteractiveLs(tempDir); err != nil {
		t.Fatalf("executeInteractiveLs failed: %v", err)
	}
}

func TestMacroInteractiveFindAndSearch(t *testing.T) {
	tempDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(tempDir)

	testFile := filepath.Join(tempDir, "query_target.txt")
	_ = os.WriteFile(testFile, []byte("alpha\nbeta target_string gamma\ndelta\n"), 0644)

	state := newInteractiveState()
	if !handleFindCmd("find query*", state) {
		t.Fatal("expected handleFindCmd to return true")
	}

	if !handleSearchCmd("search target_string", state) {
		t.Fatal("expected handleSearchCmd to return true")
	}
}

func TestMacroInteractiveReplace(t *testing.T) {
	tempDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(tempDir)

	filePath := filepath.Join(tempDir, "replace_me.txt")
	_ = os.WriteFile(filePath, []byte("foo bar foo baz"), 0644)

	if err := executeInteractiveReplace("foo replacement replace_me.txt"); err != nil {
		t.Fatalf("executeInteractiveReplace failed: %v", err)
	}

	updated, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read updated file failed: %v", err)
	}

	if string(updated) != "replacement bar replacement baz" {
		t.Fatalf("unexpected content after replace: %q", string(updated))
	}
}

func TestProcessInteractiveStepLine(t *testing.T) {
	state := newInteractiveState()
	var steps []macro.MacroStep
	stepNum := 1

	// In-builder ls does not add step, but records last inspected command
	action := processInteractiveStepLine("ls", "test-macro", state, &steps, &stepNum)
	if action != loopActionContinue || len(steps) != 0 {
		t.Fatalf("expected ls not to record step directly, steps=%d", len(steps))
	}

	// +add records last inspected command
	action = processInteractiveStepLine("+add", "test-macro", state, &steps, &stepNum)
	if action != loopActionContinue || len(steps) != 1 || steps[0].CommandLine != "ls" {
		t.Fatalf("expected +add to record ls, steps=%+v", steps)
	}

	// Regular command
	action = processInteractiveStepLine("git status", "test-macro", state, &steps, &stepNum)
	if action != loopActionContinue || len(steps) != 2 || steps[1].CommandLine != "git status" {
		t.Fatalf("expected git status to be recorded, steps=%+v", steps)
	}

	// Explicit add command
	action = processInteractiveStepLine("add ls -la", "test-macro", state, &steps, &stepNum)
	if action != loopActionContinue || len(steps) != 3 || steps[2].CommandLine != "ls -la" {
		t.Fatalf("expected add ls -la to be recorded, steps=%+v", steps)
	}

	// Done command
	action = processInteractiveStepLine("done", "test-macro", state, &steps, &stepNum)
	if action != loopActionBreak {
		t.Fatalf("expected done to break loop")
	}
}

func TestMacroInteractiveExecToggle(t *testing.T) {
	state := newInteractiveState()
	if !handleExecCmd("exec on", state) || !state.isExecEnabled {
		t.Fatal("expected exec on to enable live execution")
	}

	if !handleExecCmd("exec off", state) || state.isExecEnabled {
		t.Fatal("expected exec off to disable live execution")
	}

	if !handleExecCmd("exec", state) || !state.isExecEnabled {
		t.Fatal("expected exec toggle to enable live execution")
	}
}

func TestMacroInteractiveCdExpansion(t *testing.T) {
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	state := newInteractiveState()
	if !handleCdCmd("cd %temp%", state) || state.lastInspectedCmd != "cd %temp%" {
		t.Fatalf("expected handleCdCmd with %%temp%% to succeed")
	}

	if !handleCdCmd("cd //temp", state) || !handleCdCmd("cd /temp", state) {
		t.Fatal("expected handleCdCmd with //temp and /temp to succeed")
	}

	if !handleCdCmd("cd ~", state) {
		t.Fatal("expected handleCdCmd with ~ to succeed")
	}
}
