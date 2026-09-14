package macro

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintExecutionSummaryTree outputs the formatted execution tree to stdout.
func PrintExecutionSummaryTree(rep *ExecutionReport) {
	if rep == nil {
		return
	}

	treeStr := RenderExecutionSummaryTree(rep)
	fmt.Println()
	fmt.Print(treeStr)
	fmt.Println()
}

// RenderExecutionSummaryTree builds the visual hierarchical summary of executed macro steps.
func RenderExecutionSummaryTree(rep *ExecutionReport) string {
	if rep == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(formatSummaryRootHeader(rep))
	sb.WriteString("\n")
	stepCount := len(rep.Steps)
	for i, step := range rep.Steps {
		isLast := i == stepCount-1
		sb.WriteString(formatSummaryStepLine(step, isLast))
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatSummaryRootHeader(rep *ExecutionReport) string {
	badge := formatRootBadge(rep)
	title := fmt.Sprintf("%sMacro %s%s", constants.ColorWhite, strconv.Quote(rep.Macro), constants.ColorReset)
	metrics := formatRootMetrics(rep)

	return fmt.Sprintf("  %s %s %s", badge, title, metrics)
}

func formatRootBadge(rep *ExecutionReport) string {
	if rep.FailedSteps > 0 {
		return fmt.Sprintf("%s[✖]%s", constants.ColorRed, constants.ColorReset)
	}

	return fmt.Sprintf("%s[✔]%s", constants.ColorGreen, constants.ColorReset)
}

func formatRootMetrics(rep *ExecutionReport) string {
	if rep.FailedSteps > 0 {
		return fmt.Sprintf("%s(%d/%d completed · %d failed · %.1fs)%s",
			constants.ColorDim, rep.ExecutedSteps, rep.TotalSteps, rep.FailedSteps, rep.ElapsedSeconds, constants.ColorReset)
	}

	return fmt.Sprintf("%s(%d/%d steps · %.1fs · ok)%s",
		constants.ColorDim, rep.ExecutedSteps, rep.TotalSteps, rep.ElapsedSeconds, constants.ColorReset)
}

func formatSummaryStepLine(step StepExecution, isLast bool) string {
	connector := selectSummaryConnector(isLast)
	badge := formatStepStatusBadge(step.Status)
	meta := fmt.Sprintf("%s(code %d, %.1fs)%s", constants.ColorDim, step.ExitCode, step.ElapsedSeconds, constants.ColorReset)
	cmdText := fmt.Sprintf("%s%s%s", constants.ColorWhite, step.CommandLine, constants.ColorReset)
	line := fmt.Sprintf("  %s %s %s %s", connector, badge, meta, cmdText)
	if step.FailureLogFile != "" {
		line += fmt.Sprintf("  %s(log: %s%s%s)%s", constants.ColorDim, constants.ColorYellow, step.FailureLogFile, constants.ColorDim, constants.ColorReset)
	}
	if step.WorkingDir != "" {
		line += fmt.Sprintf(" %s(dir: %s)%s", constants.ColorDim, step.WorkingDir, constants.ColorReset)
	}

	return line
}

func selectSummaryConnector(isLast bool) string {
	if isLast {
		return fmt.Sprintf("%s%s%s", constants.ColorCyan, constants.TreeCorner, constants.ColorReset)
	}

	return fmt.Sprintf("%s%s%s", constants.ColorCyan, constants.TreeBranch, constants.ColorReset)
}

func formatStepStatusBadge(status string) string {
	switch strings.ToLower(status) {
	case "success", "ok":
		return fmt.Sprintf("%s[✔] ok%s", constants.ColorGreen, constants.ColorReset)
	case "timeout":
		return fmt.Sprintf("%s[!] timeout%s", constants.ColorYellow, constants.ColorReset)
	case "dry-run":
		return fmt.Sprintf("%s[-] dry-run%s", constants.ColorDim, constants.ColorReset)
	default:
		return fmt.Sprintf("%s[✖] failed%s", constants.ColorRed, constants.ColorReset)
	}
}
