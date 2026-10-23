package cmdcommit

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// parseRewriteFlags extracts target SHA and push flag from args.
func parseRewriteFlags(args []string) (string, bool) {
	isPush := true
	targetSha := ""

	for _, arg := range args {
		targetSha, isPush = evaluateRewriteFlag(arg, targetSha, isPush)
	}

	return targetSha, isPush
}

func evaluateRewriteFlag(arg, targetSha string, isPush bool) (string, bool) {
	if arg == "--no-push" || arg == "--local" || arg == "-n" {
		return targetSha, false
	}

	if strings.HasPrefix(arg, "-") {
		return targetSha, isPush
	}

	if targetSha == "" {
		return arg, isPush
	}

	return targetSha, isPush
}

// RunRmGit removes a commit by its SHA prefix and synchronizes remote tracking.
func RunRmGit(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdRmGit, []string{"--help"})

		return nil
	}

	targetSha, isPush := parseRewriteFlags(args)
	if len(targetSha) < 4 {
		return apperror.NewSimple("Usage: gitmap rm-git <sha-fragment> [--no-push]", "E9000")
	}

	return executeRmGit(targetSha, isPush)
}

func runRmGit(args []string) error {
	return RunRmGit(args)
}

func executeRmGit(targetSha string, isPush bool) error {
	fmt.Println()
	fullSha, errResolve := resolveSHAFragment(targetSha)
	if errResolve != nil {
		return errResolve
	}

	printPaddedInfo("Resolved SHA: %s", fullSha)
	printPaddedWarning("This will rewrite history. Use with caution.")

	return finishRmGit(fullSha, isPush)
}

func finishRmGit(fullSha string, isPush bool) error {
	if errDrop := executeDropCommit(fullSha); errDrop != nil {
		return errDrop
	}

	printPaddedSuccess("Commit %s removed successfully.", fullSha[:8])

	if isPush {
		syncRemoteAfterRewrite()
	}

	printCommandVersionFooter()

	return nil
}

// RunGitReset resets the current branch to a target SHA and synchronizes remote.
func RunGitReset(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdGitReset, []string{"--help"})

		return nil
	}

	targetSha, isPush := parseRewriteFlags(args)
	if len(targetSha) < 4 {
		return apperror.NewSimple("Usage: gitmap git-reset <target-sha> [--no-push]", "E9000")
	}

	return executeGitReset(targetSha, isPush)
}

func executeGitReset(targetSha string, isPush bool) error {
	fmt.Println()
	fullSha, errResolve := resolveSHAFragment(targetSha)
	if errResolve != nil {
		return errResolve
	}

	subject := getCommitSubject(fullSha)
	preSha, _ := execGitOutputCP("rev-parse", "HEAD")

	printPaddedInfo("Resolved target SHA: %s", fullSha)
	printPaddedWarning("Resetting branch history to %s (%s).", fullSha[:8], subject)

	return finishGitReset(fullSha, subject, preSha, isPush)
}

func finishGitReset(fullSha, subject, preSha string, isPush bool) error {
	if errReset := execGitPadded("reset", "--hard", fullSha); errReset != nil {
		printPaddedError("Failed to reset branch: %v", errReset)

		return apperror.WrapSimple(errReset, "reset branch")
	}

	printPaddedSuccess("Branch reset to %s (%s).", fullSha[:8], subject)

	if isPush {
		syncRemoteAfterRewrite()
	}

	printUndoResetHint(preSha)
	printCommandVersionFooter()

	return nil
}

func printUndoResetHint(preSha string) {
	if len(preSha) >= 8 {
		printPaddedInfo("To undo this reset locally: git reset --hard %s", strings.TrimSpace(preSha)[:8])
	}
}

func executeDropCommit(fullSha string) error {
	headSha, errHead := execGitOutputCP("rev-parse", "HEAD")
	isHead := errHead == nil && strings.HasPrefix(strings.TrimSpace(headSha), fullSha)

	if isHead {
		return execGitPadded("reset", "--hard", "HEAD~1")
	}

	return rebaseDropCommit(fullSha)
}

func rebaseDropCommit(fullSha string) error {
	if errRebase := execGitPadded("rebase", "--onto", fullSha+"^", fullSha, "HEAD"); errRebase != nil {
		execGitOutputCP("rebase", "--abort")
		printPaddedError("Failed to rebase commit %s: %v", fullSha[:8], errRebase)

		return apperror.WrapSimple(errRebase, "rebase commit")
	}

	return nil
}

func syncRemoteAfterRewrite() {
	branch, errBranch := getCurrentBranchName()
	if errBranch != nil || !hasRemoteTracking(branch) {
		return
	}

	printPaddedInfo("Synchronizing remote branch 'origin/%s' (force-pushing with lease)...", branch)
	if errPush := forcePushRemote(branch); errPush != nil {
		printPaddedWarning("Remote push with lease failed: %v", errPush)
		printPaddedWarning("Run 'git push --force origin %s' if you wish to overwrite remote history.", branch)

		return
	}

	printPaddedSuccess("Remote 'origin/%s' synchronized successfully.", branch)
}

func getCurrentBranchName() (string, error) {
	out, err := execGitOutputCP("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func hasRemoteTracking(branch string) bool {
	remote, err := execGitOutputCP("config", "--get", fmt.Sprintf("branch.%s.remote", branch))
	if err != nil || strings.TrimSpace(remote) == "" {
		return false
	}

	return true
}

func forcePushRemote(branch string) error {
	return execGitPadded("push", "--force-with-lease", "origin", branch)
}

func getCommitSubject(sha string) string {
	out, err := execGitOutputCP("log", "-1", "--format=%s", sha)
	if err != nil {
		return "commit"
	}

	return strings.TrimSpace(out)
}

func resolveSHAFragment(shaFragment string) (string, error) {
	fullSha, err := execGitOutputCP("rev-parse", shaFragment)
	if err == nil && fullSha != "" {
		return strings.TrimSpace(fullSha), nil
	}

	return searchSHAFragment(shaFragment)
}

func searchSHAFragment(shaFragment string) (string, error) {
	fullSha, err := execGitOutputCP("log", "--all", "--format=%H", "--grep="+shaFragment)
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve SHA fragment: "+shaFragment)
	}

	if fullSha == "" {
		return "", apperror.NewSimple("Could not resolve SHA fragment: "+shaFragment, "E9000")
	}

	firstMatch := strings.Split(strings.TrimSpace(fullSha), "\n")[0]

	return strings.TrimSpace(firstMatch), nil
}
