// Package cmd — macro_edit.go: interactive editor for existing macros with live execution.
package cmdmacro

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func handleMacroEdit(args []string) error {
	name, isExec := parseMacroEditArgs(args)
	if name == "" {
		printMacroEditUsage()

		return apperror.NewValidationError("macro name required for edit")
	}

	m, err := macro.LoadMacro(name)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("load macro %q", name))
	}

	return runInteractiveMacroEdit(m, isExec)
}

func parseMacroEditArgs(args []string) (string, bool) {
	name := ""
	isExec := true
	for _, a := range args {
		if a == "--no-exec" {
			isExec = false
			continue
		}

		if a == "--exec" || a == "-e" {
			isExec = true
			continue
		}

		if !strings.HasPrefix(a, "-") && name == "" {
			name = a
		}
	}

	return name, isExec
}

func printMacroEditUsage() {
	fmt.Println("Usage: gitmap macro edit <name> [--no-exec]")
	fmt.Println()
	fmt.Println("Interactive editing of existing macros with live command execution.")
	fmt.Println()
}

func runInteractiveMacroEdit(m *macro.Macro, isExec bool) error {
	printMacroEditHeader(m)
	state := newInteractiveState()
	state.isExecEnabled = isExec
	reader := newInteractiveLineReader()

	steps := m.Steps
	stepNum := len(steps) + 1

	return loopInteractiveMacroEdit(m, state, reader, &steps, &stepNum)
}

func printMacroEditHeader(m *macro.Macro) {
	fmt.Println()
	fmt.Printf("  %s● Interactive Macro Editor: %s%q%s (%d steps)\n",
		constants.ColorCyan, constants.ColorWhite, m.Name, constants.ColorReset, len(m.Steps))
	printCurrentEditSteps(m.Steps)
	printMacroEditInstructions()
}

func printCurrentEditSteps(steps []macro.MacroStep) {
	if len(steps) == 0 {
		fmt.Printf("  %s(No steps currently defined)%s\n", constants.ColorDim, constants.ColorReset)

		return
	}

	fmt.Println("  Current steps:")
	for i, s := range steps {
		fmt.Printf("    %d. %s\n", i+1, s.CommandLine)
	}

	fmt.Println()
}

func printMacroEditInstructions() {
	fmt.Println("  Commands: 'del <n>', 'replace <n> <cmd>', 'insert <n> <cmd>', 'list', 'done', 'cancel'")
	fmt.Println("  In-builder: 'cat <file>', 'touch <file>', 'mkfile <file>', 'copy', 'paste', 'explorer', 'browse'")
	fmt.Println("  Type any shell command to append it (executes live in terminal):")
	fmt.Println()
}

func loopInteractiveMacroEdit(m *macro.Macro, state *interactiveSessionState, reader *interactiveLineReader, steps *[]macro.MacroStep, stepNum *int) error {
	for {
		printMacroPromptPwd()
		prompt := fmt.Sprintf("  Edit [%d]> ", *stepNum)
		line, isEof, err := reader.readLine(prompt)
		if err != nil || isEof {
			break
		}

		action := processEditStepLine(line, m, state, steps, stepNum)
		if action == loopActionBreak {
			break
		}
	}

	return nil
}

func processEditStepLine(rawLine string, m *macro.Macro, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) interactiveLoopAction {
	line := strings.TrimSpace(rawLine)
	if isDoneStep(line) {
		return saveAndFinishEdit(m, *steps)
	}

	if isAbortStep(line) {
		fmt.Printf("  %s▲ Macro %q editing aborted (no changes saved).%s\n\n",
			constants.ColorYellow, m.Name, constants.ColorReset)

		return loopActionBreak
	}

	if handleEditSpecialAction(line, state, steps, stepNum) {
		return loopActionContinue
	}

	return recordStepLine(line, m.Name, state, steps, stepNum)
}

func handleEditSpecialAction(line string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) bool {
	if isListStepsCmd(line) {
		printCurrentEditSteps(*steps)

		return true
	}

	if isDeleteStepCmd(line) {
		return handleDeleteStepCmd(line, steps, stepNum)
	}

	if isReplaceStepCmd(line) {
		return handleReplaceStepCmd(line, state, steps)
	}

	if isInsertStepCmd(line) {
		return handleInsertStepCmd(line, state, steps, stepNum)
	}

	return processInBuilderCommand(line, state, steps, stepNum)
}

func isListStepsCmd(line string) bool {
	low := strings.ToLower(line)

	return low == "list" || low == "show" || low == ":list" || low == ":show"
}

func isDeleteStepCmd(line string) bool {
	low := strings.ToLower(line)

	return strings.HasPrefix(low, "del ") || strings.HasPrefix(low, "rm ") || strings.HasPrefix(low, "delete ")
}

func isReplaceStepCmd(line string) bool {
	low := strings.ToLower(line)

	return strings.HasPrefix(low, "replace ") || strings.HasPrefix(low, ":replace ")
}

func isInsertStepCmd(line string) bool {
	low := strings.ToLower(line)

	return strings.HasPrefix(low, "insert ") || strings.HasPrefix(low, ":insert ")
}

func handleDeleteStepCmd(line string, steps *[]macro.MacroStep, stepNum *int) bool {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		fmt.Printf("  %s▲ Usage: del <step-number>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	idx, err := strconv.Atoi(parts[1])
	if err != nil || idx < 1 || idx > len(*steps) {
		fmt.Printf("  %s▲ Invalid step number %q (1..%d)%s\n\n", constants.ColorRed, parts[1], len(*steps), constants.ColorReset)

		return true
	}

	removed := (*steps)[idx-1].CommandLine
	*steps = append((*steps)[:idx-1], (*steps)[idx:]...)
	reindexMacroSteps(steps)
	*stepNum = len(*steps) + 1
	fmt.Printf("  %s✓ Removed step %d: %s%s\n\n", constants.ColorGreen, idx, removed, constants.ColorReset)

	return true
}

func handleReplaceStepCmd(line string, state *interactiveSessionState, steps *[]macro.MacroStep) bool {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		fmt.Printf("  %s▲ Usage: replace <step-number> <new-command>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	idx, err := strconv.Atoi(parts[1])
	if err != nil || idx < 1 || idx > len(*steps) {
		fmt.Printf("  %s▲ Invalid step number %q (1..%d)%s\n\n", constants.ColorRed, parts[1], len(*steps), constants.ColorReset)

		return true
	}

	newCmd := strings.TrimSpace(parts[2])
	runLiveStepIfEnabled(newCmd, state.isExecEnabled)
	(*steps)[idx-1].CommandLine = newCmd
	fmt.Printf("  %s✓ Replaced step %d with: %s%s\n\n", constants.ColorGreen, idx, newCmd, constants.ColorReset)

	return true
}

func handleInsertStepCmd(line string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) bool {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		fmt.Printf("  %s▲ Usage: insert <step-number> <command>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	idx, err := strconv.Atoi(parts[1])
	if err != nil || idx < 1 || idx > len(*steps)+1 {
		fmt.Printf("  %s▲ Invalid insert index %q (1..%d)%s\n\n", constants.ColorRed, parts[1], len(*steps)+1, constants.ColorReset)

		return true
	}

	newCmd := strings.TrimSpace(parts[2])
	runLiveStepIfEnabled(newCmd, state.isExecEnabled)
	insertStepAt(steps, idx-1, makeMacroStep(idx, newCmd))
	reindexMacroSteps(steps)
	*stepNum = len(*steps) + 1
	fmt.Printf("  %s✓ Inserted step %d: %s%s\n\n", constants.ColorGreen, idx, newCmd, constants.ColorReset)

	return true
}

func insertStepAt(steps *[]macro.MacroStep, pos int, s macro.MacroStep) {
	if pos >= len(*steps) {
		*steps = append(*steps, s)

		return
	}

	*steps = append((*steps)[:pos], append([]macro.MacroStep{s}, (*steps)[pos:]...)...)
}

func reindexMacroSteps(steps *[]macro.MacroStep) {
	for i := range *steps {
		(*steps)[i].StepNum = i + 1
	}
}

func saveAndFinishEdit(m *macro.Macro, steps []macro.MacroStep) interactiveLoopAction {
	m.Steps = steps
	if err := macro.SaveMacro(m); err != nil {
		fmt.Printf("  %s▲ Failed saving macro: %v%s\n\n", constants.ColorRed, err, constants.ColorReset)

		return loopActionBreak
	}

	fmt.Printf("  %s✔ Macro %q successfully updated (%d steps)%s\n\n",
		constants.ColorGreen, m.Name, len(steps), constants.ColorReset)

	return loopActionBreak
}
