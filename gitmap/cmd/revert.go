package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// runRevert handles the "revert" command.
func runRevert(args []string) {
	checkHelp("revert", args)
	if handleRevertTxnFlags(args) {
		return
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrRevertUsage)
		os.Exit(1)
	}

	version := release.NormalizeVersion(args[0])
	validateRevertVersion(version)
	checkoutRevertTag(version)
	launchRevertHandoff()
}

// validateRevertVersion ensures the tag exists locally.
func validateRevertVersion(version string) {
	if release.TagExistsLocally(version) {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrRevertTagNotFound, version)
	os.Exit(1)
}

// checkoutRevertTag checks out the tag in the repo directory.
func checkoutRevertTag(version string) {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrNoRepoPath)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgRevertCheckout, version)
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoPath,
		constants.GitCheckout, version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRevertCheckoutFailed, err)
		os.Exit(1)
	}
}

// launchRevertHandoff creates a handoff copy and runs the revert-runner.
func launchRevertHandoff() {
	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateExecFind, err)
		os.Exit(1)
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
