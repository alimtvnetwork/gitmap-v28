package cmdsee

import (
	"fmt"
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
	if len(args) == 0 {
		return PrintSeeHelp()
	}
	joined := strings.ToLower(strings.Join(args, " "))
	return routeSeeCommand(args, joined)
}

func routeSeeCommand(args []string, joined string) *apperror.AppError {
	if strings.HasPrefix(joined, "commit pending") {
		return runSeeCommitPending()
	}
	if hasIgnorePrefix(joined) {
		return routeIgnoreIssues(args)
	}
	if hasErrorsSSHPrefix(joined) {
		return runSeeErrorsSSH(args)
	}
	return routeRemaining(args, joined)
}

func routeRemaining(args []string, joined string) *apperror.AppError {
	if strings.HasPrefix(joined, "errors") {
		err := cmderrors.RunErrorsCLI(args[1:])
		return apperror.WrapSimple(err, "errors")
	}
	if strings.HasPrefix(joined, "history") {
		return runSeeHistory(args[1:])
	}
	return apperror.NewSimple("unknown see command: "+joined, "E1020")
}

func hasIgnorePrefix(joined string) bool {
	return strings.HasPrefix(joined, "git-ignore") || strings.HasPrefix(joined, "ignore") || strings.HasPrefix(joined, "ig issues")
}

func routeIgnoreIssues(args []string) *apperror.AppError {
	err := cmdignore.RunIgnoreCLI(append([]string{"scan"}, args[1:]...))
	return apperror.WrapSimple(err, "ignore")
}

func hasErrorsSSHPrefix(joined string) bool {
	return strings.HasPrefix(joined, "errors ssh") || joined == "ses"
}

func runSeeCommitPending() *apperror.AppError {
	dirty := cmdcpar.CollectDirtyRepositories()
	if len(dirty) == 0 {
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
	if HistoryRunnerFn != nil {
		return HistoryRunnerFn(args)
	}
	err := exec.Command("gitmap", append([]string{"history"}, args...)...).Run()
	return apperror.WrapSimple(err, "history")
}

// RunSeeErrorsSSH executes see errors across local host and remote fleet nodes.
func RunSeeErrorsSSH(args []string) *apperror.AppError {
	err := cmdssh.RunFleetPASCommand("see errors", "gitmap errors --limit 10", func() error {
		return cmderrors.RunErrorsCLI(args)
	})
	return apperror.WrapSimple(err, "errors-ssh")
}

func runSeeErrorsSSH(args []string) *apperror.AppError {
	subArgs := []string{}
	if len(args) > 2 {
		subArgs = args[2:]
	}
	return RunSeeErrorsSSH(subArgs)
}

// RunRepoManageUI opens the repository management web UI.
func RunRepoManageUI() *apperror.AppError {
	fmt.Println("Launching GitMap Repository Management UI...")
	cmd := exec.Command("gitmap", "ui")
	err := cmd.Start()
	return apperror.WrapSimple(err, "repo-manage-ui")
}

// PrintSeeHelp displays help for gitmap see suite.
func PrintSeeHelp() *apperror.AppError {
	fmt.Printf(`
  GitMap Inspection Suite (see)
    • gitmap see commit pending          - List repositories with uncommitted changes
    • gitmap see git-ignore issues       - List repositories with ignore/duplicate issues
    • gitmap see errors                  - Display recent execution errors
    • gitmap see history                 - Display command history
    • gitmap see errors ssh (ses)        - Display errors across SSH fleet (PAS Formula)
    • gitmap repo-manage ui              - Open web repository manager
`)
	return nil
}
