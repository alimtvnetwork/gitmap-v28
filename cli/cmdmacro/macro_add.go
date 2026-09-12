// Package cmd — macro_add.go: create/add new macros directly from CLI arguments.
package cmdmacro

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"github.com/alimtvnetwork/gitmap-v28/cli/uipref"
)

func handleMacroAdd(args []string) error {
	checkHelp("macro", args)

	if len(args) == 0 {
		printMacroAddUsage()

		return apperror.NewValidationError("macro name required")
	}

	name := args[0]
	if isMacroHelpArg(name) {
		printMacroAddUsage()

		return nil
	}

	desc, tag, isExec, rawSteps := parseMacroAddFlags(args[1:])
	steps := parseMacroStepsList(rawSteps)

	steps, resolveErr := ensureMacroSteps(name, steps, isExec)
	if resolveErr != nil {
		return resolveErr
	}

	if len(steps) == 0 {
		return nil
	}

	m := buildNewMacro(name, desc, tag, steps)

	if err := macro.SaveMacro(&m); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("save macro %s", name))
	}

	printMacroCreatedSuccess(m)

	return nil
}

func ensureMacroSteps(name string, steps []macro.MacroStep, isExec bool) ([]macro.MacroStep, error) {
	if len(steps) > 0 {
		return steps, nil
	}

	return resolveStepsInteractively(name, isExec)
}

func parseMacroAddFlags(args []string) (string, string, bool, []string) {
	var desc, tag string
	isExec := true
	var rawSteps []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if handleMacroAddSpecialFlag(a, &isExec) {
			continue
		}

		i, desc, tag, rawSteps = parseSingleMacroFlag(args, i, desc, tag, rawSteps)
	}

	return desc, tag, isExec, rawSteps
}

func handleMacroAddSpecialFlag(flag string, isExec *bool) bool {
	if flag == "--pwd" {
		uipref.SetMacroPwdOverride(true)

		return true
	}

	if flag == "--no-pwd" {
		uipref.SetMacroPwdOverride(false)

		return true
	}

	if flag == "--no-exec" {
		*isExec = false

		return true
	}

	if flag == "--exec" || flag == "-e" {
		*isExec = true

		return true
	}

	return false
}

func parseSingleMacroFlag(args []string, i int, desc, tag string, rawSteps []string) (int, string, string, []string) {
	a := args[i]
	switch {
	case (a == "--desc" || a == "--description") && i+1 < len(args):
		return i + 1, args[i+1], tag, rawSteps
	case a == "--tag" && i+1 < len(args):
		return i + 1, desc, args[i+1], rawSteps
	case strings.HasPrefix(a, "--desc="):
		return i, strings.TrimPrefix(a, "--desc="), tag, rawSteps
	case strings.HasPrefix(a, "--tag="):
		return i, desc, strings.TrimPrefix(a, "--tag="), rawSteps
	case !strings.HasPrefix(a, "-"):
		return i, desc, tag, append(rawSteps, a)
	default:
		return i, desc, tag, rawSteps
	}
}

func parseMacroStepsList(rawSteps []string) []macro.MacroStep {
	var steps []macro.MacroStep
	stepNum := 1
	for _, raw := range rawSteps {
		if strings.Contains(raw, "&&") {
			stepNum = appendChainedSteps(&steps, raw, stepNum)
			continue
		}

		trimmed := strings.TrimSpace(raw)
		if trimmed != "" {
			steps = append(steps, makeMacroStep(stepNum, trimmed))
			stepNum++
		}
	}

	return steps
}

func appendChainedSteps(steps *[]macro.MacroStep, raw string, stepNum int) int {
	subSteps := strings.Split(raw, "&&")
	for _, sub := range subSteps {
		trimmed := strings.TrimSpace(sub)
		if trimmed != "" {
			*steps = append(*steps, makeMacroStep(stepNum, trimmed))
			stepNum++
		}
	}

	return stepNum
}

func makeMacroStep(num int, cmdLine string) macro.MacroStep {
	return macro.MacroStep{
		StepNum:         num,
		CommandLine:     cmdLine,
		TimeoutSeconds:  300,
		ContinueOnError: false,
	}
}

func buildNewMacro(name, desc, tag string, steps []macro.MacroStep) macro.Macro {
	now := time.Now()

	return macro.Macro{
		Name:        name,
		Description: desc,
		Tags:        tag,
		CreatedAt:   now,
		UpdatedAt:   now,
		TotalSteps:  len(steps),
		Steps:       steps,
	}
}

func printMacroCreatedSuccess(m macro.Macro) {
	fmt.Printf("\n\033[1;92m✓ Created macro\033[0m \033[1m%q\033[0m (%d step(s))\n", m.Name, len(m.Steps))
	for _, step := range m.Steps {
		fmt.Printf("  %2d. %s\n", step.StepNum, step.CommandLine)
	}

	fmt.Printf("\nRun with: \033[1;96mgitmap macro run %s\033[0m (or \033[1;96mgitmap %s\033[0m)\n\n", m.Name, m.Name)
}

func isMacroHelpArg(name string) bool {
	return name == "help" || name == "--help" || name == "-h"
}
