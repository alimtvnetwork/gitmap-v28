// Package cmd — schedule_tree.go renders schedule composition tree views.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printScheduleTree is an alias for printScheduleSummaryTree.
func printScheduleTree(taskName, interval, shellType string, steps []string) {
	printScheduleSummaryTree(taskName, interval, shellType, steps)
}

// printScheduleSummaryTree renders a structured tree view for a scheduled task.
func printScheduleSummaryTree(taskName, interval, shellType string, steps []string) {
	printScheduleTreeHeader(taskName, interval, shellType)
	renderScheduleStepList(steps)
	fmt.Println()
}

// printScheduleTreeHeader outputs the root node with task name and interval/shell details.
func printScheduleTreeHeader(taskName, interval, shellType string) {
	details := formatScheduleDetails(interval, shellType)
	if details != "" {
		fmt.Printf("  %s%s%s %s%s%s\n",
			constants.ColorWhite, taskName, constants.ColorReset,
			constants.ColorDim, details, constants.ColorReset)
		return
	}
	fmt.Printf("  %s%s%s\n",
		constants.ColorWhite, taskName, constants.ColorReset)
}

// formatScheduleDetails formats the interval and shell metadata string.
func formatScheduleDetails(interval, shellType string) string {
	hasInterval := interval != ""
	hasShell := shellType != ""
	if hasInterval && hasShell {
		return fmt.Sprintf("(interval: %s, shell: %s)", interval, shellType)
	}
	if hasInterval {
		return fmt.Sprintf("(interval: %s)", interval)
	}
	if hasShell {
		return fmt.Sprintf("(shell: %s)", shellType)
	}
	return ""
}

// renderScheduleStepList iterates over schedule steps and prints branch lines.
func renderScheduleStepList(steps []string) {
	stepCount := len(steps)
	if stepCount == 0 {
		return
	}
	for stepIndex, stepText := range steps {
		isLastStep := stepIndex == stepCount-1
		renderScheduleStepNode(stepText, isLastStep)
	}
}

// renderScheduleStepNode renders an individual step branch line.
func renderScheduleStepNode(stepText string, isLastStep bool) {
	connector := constants.TreeBranch
	if isLastStep {
		connector = constants.TreeCorner
	}
	fmt.Printf("  %s%s%s %s%s%s\n",
		constants.ColorCyan, connector, constants.ColorReset,
		constants.ColorWhite, stepText, constants.ColorReset)
}
