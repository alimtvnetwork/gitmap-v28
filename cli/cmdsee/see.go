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
var HistoryRunnerFn func(args []string) error

// RunSeeCLI routes gitmap see subcommands.
func RunSeeCLI(args []string) error {
	if len(args) == 0 {
		return PrintSeeHelp()
	}
	joined := strings.ToLower(strings.Join(args, " "))
	if strings.HasPrefix(joined, "commit pending") {
		return runSeeCommitPending()
	}
	if strings.HasPrefix(joined, "git-ignore") || strings.HasPrefix(joined, "ignore") || strings.HasPrefix(joined, "ig issues") {
		return runSeeIgnoreIssues(args[1:])
	}
	if strings.HasPrefix(joined, "errors ssh") || joined == "ses" {
		return runSeeErrorsSSH(args)
	}
	if strings.HasPrefix(joined, "errors") {
		return cmderrors.RunErrorsCLI(args[1:])
	}
	if strings.HasPrefix(joined, "history") {
		return runSeeHistory(args[1:])
	}
	return apperror.NewSimple("unknown see command: "+joined, "E1020")
}

func runSeeCommitPending() error {
	dirty := cmdcpar.CollectDirtyRepositories()
	if len(dirty) == 0 {
		fmt.Printf("%s✓ No repositories with pending commits.%s\n",
			constants.ColorGreen, constants.ColorReset)
		return nil
	}
	fmt.Printf("\n%s  === PENDING COMMITS (%d repositories) ===%s\n",
		constants.ColorYellow, len(dirty), constants.ColorReset)
	for _, d := range dirty {
		fmt.Printf("  • %-30s | modified: %d | untracked: %d | staged: %d\n",
			d.RepoName, d.Diagnosis.ModifiedCount, d.Diagnosis.UntrackedCount, d.Diagnosis.StagedCount)
	}
	fmt.Println()
	return nil
}

func runSeeIgnoreIssues(args []string) error {
	return cmdignore.RunIgnoreCLI(append([]string{"scan"}, args...))
}

func runSeeHistory(args []string) error {
	if HistoryRunnerFn != nil {
		return HistoryRunnerFn(args)
	}
	return exec.Command("gitmap", append([]string{"history"}, args...)...).Run()
}

// RunSeeErrorsSSH executes see errors across local host and remote fleet nodes.
func RunSeeErrorsSSH(args []string) error {
	return cmdssh.RunFleetPASCommand("see errors", "gitmap errors --limit 10", func() error {
		return cmderrors.RunErrorsCLI(args)
	})
}

func runSeeErrorsSSH(args []string) error {
	subArgs := []string{}
	if len(args) > 2 {
		subArgs = args[2:]
	}
	return RunSeeErrorsSSH(subArgs)
}

// RunRepoManageUI opens the repository management web UI.
func RunRepoManageUI() error {
	fmt.Println("Launching GitMap Repository Management UI...")
	cmd := exec.Command("gitmap", "ui")
	return cmd.Start()
}

// PrintSeeHelp displays help for gitmap see suite.
func PrintSeeHelp() error {
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
