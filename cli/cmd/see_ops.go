package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcpar"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

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
