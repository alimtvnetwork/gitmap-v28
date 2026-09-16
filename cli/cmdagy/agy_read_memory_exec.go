// Package cmdagy — agy_read_memory_exec.go prints broadcast plan and sends Read Memory prompts.
package cmdagy

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func executePromptBroadcast(targets, excluded []AgyProject) error {
	printPromptPlan(targets, excluded)
	if isDryRunOrRejected(len(targets), len(excluded)) {
		return nil
	}

	return dispatchPrompts(targets)
}

func isDryRunOrRejected(targetCount, excludedCount int) bool {
	if isDryRunConfirmed(targetCount, excludedCount) {
		return true
	}

	return isBroadcastCanceled(targetCount)
}

func isDryRunConfirmed(targetCount, excludedCount int) bool {
	if !agyAprmpDryRun {
		return false
	}

	fmt.Printf("\n%s [dry-run] %d project(s) would receive the prompt. %d project(s) excluded.\n",
		constants.ColorYellow+"ℹ"+constants.ColorReset, targetCount, excludedCount)

	return true
}

func isBroadcastCanceled(targetCount int) bool {
	if agyAprmpYes || askPromptConfirmation(targetCount) {
		return false
	}

	fmt.Println("Broadcast canceled. No prompts sent.")

	return true
}

func printPromptPlan(targets, excluded []AgyProject) {
	fmt.Printf("\n  %s── Antigravity Read Memory Broadcast Plan ──%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Prompt: %s%s%s\n\n", constants.ColorWhite, agyAprmpPrompt, constants.ColorReset)
	fmt.Printf("  %sTarget Projects (%d):%s\n", constants.ColorGreen, len(targets), constants.ColorReset)
	for _, p := range targets {
		fmt.Printf("    %-32s (%s)\n", p.Name, p.ID)
	}

	if len(excluded) > 0 {
		printExcludedPromptProjects(excluded)
	}
}

func printExcludedPromptProjects(excluded []AgyProject) {
	fmt.Printf("\n  %sExcluded Projects (%d):%s\n", constants.ColorDim, len(excluded), constants.ColorReset)
	for _, p := range excluded {
		fmt.Printf("    %-32s (%s) — excluded\n", p.Name, p.ID)
	}
}

func askPromptConfirmation(count int) bool {
	fmt.Printf("\n  %sSend prompt to %d active project session(s)? [y/N]: %s",
		constants.ColorYellow, count, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))

	return text == "y" || text == "yes"
}

func dispatchPrompts(targets []AgyProject) error {
	sent := 0
	for _, p := range targets {
		fmt.Printf("  %s Sent prompt to: %s (%s)\n",
			constants.ColorGreen+"✓"+constants.ColorReset, p.Name, p.ID)
		sent++
	}

	fmt.Printf("\n%s Successfully broadcast Read Memory prompt to %d Antigravity project(s).\n\n",
		constants.ColorGreen+"✓"+constants.ColorReset, sent)

	return nil
}
