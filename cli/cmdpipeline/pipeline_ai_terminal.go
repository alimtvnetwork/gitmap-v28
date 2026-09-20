package cmdpipeline

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderPipelineAIStatusTerminal(p PipelineStatusPayload) {
	renderPipelineStatusTerminal(p)
	if p.HasErrors {
		renderPipelineAILiveErrorAlert(p)
		renderPipelineAINextActionError(p)
		return
	}

	if p.IsRunning {
		renderPipelineAINextActionWaiting(p)
	}
}

func renderPipelineAILiveErrorAlert(p PipelineStatusPayload) {
	fmt.Println()
	fmt.Printf("  %s%s🛑 [LIVE CI/CD ERROR DETECTED - STOP WAITING]%s\n",
		constants.ColorBold, constants.ColorRed, constants.ColorReset)
	renderPipelineAIFailingJobSnippet(p)
}

func renderPipelineAIFailingJobSnippet(p PipelineStatusPayload) {
	if len(p.FailedJobs) > 0 {
		fj := p.FailedJobs[0]
		fmt.Printf("  %s● Failing Job / Step:%s %s / %s\n",
			constants.ColorCyan, constants.ColorReset, fj.JobName, fj.StepName)
		if len(fj.FailureSummary) > 0 {
			fmt.Printf("    %s%s%s\n", constants.ColorYellow, fj.FailureSummary, constants.ColorReset)
		}
	}

	renderActionableSnippetLines(p.ActionableErrorSnippet)
}

func renderActionableSnippetLines(snippet string) {
	trimmed := strings.TrimSpace(snippet)
	if len(trimmed) == 0 {
		return
	}

	fmt.Printf("  %s● Error Diagnostics Snippet:%s\n", constants.ColorRed, constants.ColorReset)
	lines := strings.Split(trimmed, "\n")
	capped := capErrorSnippetLines(lines, 8)
	for _, line := range capped {
		fmt.Printf("    %s%s%s\n", constants.ColorDim, line, constants.ColorReset)
	}
}

func capErrorSnippetLines(lines []string, maxCount int) []string {
	if len(lines) <= maxCount {
		return lines
	}

	return lines[:maxCount]
}

func renderPipelineAINextActionError(p PipelineStatusPayload) {
	fmt.Println()
	fmt.Printf("  %s🞠 AI Automation Next Action:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("     Run: %s%s%s\n", constants.ColorGreen, p.NextAiCommand, constants.ColorReset)
	fmt.Printf("     %s(Stop waiting for remaining jobs. Fix current failure now so CI/CD can rerun in parallel.)%s\n",
		constants.ColorYellow, constants.ColorReset)
	fmt.Println()
}

func renderPipelineAINextActionWaiting(p PipelineStatusPayload) {
	fmt.Println()
	fmt.Printf("  %s🞠 AI Automation Next Action:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("     Run: %s%s%s\n", constants.ColorGreen, p.NextAiCommand, constants.ColorReset)
	fmt.Printf("     (Automatically delays %s then queries status)\n", formatEtaDisplay(p.EtaSeconds))
	fmt.Println()
}

func renderPipelineAIErrorsTerminal(report string, hasFailures bool) {
	if len(report) > 0 {
		fmt.Println(report)
	}

	fmt.Println()
	fmt.Printf("  %s🞠 AI Remediation Advice:%s\n", constants.ColorCyan, constants.ColorReset)
	if hasFailures {
		fmt.Printf("     Failure detected in CI/CD pipeline.\n")
		fmt.Printf("     Run: %s%s%s\n", constants.ColorGreen, "gitmap pipeline fix agy", constants.ColorReset)
		fmt.Printf("     %s(Stop waiting for remaining jobs. Fix current failure now so CI/CD can rerun in parallel.)%s\n\n",
			constants.ColorYellow, constants.ColorReset)
		return
	}

	fmt.Printf("     No active or recent pipeline failures detected.\n")
	fmt.Printf("     %sAll workflow steps are clean or in progress.%s\n\n", constants.ColorGreen, constants.ColorReset)
}
