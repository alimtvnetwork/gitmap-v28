package cmdssh

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	"golang.org/x/term"
)

// PromptFleetPostPullRemediation presents an interactive resolution circle when fleet issues exist.
func PromptFleetPostPullRemediation(outcomes []FleetNodePullOutcome, conns []db.SSHConnection) error {
	if !IsInteractiveFleetPromptEligible() {
		return nil
	}
	issues := ExtractFleetIssues(outcomes)
	if !HasAnyFleetIssues(issues) {
		return nil
	}

	return runFleetRemediationCircle(issues, conns)
}

func runFleetRemediationCircle(initialIssues []FleetNodeIssue, conns []db.SSHConnection) error {
	issues := initialIssues
	for HasAnyFleetIssues(issues) {
		RenderFleetPostPullUnresolvedList(issues)
		RenderFleetPostPullMenu()
		choice := readFleetPromptUserChoice()
		if isQuitChoice(choice) {
			fmt.Printf("\n  %sSkipping multi-machine remediation.%s\n\n", constants.ColorDim, constants.ColorReset)
			return nil
		}
		if isHelpChoice(choice) {
			RenderFleetPostPullHelpGuide()
			continue
		}
		issues = handleFleetChoiceInCircle(choice, issues, conns)
	}
	printAllFleetIssuesResolvedBanner()

	return nil
}

func isQuitChoice(choice string) bool {
	return choice == "7" || choice == "q" || choice == "skip" || choice == "exit"
}

func isHelpChoice(choice string) bool {
	return choice == "h" || choice == "?" || choice == "help"
}

func handleFleetChoiceInCircle(choice string, issues []FleetNodeIssue, conns []db.SSHConnection) []FleetNodeIssue {
	isHandled := dispatchFleetRemediationAction(choice, issues, conns)
	if !isHandled {
		fmt.Printf("  %sUnrecognized choice '%s'.%s\n", constants.ColorYellow, choice, constants.ColorReset)
		return issues
	}
	if choice == "5" || choice == "r" || choice == "retry" {
		return issues
	}

	return promptRetryCircleDecision(issues, conns)
}

func dispatchFleetRemediationAction(choice string, issues []FleetNodeIssue, conns []db.SSHConnection) bool {
	switch choice {
	case "1", "s", "send":
		ExecuteFleetSendCommitsAction(issues, conns)
	case "2", "a", "all":
		ExecuteFleetSendAllAction(conns)
	case "3", "t", "stash":
		ExecuteFleetStashAllAction(issues, conns)
	case "4", "f", "fix":
		ExecuteFleetFixAllAction(issues, conns)
	case "5", "r", "retry":
		_ = ExecuteFleetRetryPullAction(issues, conns)
	case "6", "i", "step":
		ExecuteFleetStepThroughAction(issues, conns)
	default:
		return false
	}

	return true
}

func promptRetryCircleDecision(issues []FleetNodeIssue, conns []db.SSHConnection) []FleetNodeIssue {
	fmt.Printf("  %sAction executed. Would you like to retry pull-all now? [Y/n/circle]:%s ",
		constants.ColorCyan, constants.ColorReset)
	decision := readFleetPromptUserChoice()
	if decision == "n" || decision == "no" || decision == "q" {
		return nil
	}
	if decision == "circle" || decision == "c" || decision == "m" {
		return issues
	}

	return ExecuteFleetRetryPullAction(issues, conns)
}

func printAllFleetIssuesResolvedBanner() {
	fmt.Printf("\n  %s✓ All multi-machine fleet issues successfully resolved!%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func readFleetPromptUserChoice() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "q"
	}

	return strings.ToLower(strings.TrimSpace(line))
}

// IsInteractiveFleetPromptEligible verifies terminal TTY and standard input availability.
func IsInteractiveFleetPromptEligible() bool {
	if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") != "" || os.Getenv("GITMAP_JSON") != "" {
		return false
	}
	if term.IsTerminal(int(os.Stdin.Fd())) || term.IsTerminal(int(os.Stdout.Fd())) {
		return true
	}
	if termout.IsStdoutTTY() {
		return true
	}
	info, err := os.Stdin.Stat()

	return err == nil && (info.Mode()&os.ModeCharDevice) != 0
}
