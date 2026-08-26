// Package cmd — macro_tree.go renders macro steps composition tree views.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

// printMacroTree loads a macro by name and renders its step composition tree.
func printMacroTree(macroName string) {
	loadedMacro, loadErr := macro.LoadMacro(macroName)
	if loadErr != nil || loadedMacro == nil {
		return
	}
	printMacroStepsTree(loadedMacro)
}

// printMacroStepsTree renders the hierarchical step tree for a given macro.
func printMacroStepsTree(targetMacro *macro.Macro) {
	if targetMacro == nil || len(targetMacro.Steps) == 0 {
		return
	}
	printMacroTreeHeader(targetMacro)
	renderMacroStepList(targetMacro.Steps)
	fmt.Println()
}

// printMacroTreeHeader prints the root macro title and description.
func printMacroTreeHeader(targetMacro *macro.Macro) {
	if targetMacro.Description != "" {
		fmt.Printf("  %s%s%s %s%s%s\n",
			constants.ColorWhite, targetMacro.Name, constants.ColorReset,
			constants.ColorDim, targetMacro.Description, constants.ColorReset)
		return
	}
	fmt.Printf("  %s%s%s\n",
		constants.ColorWhite, targetMacro.Name, constants.ColorReset)
}

// renderMacroStepList iterates and renders each macro step node.
func renderMacroStepList(steps []macro.MacroStep) {
	stepCount := len(steps)
	for stepIndex, step := range steps {
		isLastStep := stepIndex == stepCount-1
		renderMacroStepNode(step, isLastStep)
	}
}

// selectTreeConnector returns the branch glyph based on whether this is the final step.
func selectTreeConnector(isLastStep bool) string {
	if isLastStep {
		return constants.TreeCorner
	}
	return constants.TreeBranch
}

// printStepBranch prints a formatted tree step branch line.
func printStepBranch(connector, commandLine, stepDescription string) {
	if stepDescription != "" {
		fmt.Printf("  %s%s%s %s%s%s %s%s%s\n",
			constants.ColorCyan, connector, constants.ColorReset,
			constants.ColorWhite, commandLine, constants.ColorReset,
			constants.ColorDim, stepDescription, constants.ColorReset)
		return
	}
	fmt.Printf("  %s%s%s %s%s%s\n",
		constants.ColorCyan, connector, constants.ColorReset,
		constants.ColorWhite, commandLine, constants.ColorReset)
}

// renderMacroStepNode formats and outputs an individual step tree branch.
func renderMacroStepNode(step macro.MacroStep, isLastStep bool) {
	connector := selectTreeConnector(isLastStep)
	stepDescription := formatStepDescription(step)
	printStepBranch(connector, step.CommandLine, stepDescription)
}

// formatStepDescription builds auxiliary description text for a step.
func formatStepDescription(step macro.MacroStep) string {
	if step.WorkingDir != "" {
		return fmt.Sprintf("(dir: %s)", step.WorkingDir)
	}
	return ""
}
