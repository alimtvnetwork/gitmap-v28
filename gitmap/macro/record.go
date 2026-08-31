package macro

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RecordInteractive starts an interactive shell session and records commands to a macro.
func RecordInteractive(name string) error {
	initialDir, _ := os.Getwd()
	dt := NewDirTracker(initialDir)
	printRecorderHeader(name, initialDir)
	m := &Macro{
		Name:      name,
		CreatedAt: time.Now(),
		Steps:     make([]MacroStep, 0),
	}
	var redoStack []MacroStep
	reader := bufio.NewReader(os.Stdin)
	return runRecorderLoop(name, m, dt, &redoStack, reader)
}

func runRecorderLoop(name string, m *Macro, dt *DirTracker, redoStack *[]MacroStep, r *bufio.Reader) error {
	for {
		stepNum := len(m.Steps) + 1
		printPrompt(name, dt.CurrentDir, stepNum)
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		cmdText := strings.TrimSpace(line)
		if len(cmdText) == 0 {
			continue
		}
		isHandled, shouldExit, shouldSave := handleSessionCommand(cmdText, m, redoStack, r, dt.CurrentDir)
		if isHandled && shouldExit && !shouldSave {
			fmt.Printf("  %s▲ Recording canceled. Macro %q was not saved.%s\n\n", constants.ColorYellow, name, constants.ColorReset)
			return nil
		}
		if isHandled && shouldExit {
			break
		}
		if isHandled {
			continue
		}
		recordSingleStep(cmdText, dt, m, redoStack)
	}
	return finalizeSavedMacro(m)
}

func printPrompt(name, currentDir string, stepNum int) {
	fmt.Printf("  %s┌─%s 📁 %s%s%s\n", constants.ColorDim, constants.ColorReset, constants.ColorCyan, currentDir, constants.ColorReset)
	fmt.Printf("  %s└─%s %s●%s %s[rec:%s %d]%s ➜ ",
		constants.ColorDim, constants.ColorReset,
		constants.ColorRed, constants.ColorReset,
		constants.ColorCyan, name, stepNum, constants.ColorReset)
}

func handleSessionCommand(cmdText string, m *Macro, redoStack *[]MacroStep, r *bufio.Reader, currentDir string) (bool, bool, bool) {
	lower := strings.ToLower(cmdText)
	if lower == "stop" || lower == "exit" || lower == "quit" {
		return true, true, true
	}
	if lower == "cancel" || lower == "abort" {
		return true, true, false
	}
	if lower == "help" || lower == "?" {
		printRecorderHelp()
		return true, false, false
	}
	if lower == "list" || lower == "steps" || lower == "show" {
		printRecordedSteps(m, currentDir)
		return true, false, false
	}
	if strings.HasPrefix(lower, "undo") {
		count, isAutoConfirm := parseUndoParams(cmdText)
		handleUndo(m, redoStack, count, isAutoConfirm, r)
		return true, false, false
	}
	if strings.HasPrefix(lower, "redo") {
		count := parseRedoParams(cmdText)
		handleRedo(m, redoStack, count)
		return true, false, false
	}
	return false, false, false
}

func recordSingleStep(cmdText string, dt *DirTracker, m *Macro, redoStack *[]MacroStep) {
	expandedCmd := ExpandPathAndEnv(cmdText)
	dt.ProcessCd(expandedCmd)
	*redoStack = (*redoStack)[:0]
	elapsed, isSuccess := execLive(expandedCmd, dt.CurrentDir)
	m.Steps = append(m.Steps, MacroStep{
		StepNum:        len(m.Steps) + 1,
		CommandLine:    expandedCmd,
		WorkingDir:     dt.CurrentDir,
		TimeoutSeconds: 300,
	})
	printStepExecutionResult(len(m.Steps), elapsed, isSuccess)
}

func printStepExecutionResult(stepNum int, elapsed time.Duration, isSuccess bool) {
	if isSuccess {
		fmt.Printf("  %s✔ Recorded step %d%s %s(%.1fs)%s\n\n",
			constants.ColorGreen, stepNum, constants.ColorReset,
			constants.ColorDim, elapsed.Seconds(), constants.ColorReset)
		return
	}
	fmt.Printf("  %s▲ Recorded step %d (non-zero exit)%s %s(%.1fs)%s\n\n",
		constants.ColorYellow, stepNum, constants.ColorReset,
		constants.ColorDim, elapsed.Seconds(), constants.ColorReset)
}

func finalizeSavedMacro(m *Macro) error {
	if err := SaveMacro(m); err != nil {
		return fmt.Errorf("could not save macro: %w", err)
	}
	fmt.Printf("\n  %s✔ Saved macro %q with %d steps.%s\n\n",
		constants.ColorGreen, m.Name, len(m.Steps), constants.ColorReset)
	return nil
}

func execLive(cmdText, dir string) (time.Duration, bool) {
	start := time.Now()
	var cmd *exec.Cmd
	if runtime.GOOS == constants.OSWindows {
		cmd = exec.CommandContext(context.Background(), "powershell", "-NoProfile", "-Command", cmdText)
	} else {
		cmd = exec.CommandContext(context.Background(), "sh", "-c", cmdText)
	}
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	elapsed := time.Since(start)
	isSuccess := err == nil
	return elapsed, isSuccess
}

func printRecorderHeader(name, initialDir string) {
	fmt.Println()
	fmt.Printf("  %s● [REC] Recording session: %s%q%s\n", constants.ColorRed, constants.ColorCyan, name, constants.ColorReset)
	fmt.Printf("  📁 %sInitial Directory:%s %s%s%s\n", constants.ColorWhite, constants.ColorReset, constants.ColorCyan, initialDir, constants.ColorReset)
	fmt.Println("  Type shell commands. They will execute live and be saved to the macro.")
	printRecorderHelp()
	fmt.Println()
}

func printRecorderHelp() {
	fmt.Println()
	fmt.Printf("  %sSession Commands:%s\n", constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("    %sstop%s / %sexit%s / %squit%s   Save and finish recording\n", constants.ColorGreen, constants.ColorReset, constants.ColorGreen, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %scancel%s / %sabort%s       Abort recording without saving\n", constants.ColorRed, constants.ColorReset, constants.ColorRed, constants.ColorReset)
	fmt.Printf("    %sundo%s                 Remove the last recorded command\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    %sundo-steps <N> [-y]%s  Remove last N recorded commands\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    %sredo%s                 Restore previously undone command\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sredo-steps <N>%s       Restore last N undone commands\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %slist%s / %ssteps%s         Show currently recorded steps\n", constants.ColorWhite, constants.ColorReset, constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    %shelp%s / %s?%s             Show this help message\n\n", constants.ColorDim, constants.ColorReset, constants.ColorDim, constants.ColorReset)
}

func printRecordedSteps(m *Macro, currentDir string) {
	if len(m.Steps) == 0 {
		fmt.Printf("  %s▲ No steps recorded yet.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("\n  %sRecorded steps in %q (%d steps) · 📁 %s:%s\n",
		constants.ColorCyan, m.Name, len(m.Steps), currentDir, constants.ColorReset)
	for _, s := range m.Steps {
		dirLabel := ""
		if len(s.WorkingDir) > 0 {
			dirLabel = fmt.Sprintf(" %s(dir: %s)%s", constants.ColorDim, s.WorkingDir, constants.ColorReset)
		}
		fmt.Printf("    %2d. %s➜%s %s%s\n", s.StepNum, constants.ColorGreen, constants.ColorReset, s.CommandLine, dirLabel)
	}
	fmt.Println()
}
