// Package cmd — fixgit_ui.go: terminal rendering and reporting for fix-git.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func renderFixGitOutput(res FixGitResult, isJSON bool) error {
	if isJSON {
		return renderFixGitJSON(res)
	}

	renderFixGitTerminal(res)

	return nil
}

func renderFixGitJSON(res FixGitResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(res)
}

func renderFixGitTerminal(res FixGitResult) {
	printFixGitHeader(res.TargetDir)

	if res.IsClean && len(res.Issues) == 0 {
		fmt.Printf("  %s✓ Git repository is completely healthy.%s\n", constants.ColorGreen, constants.ColorReset)
		fmt.Printf("  %sAll permissions, lock files, and index structure nominal.%s\n\n", constants.ColorDim, constants.ColorReset)

		return
	}

	printFixGitIssuesList(res.Issues)
	printFixGitSummary(res)
}

func printFixGitHeader(targetDir string) {
	fmt.Println()
	fmt.Printf("  %s[FIX-GIT] Git Diagnostic & Self-Healing Engine%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %sRepository:%s %s%s%s\n\n", constants.ColorWhite, constants.ColorReset, constants.ColorCyan, targetDir, constants.ColorReset)
}

func printFixGitIssuesList(issues []FixGitIssue) {
	for _, issue := range issues {
		statusIcon := fmt.Sprintf("%s✓%s", constants.ColorGreen, constants.ColorReset)
		if !issue.IsFixed {
			statusIcon = fmt.Sprintf("%s!%s", constants.ColorYellow, constants.ColorReset)
		}

		catBadge := fmt.Sprintf("%s[%-13s]%s", constants.ColorMagenta, issue.Category, constants.ColorReset)
		fmt.Printf("  %s %s %s\n", statusIcon, catBadge, issue.Description)
		fmt.Printf("     %s↳ Action:%s %s\n", constants.ColorDim, constants.ColorReset, issue.Remedy)

		if issue.ErrorDetail != "" {
			fmt.Printf("     %s↳ Error:%s  %s%s%s\n", constants.ColorRed, constants.ColorReset, constants.ColorRed, issue.ErrorDetail, constants.ColorReset)
		}
	}

	fmt.Println()
}

func printFixGitSummary(res FixGitResult) {
	if res.IssuesFixed > 0 {
		fmt.Printf("  %s✓ Successfully remediated %d Git issue(s)!%s\n", constants.ColorGreen, res.IssuesFixed, constants.ColorReset)
		fmt.Printf("  %sRepository ready for staging and committing.%s\n\n", constants.ColorWhite, constants.ColorReset)

		return
	}

	fmt.Printf("  %sDiagnostic complete. Found %d issue(s).%s\n\n", constants.ColorYellow, res.IssuesFound, constants.ColorReset)
}
