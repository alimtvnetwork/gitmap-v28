package macro

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RecordInteractive starts an interactive shell session and records commands to a macro.
func RecordInteractive(name string) error {
	printRecorderHeader(name)
	cwd, _ := os.Getwd()
	m := &Macro{
		Name:      name,
		CreatedAt: time.Now(),
		Steps:     make([]MacroStep, 0),
	}
	var redoStack []MacroStep
	reader := bufio.NewReader(os.Stdin)
	for {
		stepNum := len(m.Steps) + 1
		fmt.Printf("  %s[rec:%s %d]> %s", constants.ColorCyan, name, stepNum, constants.ColorReset)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		cmdText := strings.TrimSpace(line)
		if len(cmdText) == 0 {
			continue
		}
		isHandled, shouldExit, shouldSave := handleSessionCommand(cmdText, m, &redoStack, reader)
		if isHandled && shouldExit && !shouldSave {
			fmt.Printf("  %s▲ Recording cancelled. Macro %q was not saved.%s\n\n", constants.ColorYellow, name, constants.ColorReset)
			return nil
		}
		if isHandled && shouldExit {
			break
		}
		if isHandled {
			continue
		}
		cwd = recordSingleStep(cmdText, cwd, m, &redoStack)
	}
	return finalizeSavedMacro(m)
}

func handleSessionCommand(cmdText string, m *Macro, redoStack *[]MacroStep, r *bufio.Reader) (bool, bool, bool) {
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
		printRecordedSteps(m)
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

func recordSingleStep(cmdText, cwd string, m *Macro, redoStack *[]MacroStep) string {
	expandedCmd := ExpandPathAndEnv(cmdText)
	cwd = updateWorkingDirIfCd(expandedCmd, cwd)
	*redoStack = (*redoStack)[:0]
	execLive(expandedCmd, cwd)
	m.Steps = append(m.Steps, MacroStep{
		StepNum:        len(m.Steps) + 1,
		CommandLine:    expandedCmd,
		WorkingDir:     cwd,
		TimeoutSeconds: 300,
	})
	return cwd
}

func finalizeSavedMacro(m *Macro) error {
	if err := SaveMacro(m); err != nil {
		return fmt.Errorf("could not save macro: %w", err)
	}
	fmt.Printf("\n  %s✔ Saved macro %q with %d steps.%s\n\n",
		constants.ColorGreen, m.Name, len(m.Steps), constants.ColorReset)
	return nil
}

func updateWorkingDirIfCd(cmdText, currentDir string) string {
	target := parseCdTarget(cmdText)
	if target == "" {
		return currentDir
	}
	resolved := target
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(currentDir, target)
	}
	resolved = filepath.Clean(resolved)
	info, err := os.Stat(resolved)
	if err == nil && info.IsDir() {
		return resolved
	}
	return currentDir
}

func parseCdTarget(cmdText string) string {
	trimmed := strings.TrimSpace(cmdText)
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "cd ") || strings.HasPrefix(lower, "cd\t") {
		return strings.TrimSpace(trimmed[3:])
	}
	if strings.HasPrefix(lower, "chdir ") || strings.HasPrefix(lower, "chdir\t") {
		return strings.TrimSpace(trimmed[6:])
	}
	return ""
}

func execLive(cmdText, dir string) {
	var cmd *exec.Cmd
	if runtime.GOOS == constants.OSWindows {
		cmd = exec.CommandContext(context.Background(), "powershell", "-NoProfile", "-Command", cmdText)
	} else {
		cmd = exec.CommandContext(context.Background(), "sh", "-c", cmdText)
	}
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func printRecorderHeader(name string) {
	fmt.Println()
	fmt.Printf("  %s● [REC] Recording session: %q%s\n", constants.ColorRed, name, constants.ColorReset)
	fmt.Println("  Type shell commands. They will execute live and be saved to the macro.")
	printRecorderHelp()
	fmt.Println()
}

func printRecorderHelp() {
	fmt.Println("  Session Commands:")
	fmt.Println("    stop / exit / quit   Save and finish recording")
	fmt.Println("    cancel / abort       Abort recording without saving")
	fmt.Println("    undo                 Remove the last recorded command")
	fmt.Println("    undo-steps <N> [-y]  Remove last N recorded commands")
	fmt.Println("    redo                 Restore previously undone command")
	fmt.Println("    redo-steps <N>       Restore last N undone commands")
	fmt.Println("    list / steps         Show currently recorded steps")
	fmt.Println("    help / ?             Show this help message")
}

func printRecordedSteps(m *Macro) {
	if len(m.Steps) == 0 {
		fmt.Printf("  %s▲ No steps recorded yet.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("\n  %sRecorded steps in %q (%d steps):%s\n", constants.ColorCyan, m.Name, len(m.Steps), constants.ColorReset)
	for _, s := range m.Steps {
		fmt.Printf("    %2d. %s\n", s.StepNum, s.CommandLine)
	}
	fmt.Println()
}
