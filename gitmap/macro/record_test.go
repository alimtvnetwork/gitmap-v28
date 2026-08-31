package macro

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestParseUndoParams(t *testing.T) {
	c1, y1 := parseUndoParams("undo")
	if c1 != 1 || y1 {
		t.Errorf("parseUndoParams('undo') = (%d, %v), want (1, false)", c1, y1)
	}

	c2, y2 := parseUndoParams("undo-steps 3 -y")
	if c2 != 3 || !y2 {
		t.Errorf("parseUndoParams('undo-steps 3 -y') = (%d, %v), want (3, true)", c2, y2)
	}

	c3, y3 := parseUndoParams("undo 5 --yes")
	if c3 != 5 || !y3 {
		t.Errorf("parseUndoParams('undo 5 --yes') = (%d, %v), want (5, true)", c3, y3)
	}
}

func TestParseRedoParams(t *testing.T) {
	c1 := parseRedoParams("redo")
	if c1 != 1 {
		t.Errorf("parseRedoParams('redo') = %d, want 1", c1)
	}

	c2 := parseRedoParams("redo-steps 4")
	if c2 != 4 {
		t.Errorf("parseRedoParams('redo-steps 4') = %d, want 4", c2)
	}
}

func TestUndoAndRedoLogic(t *testing.T) {
	m := &Macro{
		Name: "test-undo-redo",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "echo step 1"},
			{StepNum: 2, CommandLine: "echo step 2"},
			{StepNum: 3, CommandLine: "echo step 3"},
		},
	}
	var redoStack []MacroStep
	reader := bufio.NewReader(strings.NewReader("y\n"))

	// Undo 1 step
	handleUndo(m, &redoStack, 1, false, reader)
	if len(m.Steps) != 2 {
		t.Fatalf("expected 2 steps remaining, got %d", len(m.Steps))
	}
	if len(redoStack) != 1 {
		t.Fatalf("expected 1 step on redoStack, got %d", len(redoStack))
	}

	// Redo 1 step
	handleRedo(m, &redoStack, 1)
	if len(m.Steps) != 3 {
		t.Fatalf("expected 3 steps after redo, got %d", len(m.Steps))
	}
	if len(redoStack) != 0 {
		t.Fatalf("expected empty redoStack, got %d", len(redoStack))
	}

	// Undo 2 steps with auto confirm
	handleUndo(m, &redoStack, 2, true, reader)
	if len(m.Steps) != 1 {
		t.Fatalf("expected 1 step remaining after undo 2, got %d", len(m.Steps))
	}
	if len(redoStack) != 2 {
		t.Fatalf("expected 2 steps on redoStack, got %d", len(redoStack))
	}

	// Redo 2 steps
	handleRedo(m, &redoStack, 2)
	if len(m.Steps) != 3 {
		t.Fatalf("expected 3 steps after redo 2, got %d", len(m.Steps))
	}
}

func TestHandleSessionCommand(t *testing.T) {
	m := &Macro{Name: "session-test"}
	var redoStack []MacroStep
	reader := bufio.NewReader(strings.NewReader(""))
	currentDir, _ := os.Getwd()

	// Test stop
	isHandled, shouldExit, shouldSave := handleSessionCommand("stop", m, &redoStack, reader, currentDir)
	if !isHandled || !shouldExit || !shouldSave {
		t.Errorf("handleSessionCommand('stop') unexpected result: (%v, %v, %v)", isHandled, shouldExit, shouldSave)
	}

	// Test cancel
	isHandled, shouldExit, shouldSave = handleSessionCommand("cancel", m, &redoStack, reader, currentDir)
	if !isHandled || !shouldExit || shouldSave {
		t.Errorf("handleSessionCommand('cancel') unexpected result: (%v, %v, %v)", isHandled, shouldExit, shouldSave)
	}

	// Test help
	isHandled, shouldExit, _ = handleSessionCommand("help", m, &redoStack, reader, currentDir)
	if !isHandled || shouldExit {
		t.Errorf("handleSessionCommand('help') unexpected result: (%v, %v)", isHandled, shouldExit)
	}

	// Test non-command
	isHandled, _, _ = handleSessionCommand("git status", m, &redoStack, reader, currentDir)
	if isHandled {
		t.Errorf("handleSessionCommand('git status') should not be handled as session command")
	}
}

func TestDirTracker_ProcessCd(t *testing.T) {
	tmpDir := t.TempDir()
	dt := NewDirTracker(tmpDir)

	// Valid sub directory
	subDir := t.TempDir()
	isChanged := dt.ProcessCd("cd " + subDir)
	if !isChanged || dt.CurrentDir != subDir {
		t.Fatalf("expected currentDir %s, got %s", subDir, dt.CurrentDir)
	}

	// cd - to go back
	isSwapped := dt.ProcessCd("cd -")
	if !isSwapped || dt.CurrentDir != tmpDir {
		t.Fatalf("expected currentDir %s after cd -, got %s", tmpDir, dt.CurrentDir)
	}
}
