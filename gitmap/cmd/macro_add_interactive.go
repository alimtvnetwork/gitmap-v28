// Package cmd — macro_add_interactive.go: interactive prompt/recorder fallback for macro add.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/uipref"
)

type interactiveLoopAction int

const (
	loopActionContinue interactiveLoopAction = iota
	loopActionBreak
)

func resolveStepsInteractively(name string, isExec bool) ([]macro.MacroStep, error) {
	if isTerminalInput() {
		return promptInteractiveMacroSteps(name, isExec)
	}

	steps, err := readPipedMacroSteps()
	if err != nil {
		return nil, apperror.WrapSimple(err, "read piped macro steps")
	}

	if len(steps) == 0 {
		printMacroAddUsage()

		return nil, apperror.NewValidationError("macro name and at least one command required")
	}

	return steps, nil
}

func readPipedMacroSteps() ([]macro.MacroStep, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var steps []macro.MacroStep
	stepNum := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if isDoneStep(line) || isAbortStep(line) {
			break
		}

		if strings.Contains(line, "&&") {
			stepNum = appendChainedSteps(&steps, line, stepNum)
			continue
		}

		steps = append(steps, makeMacroStep(stepNum, line))
		stepNum++
	}

	return steps, scanner.Err()
}

func isTerminalInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

func promptInteractiveMacroSteps(name string, isExec bool) ([]macro.MacroStep, error) {
	printInteractiveMacroHeader(name)

	steps, err := collectInteractiveMacroSteps(name, isExec)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read interactive macro input")
	}

	if len(steps) == 0 {
		printNoCommandsEntered(name)

		return nil, nil
	}

	return steps, nil
}

func collectInteractiveMacroSteps(name string, isExec bool) ([]macro.MacroStep, error) {
	var steps []macro.MacroStep
	stepNum := 1
	state := newInteractiveState()
	state.isExecEnabled = isExec
	reader := newInteractiveLineReader()

	for {
		printMacroPromptPwd()
		prompt := fmt.Sprintf("  Step %d> ", stepNum)
		line, isEof, err := reader.readLine(prompt)
		if err != nil || isEof {
			break
		}

		if processInteractiveStepLine(line, name, state, &steps, &stepNum) == loopActionBreak {
			break
		}
	}

	return steps, nil
}

func printMacroPromptPwd() {
	if !uipref.IsMacroPwdVisible() {
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	fmt.Printf("  %s[PWD: %s]%s\n", constants.ColorCyan, cwd, constants.ColorReset)
}

func processInteractiveStepLine(rawLine, name string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) interactiveLoopAction {
	line := strings.TrimSpace(rawLine)
	if processInBuilderCommand(line, state, steps, stepNum) {
		return loopActionContinue
	}

	if isDoneStep(line) {
		return loopActionBreak
	}

	if isAbortStep(line) {
		printMacroAbortedMessage(name)
		*steps = nil

		return loopActionBreak
	}

	return recordStepLine(line, name, state, steps, stepNum)
}

func processInBuilderCommand(line string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) bool {
	if !isInteractiveHelper(line) {
		return false
	}

	return handleInteractiveHelper(line, state, steps, stepNum)
}

func recordStepLine(line, name string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) interactiveLoopAction {
	if isRecordStep(line) {
		_ = macro.RecordInteractive(name)

		return loopActionBreak
	}

	runLiveStepIfEnabled(line, state.isExecEnabled)

	return appendRecordedStep(line, steps, stepNum)
}

func runLiveStepIfEnabled(line string, isExecEnabled bool) {
	if !isExecEnabled {
		return
	}

	if err := executeLiveCommand(line); err != nil {
		ensureTerminalVisibility()
	}
}

func ensureTerminalVisibility() {
	_ = os.Stdout.Sync()
	_ = os.Stderr.Sync()
}

func appendRecordedStep(line string, steps *[]macro.MacroStep, stepNum *int) interactiveLoopAction {
	if strings.Contains(line, "&&") {
		*stepNum = appendChainedSteps(steps, line, *stepNum)

		return loopActionContinue
	}

	*steps = append(*steps, makeMacroStep(*stepNum, line))
	fmt.Printf("  %s✓ Recorded Step %d: %s%s (will run when macro is executed)\n\n",
		constants.ColorGreen, *stepNum, line, constants.ColorReset)
	*stepNum++

	return loopActionContinue
}

func executeLiveCommand(cmdText string) error {
	cmd := buildLiveExecCmd(cmdText)
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s▲ Live command returned error: %v%s\n\n",
			constants.ColorYellow, err, constants.ColorReset)
	}

	return err
}

func buildLiveExecCmd(cmdText string) *exec.Cmd {
	exeCmd := resolveLiveCommandText(cmdText)
	var cmd *exec.Cmd
	if runtime.GOOS == constants.OSWindows {
		cmd = exec.Command("powershell", "-NoProfile", "-Command", exeCmd)
	} else {
		cmd = exec.Command("sh", "-c", exeCmd)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func resolveLiveCommandText(cmdText string) string {
	trimmed := strings.TrimSpace(cmdText)
	if !strings.Contains(trimmed, "&&") {
		return resolveSingleLiveCmd(trimmed)
	}

	parts := splitAndResolveParts(trimmed)
	if runtime.GOOS == constants.OSWindows {
		return chainWindowsCompoundCommands(parts)
	}

	return strings.Join(parts, " && ")
}

func splitAndResolveParts(cmdText string) []string {
	rawParts := strings.Split(cmdText, "&&")
	parts := make([]string, 0, len(rawParts))
	for _, raw := range rawParts {
		if sub := strings.TrimSpace(raw); sub != "" {
			parts = append(parts, resolveSingleLiveCmd(sub))
		}
	}

	return parts
}

func chainWindowsCompoundCommands(parts []string) string {
	if len(parts) == 0 {
		return ""
	}

	res := parts[len(parts)-1]
	for i := len(parts) - 2; i >= 0; i-- {
		res = fmt.Sprintf("%s; if ($?) { %s }", parts[i], res)
	}

	return res
}

func resolveSingleLiveCmd(cmdText string) string {
	trimmed := strings.TrimSpace(cmdText)
	lower := strings.ToLower(trimmed)
	if !strings.HasPrefix(lower, "gitmap ") && lower != "gitmap" {
		return trimmed
	}

	exe, err := os.Executable()
	if err != nil {
		return trimmed
	}

	args := strings.TrimSpace(trimmed[len("gitmap"):])

	return formatLiveExeCmd(exe, args)
}

func formatLiveExeCmd(exe, args string) string {
	if runtime.GOOS == constants.OSWindows {
		return formatWindowsLiveExe(exe, args)
	}

	return formatUnixLiveExe(exe, args)
}

func formatWindowsLiveExe(exe, args string) string {
	if args == "" {
		return fmt.Sprintf("& %q", exe)
	}

	return fmt.Sprintf("& %q %s", exe, args)
}

func formatUnixLiveExe(exe, args string) string {
	if args == "" {
		return fmt.Sprintf("%q", exe)
	}

	return fmt.Sprintf("%q %s", exe, args)
}

func printNoCommandsEntered(name string) {
	fmt.Printf("  %s▲ No commands entered. Macro %q was not saved.%s\n\n", constants.ColorYellow, name, constants.ColorReset)
}

func printMacroAbortedMessage(name string) {
	fmt.Printf("  %s▲ Macro %q creation aborted.%s\n\n", constants.ColorYellow, name, constants.ColorReset)
}

func isDoneStep(line string) bool {
	return line == "" || strings.EqualFold(line, "done") || strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit")
}

func isAbortStep(line string) bool {
	return strings.EqualFold(line, "cancel") || strings.EqualFold(line, "abort")
}

func isRecordStep(line string) bool {
	return strings.EqualFold(line, "rec") || strings.EqualFold(line, "record")
}

func printInteractiveMacroHeader(name string) {
	fmt.Println()
	fmt.Printf("  %s● Interactive Macro Builder: %s%q%s\n", constants.ColorCyan, constants.ColorWhite, name, constants.ColorReset)
	fmt.Println("  Enter commands one per line (empty line or 'done' to save, 'cancel' to abort):")
	fmt.Printf("  %s(Commands: 'ls', 'mkdir', 'cd', 'pwd on/off', 'find', 'search', 'replace', 'exec on/off', 'help', 'rec')%s\n\n", constants.ColorDim, constants.ColorReset)
}

func printMacroAddUsage() {
	fmt.Println("Usage: gitmap macro add <name> <command1> [command2...] [--desc <text>] [--tag <tag>] [--pwd|--no-pwd]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap macro add build \"go build -o app.exe .\" \"go test ./...\"")
	fmt.Println("  gitmap macro add deploy \"git push origin main\" --desc \"Deploy to main\"")
	fmt.Println("  gitmap macro add alim               # enter commands interactively")
	fmt.Println()
}
