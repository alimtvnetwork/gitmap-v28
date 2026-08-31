package macro

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func parseUndoParams(cmdText string) (int, bool) {
	parts := strings.Fields(cmdText)
	count := 1
	isAutoConfirm := false
	for _, part := range parts[1:] {
		if part == "-y" || part == "--yes" {
			isAutoConfirm = true
			continue
		}
		if n, err := strconv.Atoi(part); err == nil && n > 0 {
			count = n
		}
	}
	return count, isAutoConfirm
}

func parseRedoParams(cmdText string) int {
	parts := strings.Fields(cmdText)
	count := 1
	for _, part := range parts[1:] {
		if n, err := strconv.Atoi(part); err == nil && n > 0 {
			count = n
		}
	}
	return count
}

func handleUndo(m *Macro, redoStack *[]MacroStep, count int, isAutoConfirm bool, reader *bufio.Reader) {
	if len(m.Steps) == 0 {
		fmt.Printf("  %s▲ No steps to undo.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	count = clampCount(count, len(m.Steps))
	if count > 1 && !isAutoConfirm && !confirmUndoPrompt(count, reader) {
		fmt.Println("  ▲ Undo canceled.")
		return
	}
	applyUndo(m, redoStack, count)
}

func clampCount(requested, maxLimit int) int {
	if requested < 1 {
		return 1
	}
	if requested > maxLimit {
		return maxLimit
	}
	return requested
}

func confirmUndoPrompt(count int, reader *bufio.Reader) bool {
	fmt.Printf("  Undo last %d steps? [y/N]: ", count)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "y" || ans == "yes"
}

func applyUndo(m *Macro, redoStack *[]MacroStep, count int) {
	splitIdx := len(m.Steps) - count
	undone := m.Steps[splitIdx:]
	m.Steps = m.Steps[:splitIdx]
	for i := len(undone) - 1; i >= 0; i-- {
		*redoStack = append(*redoStack, undone[i])
	}
	printUndoResult(undone, count, len(m.Steps))
}

func printUndoResult(undone []MacroStep, count, remaining int) {
	if count == 1 {
		last := undone[0]
		fmt.Printf("  %s✔ Undone step %d: %q%s\n", constants.ColorGreen, last.StepNum, last.CommandLine, constants.ColorReset)
		fmt.Printf("  %s💡 Tip: Type 'redo' to restore undone step, or 'undo-steps <N>' to undo multiple steps.%s\n", constants.ColorCyan, constants.ColorReset)
		return
	}
	fmt.Printf("  %s✔ Undone %d steps. (Current steps: %d)%s\n", constants.ColorGreen, count, remaining, constants.ColorReset)
	fmt.Printf("  %s💡 Tip: Type 'redo' or 'redo-steps %d' to restore undone steps.%s\n", constants.ColorCyan, count, constants.ColorReset)
}

func handleRedo(m *Macro, redoStack *[]MacroStep, count int) {
	if len(*redoStack) == 0 {
		fmt.Printf("  %s▲ No steps to redo.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	count = clampCount(count, len(*redoStack))
	for i := 0; i < count; i++ {
		lastIdx := len(*redoStack) - 1
		step := (*redoStack)[lastIdx]
		*redoStack = (*redoStack)[:lastIdx]
		step.StepNum = len(m.Steps) + 1
		m.Steps = append(m.Steps, step)
	}
	printRedoResult(m, count)
}

func printRedoResult(m *Macro, count int) {
	if count == 1 {
		last := m.Steps[len(m.Steps)-1]
		fmt.Printf("  %s✔ Restored step %d: %q%s\n", constants.ColorGreen, last.StepNum, last.CommandLine, constants.ColorReset)
		return
	}
	fmt.Printf("  %s✔ Restored %d steps. (Current steps: %d)%s\n", constants.ColorGreen, count, len(m.Steps), constants.ColorReset)
}
