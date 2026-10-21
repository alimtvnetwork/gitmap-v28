package cmdpull

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderPostPullUnresolvedList(unresolved []UnresolvedPullRepo) {
	fmt.Printf("\n  %s[!]%s %s%d repository(ies) require attention after pull:%s\n",
		constants.ColorYellow, constants.ColorReset,
		constants.ColorBold, len(unresolved), constants.ColorReset)

	for _, item := range unresolved {
		reason := "working tree dirty / autostash conflict"
		if item.IsFailed && item.ErrorMsg != "" {
			reason = item.ErrorMsg
		}
		fmt.Printf("      • %s%s%s: %s%s%s\n",
			constants.ColorBold, item.RepoName, constants.ColorReset,
			constants.ColorDim, reason, constants.ColorReset)
	}
}

func renderPostPullMenu() {
	fmt.Printf("\n  %sSome repositories require attention after pull. How would you like to resolve?%s\n", constants.ColorBold, constants.ColorReset)
	renderPostPullMenuChoicesPrimary()
	renderPostPullMenuChoicesSecondary()
	fmt.Printf("  %sChoice [1/2/3/4/5/6/h/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func renderPostPullMenuChoicesPrimary() {
	fmt.Printf("    %s[1/s/send]%s    %sSend commits (Current Machine)%s  - Commit & push on current machine (asks to send all after)\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %s[2/a/all]%s     %sSend all (All Repositories)%s     - Commit & push across all fleet repositories\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %s[3/t/stash]%s   %sStash All%s                      - Run %sgit stash%s across dirty repositories\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset, constants.ColorBold, constants.ColorReset)
}

func renderPostPullMenuChoicesSecondary() {
	fmt.Printf("    %s[4/f/fix]%s     %sFix All%s                        - Run %sgitmap fix --all%s across failed repositories\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	fmt.Printf("    %s[5/i/step]%s    %sStep through individually%s      - Tree-view diff & colored action previews\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("    %s[h/?/help]%s    %sHelp & Explanations%s            - Detailed guide on resolution strategies & risks\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %s[6/q/skip]%s    %sSkip and leave as-is%s           - No further action taken\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorDim, constants.ColorReset)
}

func dispatchPostPullChoice(choice string, unresolved []UnresolvedPullRepo) error {
	currentChoice := choice
	for {
		res, isHandled := tryDispatchAnyPostPullChoice(currentChoice, unresolved)
		if isHandled {
			return res
		}
		if !isHelpChoice(currentChoice) {
			return nil
		}
		currentChoice = handlePostPullHelpLoop()
	}
}

func handlePostPullHelpLoop() string {
	renderPostPullHelpGuide()
	renderPostPullMenu()

	return readPromptUserChoice()
}

func tryDispatchAnyPostPullChoice(choice string, unresolved []UnresolvedPullRepo) (error, bool) {
	if res, isHandled := tryDispatchPrimaryPostPullChoice(choice, unresolved); isHandled {
		return res, true
	}

	return tryDispatchSecondaryPostPullChoice(choice, unresolved)
}

func tryDispatchPrimaryPostPullChoice(choice string, unresolved []UnresolvedPullRepo) (error, bool) {
	switch choice {
	case "1", "s", "send", "c", "commit", "current", "cm", "y", "yes":
		return executeSendCommitsCurrentMachine(unresolved), true
	case "2", "a", "all", "send-all", "fleet", "sa":
		return executeSendAllCommits(), true
	case "3", "t", "stash", "stash-all":
		return executeStashAll(unresolved), true
	}

	return nil, false
}

func tryDispatchSecondaryPostPullChoice(choice string, unresolved []UnresolvedPullRepo) (error, bool) {
	switch choice {
	case "4", "f", "fix", "fix-all":
		return executeFixAll(), true
	case "5", "i", "step", "single":
		return executeStepByStepRemediation(unresolved), true
	}

	return nil, false
}

func isHelpChoice(choice string) bool {
	return choice == "h" || choice == "help" || choice == "?"
}

func renderPostPullHelpGuide() {
	renderPostPullHelpHeader()
	renderPostPullOptionDetails()
}

func renderPostPullHelpHeader() {
	fmt.Printf("\n  %s╔══════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║%s  %s📖 Post-Pull Resolution & Safety Guide%s                            %s║%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func renderPostPullOptionDetails() {
	renderPostPullOptionDetailsPrimary()
	renderPostPullOptionDetailsSecondary()
}

func renderPostPullOptionDetailsPrimary() {
	fmt.Printf("  %s[1/s/send] Send Commits (Current Machine):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Runs gitmap sends cp <repo> \"wip: save changes\" for each dirty repository on this machine.")
	fmt.Println("    • Safety: Automatically stages, commits with semantic categories, and pushes to remote.")
	fmt.Printf("\n  %s[2/a/all] Send All (All Repositories):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Runs gitmap sends cp all \"wip: save changes\" across the entire workspace fleet.")
	fmt.Println("    • Safety: Discovers all repos with changes, commits and pushes concurrently.")
	fmt.Printf("\n  %s[3/t/stash] Stash All:%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Runs git stash across dirty repositories to clear working trees.")
	fmt.Println("    • Safety: Stashed changes remain restorable via git stash pop.")
}

func renderPostPullOptionDetailsSecondary() {
	fmt.Printf("\n  %s[4/f/fix] Fix All:%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Runs gitmap fix --all across failed repositories to heal remotes and resolve errors.")
	fmt.Println("    • Safety: Non-destructive; aborts on conflict without losing changes.")
	fmt.Printf("\n  %s[5/i/step] Step through individually:%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Inspects each unresolved repository with diff previews, VS Code merge editor, or stash.")
	fmt.Println("    • Safety: Granular manual review and verification per repository.")
	fmt.Printf("\n  %s[6/q/skip] Skip and leave as-is:%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Exits menu without modifying any repositories or files.")
}
