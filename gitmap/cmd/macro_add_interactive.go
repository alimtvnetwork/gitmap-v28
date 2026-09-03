// Package cmd — macro_add_interactive.go: interactive prompt/recorder fallback for macro add.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func resolveStepsInteractively(name string) ([]macro.MacroStep, error) {
	if isTerminalInput() {
		return promptInteractiveMacroSteps(name)
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

func promptInteractiveMacroSteps(name string) ([]macro.MacroStep, error) {
	printInteractiveMacroHeader(name)

	scanner := bufio.NewScanner(os.Stdin)
	var steps []macro.MacroStep
	stepNum := 1

	for {
		fmt.Printf("  Step %d> ", stepNum)
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if isDoneStep(line) {
			break
		}

		if isAbortStep(line) {
			fmt.Printf("  %s▲ Macro %q creation aborted.%s\n\n", constants.ColorYellow, name, constants.ColorReset)

			return nil, nil
		}

		if isRecordStep(line) {
			return nil, macro.RecordInteractive(name)
		}

		if strings.Contains(line, "&&") {
			stepNum = appendChainedSteps(&steps, line, stepNum)
			continue
		}

		steps = append(steps, makeMacroStep(stepNum, line))
		stepNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "read interactive macro input")
	}

	if len(steps) == 0 {
		fmt.Printf("  %s▲ No commands entered. Macro %q was not saved.%s\n\n", constants.ColorYellow, name, constants.ColorReset)

		return nil, nil
	}

	return steps, nil
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
	fmt.Printf("  %s(Tip: type 'rec' or 'record' to launch live command recording session)%s\n\n", constants.ColorDim, constants.ColorReset)
}

func printMacroAddUsage() {
	fmt.Println("Usage: gitmap macro add <name> <command1> [command2...] [--desc <text>] [--tag <tag>]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap macro add build \"go build -o app.exe .\" \"go test ./...\"")
	fmt.Println("  gitmap macro add deploy \"git push origin main\" --desc \"Deploy to main\"")
	fmt.Println("  gitmap macro add alim               # enter commands interactively")
	fmt.Println()
}
