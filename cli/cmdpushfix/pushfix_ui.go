package cmdpushfix

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintPushFixHeader displays the diagnosis and remediation HUD header.
func PrintPushFixHeader(repoDir string) {
	fmt.Printf("\n  %s%s%s %s[PUSH-FIX]%s %sAutonomous Push & Auth Recovery%s\n",
		constants.ColorBold, constants.ColorCyan, "⚡",
		constants.ColorBold, constants.ColorReset,
		constants.ColorDim, constants.ColorReset,
	)
	fmt.Printf("  Target Directory: %s%s%s\n\n", constants.ColorDim, repoDir, constants.ColorReset)
}

// PrintPhaseStep outputs a numbered execution stage card.
func PrintPhaseStep(step int, name, detail string) {
	fmt.Printf("  [%d/4] %s%-24s%s • %s\n",
		step,
		constants.ColorCyan, name, constants.ColorReset,
		detail,
	)
}

// PrintRemediationApplied prints an auto-fix action badge.
func PrintRemediationApplied(action string) {
	fmt.Printf("        %s✔ [FIX APPLIED]%s %s\n",
		constants.ColorGreen, constants.ColorReset,
		action,
	)
}

// PrintPushFixSuccess renders the final completion summary card.
func PrintPushFixSuccess(remote, branch, sha string, isDryRun bool) {
	status := "Pushed to GitHub successfully"
	if isDryRun {
		status = "Simulated push (dry-run mode) — ready for live push"
	}
	fmt.Printf("\n  %s✔ %s%s\n", constants.ColorGreen, status, constants.ColorReset)
	fmt.Printf("    • Remote: %s%s%s\n", constants.ColorCyan, remote, constants.ColorReset)
	fmt.Printf("    • Branch: %s%s%s\n", constants.ColorCyan, branch, constants.ColorReset)
	if len(sha) > 0 {
		fmt.Printf("    • Commit: %s%s%s\n", constants.ColorDim, sha, constants.ColorReset)
	}
	fmt.Println()
}
