package cmdmacro

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func TestParseMacroEditArgs(t *testing.T) {
	name, isExec := parseMacroEditArgs([]string{"deploy"})
	if name != "deploy" || !isExec {
		t.Fatalf("expected name='deploy', isExec=true; got name=%q, isExec=%v", name, isExec)
	}

	name2, isExec2 := parseMacroEditArgs([]string{"build", "--no-exec"})
	if name2 != "build" || isExec2 {
		t.Fatalf("expected name='build', isExec=false; got name=%q, isExec=%v", name2, isExec2)
	}
}

func TestMacroEditStepCommands(t *testing.T) {
	if !isListStepsCmd("list") || !isListStepsCmd(":show") {
		t.Fatal("expected isListStepsCmd to recognize list and :show")
	}

	if !isDeleteStepCmd("del 1") || !isDeleteStepCmd("rm 2") {
		t.Fatal("expected isDeleteStepCmd to recognize del and rm")
	}

	if !isReplaceStepCmd("replace 1 echo hello") {
		t.Fatal("expected isReplaceStepCmd to recognize replace")
	}

	if !isInsertStepCmd("insert 2 echo inserted") {
		t.Fatal("expected isInsertStepCmd to recognize insert")
	}
}

func TestInsertStepAtAndReindex(t *testing.T) {
	steps := []macro.MacroStep{
		{StepNum: 1, CommandLine: "cmd1"},
		{StepNum: 2, CommandLine: "cmd2"},
	}

	insertStepAt(&steps, 1, macro.MacroStep{CommandLine: "cmd_inserted"})
	if len(steps) != 3 || steps[1].CommandLine != "cmd_inserted" {
		t.Fatalf("expected cmd_inserted at index 1, got %+v", steps)
	}

	reindexMacroSteps(&steps)
	if steps[0].StepNum != 1 || steps[1].StepNum != 2 || steps[2].StepNum != 3 {
		t.Fatalf("expected step numbers 1, 2, 3; got %+v", steps)
	}
}
