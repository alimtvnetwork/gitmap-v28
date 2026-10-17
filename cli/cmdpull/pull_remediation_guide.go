package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func printRemediationPromptMenu() {
	fmt.Printf("\n  %s[?]%s %sRemediate failed repositories?%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, constants.ColorReset)
	fmt.Println("    [1/a] Resolve all failed repositories at once (auto-merge)")
	fmt.Println("    [2/s] Step through one by one")
	fmt.Println("    [3/r] Rebase all failed repositories (pull --rebase)")
	fmt.Println("    [h/?/help] Help & explanations of resolution options")
	fmt.Println("    [q/n] Skip / Exit")
	fmt.Printf("  %sChoice [1/2/3/h/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func isChoiceHelp(choice string) bool {
	return choice == "h" || choice == "help" || choice == "?"
}

func renderRemediationHelpGuide() {
	renderRemediationHelpHeader()
	renderAutoMergeExplanation()
	renderStepThroughExplanation()
	renderRebaseExplanation()
}

func renderRemediationHelpHeader() {
	fmt.Printf("\n  %s╔══════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║%s  %s📖 Remediation Options & Safety Guide%s                              %s║%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func renderAutoMergeExplanation() {
	fmt.Printf("  %s[1/a] Auto-Merge (Batch Resolution):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Operation: Attempts git pull with automatic conflict detection and")
	fmt.Println("      credential auto-conversion (HTTPS -> SSH) across all failed repos.")
	fmt.Println("    • Safety Guarantee: Non-destructive. Stashes dirty worktrees when needed.")
	fmt.Println("      Never force-pushes or discards uncommitted code.")
	fmt.Println("    • Risk Profile: Low. If merge conflicts arise, repository is safely")
	fmt.Println("      aborted and flagged for manual resolution.")
	fmt.Println()
}

func renderStepThroughExplanation() {
	fmt.Printf("  %s[2/s] Step-Through (Interactive Resolution):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Operation: Interactively iterates through each failure one by one,")
	fmt.Println("      allowing choice of auto-merge, rebase, skip, or quit per repository.")
	fmt.Println("    • Safety Guarantee: Complete operator oversight for every action.")
	fmt.Println("    • Risk Profile: Lowest. Recommended when repositories require different")
	fmt.Println("      handling strategies.")
	fmt.Println()
}

func renderRebaseExplanation() {
	fmt.Printf("  %s[3/r] Rebase (git pull --rebase):%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("    • Operation: Replays local unpushed commits on top of incoming remote")
	fmt.Println("      commits to maintain a linear git history.")
	fmt.Println("    • Safety Guarantee: Automatically aborts (git rebase --abort) on failure,")
	fmt.Println("      restoring repository to original working state.")
	fmt.Println("    • Risk Profile: Moderate. Rewrites local commit history. Not suited")
	fmt.Println("      for repositories with complex branch merges.")
	fmt.Println()
}

func readPromptChoice() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "q"
	}

	return strings.ToLower(strings.TrimSpace(line))
}

func dispatchRemediationChoice(choice string, failures []PullFailureSummary) error {
	if isChoiceResolveAll(choice) {
		return RunBatchRemediateAll(failures)
	}
	if isChoiceStepByStep(choice) {
		return RunStepByStepRemediate(failures)
	}
	if isChoiceRebaseAll(choice) {
		return RunBatchRemediateRebase(failures)
	}

	return nil
}

func isChoiceResolveAll(choice string) bool {
	return choice == "1" || choice == "a" || choice == "all" || choice == "y" || choice == "yes"
}

func isChoiceStepByStep(choice string) bool {
	return choice == "2" || choice == "s" || choice == "step"
}

func isChoiceRebaseAll(choice string) bool {
	return choice == "3" || choice == "r" || choice == "rebase"
}
