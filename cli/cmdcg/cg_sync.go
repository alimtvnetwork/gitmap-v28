// Package cmdcg provides unified synchronization for coding guidelines and prompt architecture.
package cmdcg

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	cgSyncHeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	cgSyncSuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
)

func printCGSyncHeader(repoCount int, isDryRun bool) {
	dryLabel := ""
	if isDryRun {
		dryLabel = " (DRY RUN)"
	}

	fmt.Println(cgSyncHeaderStyle.Render(fmt.Sprintf("=== Coding Guidelines & Prompt Architect Sync%s ===", dryLabel)))
	fmt.Printf("Synchronizing across %d target repository(ies)...\n\n", repoCount)
}

func printCGSyncStep(stepName string) {
	fmt.Printf("► %s\n", cgAuditTitleStyle.Render(stepName))
}

func printCGSyncSummary(repoCount int) {
	fmt.Println()
	fmt.Println(cgSyncSuccessStyle.Render("✔ Unified sync completed successfully."))
	fmt.Printf("  Synchronized %d repository(ies) with canonical standards and prompts.\n", repoCount)
}

func runSyncSteps(repos []string, isDryRun bool) error {
	printCGSyncStep("Step 1: Updating Coding Guidelines Specifications")
	runCGUpdateAction(repos)

	printCGSyncStep("Step 2: Synchronizing Prompt Architect")

	return runCGInstallPrompts(repos, isDryRun)
}

// RunCGSync executes a full synchronization of coding guidelines and prompt architect across repos.
func RunCGSync(repos []string, isDryRun bool) error {
	printCGSyncHeader(len(repos), isDryRun)

	errSync := runSyncSteps(repos, isDryRun)
	if errSync != nil {
		return errSync
	}

	printCGSyncSummary(len(repos))

	return nil
}
