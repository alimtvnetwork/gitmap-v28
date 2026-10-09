package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcpar"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// HistoryRunnerFn is a hook to delegate history execution from cmd package.
var HistoryRunnerFn func(args []string) *apperror.AppError

// RunSeeCLI routes gitmap see subcommands.
func RunSeeCLI(args []string) *apperror.AppError {
	isZero := len(args) == 0
	if isZero {
		return PrintSeeHelp()
	}
	joined := strings.ToLower(strings.Join(args, " "))
	return routeSeeCommand(args, joined)
}

func routeSeeCommand(args []string, joined string) *apperror.AppError {
	if isCommitPendingCmd(joined) {
		return runSeeCommitPending()
	}
	if hasIgnorePrefix(joined) {
		return routeIgnoreIssues(args)
	}
	if hasErrorsSSHPrefix(joined) {
		return runSeeErrorsSSH(args)
	}
	if hasHistorySSHPrefix(joined) {
		return runSeeHistorySSH(args)
	}
	return routeRemaining(args, joined)
}

func isCommitPendingCmd(joined string) bool {
	return strings.HasPrefix(joined, "commit pending") ||
		strings.HasPrefix(joined, "commit-pending") ||
		strings.HasPrefix(joined, "cp")
}

func hasIgnorePrefix(joined string) bool {
	return strings.HasPrefix(joined, "git-ignore") ||
		strings.HasPrefix(joined, "ignore") ||
		strings.HasPrefix(joined, "ig")
}

func hasErrorsSSHPrefix(joined string) bool {
	return strings.HasPrefix(joined, "errors ssh") ||
		strings.HasPrefix(joined, "ses") ||
		joined == "ses"
}

func hasHistorySSHPrefix(joined string) bool {
	return strings.HasPrefix(joined, "history ssh") ||
		strings.HasPrefix(joined, "nodes history") ||
		strings.HasPrefix(joined, "nodes histories")
}

func routeRemaining(args []string, joined string) *apperror.AppError {
	if strings.HasPrefix(joined, "errors") {
		err := cmderrors.RunErrorsCLI(args[1:])
		return apperror.WrapSimple(err, "errors")
	}
	if strings.HasPrefix(joined, "history") {
		return runSeeHistory(args[1:])
	}
	if strings.HasPrefix(joined, "repo-manage") || strings.HasPrefix(joined, "ui") {
		return RunRepoManageUI()
	}
	return apperror.NewSimple("unknown see command: "+joined, "E1020")
}

func routeIgnoreIssues(args []string) *apperror.AppError {
	var scanArgs []string
	isLonger := len(args) > 1
	if isLonger {
		scanArgs = filterIgnoreFlags(args[1:])
	}
	err := cmdignore.RunIgnoreCLI(append([]string{"scan"}, scanArgs...))
	return apperror.WrapSimple(err, "ignore")
}

func filterIgnoreFlags(subArgs []string) []string {
	var filtered []string
	for _, a := range subArgs {
		isIssueWord := strings.EqualFold(a, "issues") || strings.EqualFold(a, "issue")
		if !isIssueWord {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func runSeeCommitPending() *apperror.AppError {
	dirty := cmdcpar.CollectDirtyRepositories()
	isClean := len(dirty) == 0
	if isClean {
		fmt.Printf("%s✓ No repositories with pending commits.%s\n", constants.ColorGreen, constants.ColorReset)
		return nil
	}
	printPendingCommits(dirty)
	return nil
}

func printPendingCommits(dirty []cmdcpar.DirtyRepoSummary) {
	fmt.Printf("\n%s  === PENDING COMMITS (%d repositories) ===%s\n", constants.ColorYellow, len(dirty), constants.ColorReset)
	for _, d := range dirty {
		fmt.Printf("  • %-30s | modified: %d | untracked: %d | staged: %d\n",
			d.RepoName, d.Diagnosis.ModifiedCount, d.Diagnosis.UntrackedCount, d.Diagnosis.StagedCount)
	}
	fmt.Println()
}

func runSeeHistory(args []string) *apperror.AppError {
	hasHook := HistoryRunnerFn != nil
	if hasHook {
		return HistoryRunnerFn(args)
	}
	err := exec.Command("gitmap", append([]string{"history"}, args...)...).Run()
	return apperror.WrapSimple(err, "history")
}

func runSeeHistorySSH(args []string) *apperror.AppError {
	subArgs := resolveHistorySSHArgs(args)
	err := cmdssh.RunFleetPASCommand("history", "gitmap history --limit 10", func() error {
		hasHook := HistoryRunnerFn != nil
		if hasHook {
			return HistoryRunnerFn(subArgs)
		}
		return exec.Command("gitmap", append([]string{"history"}, subArgs...)...).Run()
	})
	return apperror.WrapSimple(err, "history-ssh")
}

func resolveHistorySSHArgs(args []string) []string {
	isLonger := len(args) > 2
	if isLonger {
		return args[2:]
	}
	return []string{}
}

// RunSeeErrorsSSH executes see errors across local host and remote fleet nodes.
func RunSeeErrorsSSH(args []string) *apperror.AppError {
	err := cmdssh.RunFleetPASCommand("see errors", "gitmap errors --limit 10", func() error {
		return cmderrors.RunErrorsCLI(args)
	})
	return apperror.WrapSimple(err, "errors-ssh")
}

func runSeeErrorsSSH(args []string) *apperror.AppError {
	subArgs := resolveErrorsSSHArgs(args)
	return RunSeeErrorsSSH(subArgs)
}

func resolveErrorsSSHArgs(args []string) []string {
	isLonger := len(args) > 2
	if isLonger {
		return args[2:]
	}
	return []string{}
}

// RunRepoManageUI opens the repository management interactive terminal dashboard.
func RunRepoManageUI() *apperror.AppError {
	fmt.Println("Launching GitMap Repository Management Dashboard...")
	cmd := exec.Command("gitmap", "interactive")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	hasErr := err != nil
	if hasErr {
		cmdWeb := exec.Command("gitmap", "ui")
		_ = cmdWeb.Start()
	}
	return nil
}

// PrintSeeHelp displays help for gitmap see suite.
func PrintSeeHelp() *apperror.AppError {
	fmt.Printf(`
  GitMap Inspection Suite (see / c)
    • gitmap see commit pending          - List repositories with uncommitted changes
    • gitmap see git-ignore issues       - List repositories with ignore/duplicate issues
    • gitmap see errors                  - Display recent execution errors
    • gitmap see history                 - Display command history
    • gitmap see errors ssh (ses)        - Display errors across SSH fleet (PAS Formula)
    • gitmap history ssh                 - Display history across SSH fleet (PAS Formula)
    • gitmap nodes history               - Display history across SSH fleet nodes
    • gitmap repo-manage ui              - Launch interactive terminal dashboard
`)
	return nil
}
