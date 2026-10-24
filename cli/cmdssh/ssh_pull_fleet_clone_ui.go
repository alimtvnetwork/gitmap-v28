package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RenderFleetPostPullUnresolvedList renders an identical clone of local post-pull unresolved issues.
func RenderFleetPostPullUnresolvedList(issues []FleetNodeIssue) {
	total := CountTotalFleetIssues(issues)
	fmt.Printf("\n  %s[!]%s %s%d multi-machine node(s) / %d issue(s) require attention after pull:%s\n",
		constants.ColorYellow, constants.ColorReset,
		constants.ColorBold, len(issues), total, constants.ColorReset)
	for _, issue := range issues {
		renderSingleNodeIssue(issue)
	}
}

func renderSingleNodeIssue(issue FleetNodeIssue) {
	nodeLabel := "Remote Node"
	color := constants.ColorYellow
	if issue.IsLocal {
		nodeLabel = "Local Machine"
		color = constants.ColorGreen
	}
	fmt.Printf("    %s▶ %s [%s] (%s):%s\n", color, nodeLabel, issue.NodeName, issue.IP, constants.ColorReset)
	renderSingleNodeIssueDetails(issue)
}

func renderSingleNodeIssueDetails(issue FleetNodeIssue) {
	if issue.HasNodeError {
		fmt.Printf("        %s• Node error:%s %s%s%s\n",
			constants.ColorBold, constants.ColorReset,
			constants.ColorDim, issue.NodeError, constants.ColorReset)
		return
	}
	for _, r := range issue.RepoIssues {
		fmt.Printf("        • %s%s%s: %s%s%s\n",
			constants.ColorBold, r.RepoName, constants.ColorReset,
			constants.ColorDim, r.Reason, constants.ColorReset)
	}
}

// RenderFleetPostPullMenu presents the interactive remediation choices cloned from local gitmap pa.
func RenderFleetPostPullMenu() {
	fmt.Printf("\n  %sSome multi-machine nodes require attention after pull. How would you like to resolve?%s\n",
		constants.ColorBold, constants.ColorReset)
	renderFleetMenuChoicesPrimary()
	renderFleetMenuChoicesSecondary()
	fmt.Printf("  %sChoice [1/2/3/4/5/6/7/h/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func renderFleetMenuChoicesPrimary() {
	fmt.Printf("    %s[1/s/send]%s    %sSend commits across affected nodes%s  - Commit & push changes on affected nodes\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %s[2/a/all]%s     %sSend all (All Fleet Repositories)%s   - Commit & push across all fleet repositories\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %s[3/t/stash]%s   %sStash All across affected nodes%s     - Run %sgit stash%s across dirty repositories\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset, constants.ColorBold, constants.ColorReset)
}

func renderFleetMenuChoicesSecondary() {
	fmt.Printf("    %s[4/f/fix]%s     %sFix All across affected nodes%s       - Run %sgitmap fix --all%s across failed nodes\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	fmt.Printf("    %s[5/r/retry]%s   %sRetry Pull-All%s                      - Re-run %sgitmap pa%s across affected nodes\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorMagenta, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	fmt.Printf("    %s[6/i/step]%s    %sStep through node-by-node%s           - Tree-view & per-node remediation choices\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("    %s[h/?/help]%s    %sHelp & Explanations%s                 - Multi-machine remediation guide & strategies\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %s[7/q/skip]%s    %sSkip and leave as-is%s                - No further action taken\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorDim, constants.ColorReset)
}

// RenderFleetPostPullHelpGuide displays clear explanations for all multi-machine remediation strategies.
func RenderFleetPostPullHelpGuide() {
	fmt.Printf("\n  %s=== Multi-Machine Fleet Remediation Guide ===%s\n", constants.ColorBold, constants.ColorReset)
	renderFleetHelpOptionSend()
	renderFleetHelpOptionFix()
	renderFleetHelpOptionRetry()
	fmt.Println()
}

func renderFleetHelpOptionSend() {
	fmt.Printf("  • %sSend Commits [1/2]%s: Runs %sgitmap sends cp all%s on target machines.\n",
		constants.ColorBold, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Println("    Safely commits untracked/dirty changes with WIP messages and pushes to remote.")
	fmt.Printf("  • %sStash All [3]%s: Runs %sgit stash%s on target machines.\n",
		constants.ColorBold, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Println("    Saves local working state to stash queue so git pull can proceed cleanly.")
}

func renderFleetHelpOptionFix() {
	fmt.Printf("  • %sFix All [4]%s: Runs %sgitmap fix --all%s on target machines.\n",
		constants.ColorBold, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Println("    Normalizes README casing, recovers detached HEADs, stashes, and retries.")
}

func renderFleetHelpOptionRetry() {
	fmt.Printf("  • %sRetry Pull-All [5]%s: Re-runs %sgitmap pa%s on affected machines to re-verify status.\n",
		constants.ColorBold, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("  • %sStep Through [6]%s: Lets you decide actions interactively node-by-node.\n",
		constants.ColorBold, constants.ColorReset)
}
