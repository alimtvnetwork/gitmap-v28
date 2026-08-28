package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// runRevert handles the "revert" command.
func runRevert(args []string) error {
	checkHelp("revert", args)
	if handleRevertTxnFlags(args) {
		return nil
	}
	if len(args) == 0 {
		return apperror.New(constants.ErrRevertUsage, "E9000", nil)
	}

	version := release.NormalizeVersion(args[0])
	validateRevertVersion(version)
	checkoutRevertTag(version)
	launchRevertHandoff()
	return nil
}

// validateRevertVersion ensures the tag exists locally.
func validateRevertVersion(version string) {
	if release.TagExistsLocally(version) {
		return
	}

	return apperror.New(constants.ErrRevertTagNotFound, "E9000", nil)
}

// checkoutRevertTag checks out the tag in the repo directory.
func checkoutRevertTag(version string) {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrNoRepoPath)
		return apperror.New("fatal error", "E9000", nil)
	}

	fmt.Printf(constants.MsgRevertCheckout, version)
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoPath,
		constants.GitCheckout, version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, constants.ErrRevertCheckoutFailed, nil)
	}
}

// launchRevertHandoff creates a handoff copy and runs the revert-runner.
func launchRevertHandoff() {
	selfPath, err := os.Executable()
	if err != nil {
		return apperror.Wrap(err, constants.ErrUpdateExecFind, nil)
	}

	copyPath := createHandoffCopy(selfPath)
	fmt.Printf(constants.MsgUpdateActive, selfPath, copyPath)
	launchRevertRunner(copyPath)
}

// launchRevertRunner runs the handoff binary with revert-runner command.
func launchRevertRunner(copyPath string) {
	copyArgs := []string{constants.CmdRevertRunner}
	if hasFlag(constants.FlagVerbose) {
		copyArgs = append(copyArgs, constants.FlagVerbose)
	}

	cmd := exec.Command(copyPath, copyArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		handleHandoffError(err)
	}
}
