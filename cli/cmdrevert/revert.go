package cmdrevert

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

// RunRevert handles the "revert" command.
func RunRevert(args []string) error {
	checkHelp(constants.CmdRevert, args)
	isHandled, err := handleRevertTxnFlags(args)
	if err != nil {
		return err
	}
	if isHandled {
		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple(constants.MsgRevertUsage, constants.ErrCodeE9000)
	}

	version := release.NormalizeVersion(args[0])
	if err := validateRevertVersion(version); err != nil {
		return err
	}
	if err := checkoutRevertTag(version); err != nil {
		return err
	}

	return launchRevertHandoff()
}

// validateRevertVersion ensures the tag exists locally.
func validateRevertVersion(version string) error {
	if release.TagExistsLocally(version) {
		return nil
	}

	msg := fmt.Sprintf(constants.ErrTagNotFound, version)

	return apperror.NewSimple(msg, constants.ErrCodeE9000)
}

// checkoutRevertTag checks out the tag in the repo directory.
func checkoutRevertTag(version string) error {
	repoPath := getRepoPath()
	if len(repoPath) == 0 {
		return newRepoPathError(version)
	}

	fmt.Printf(constants.MsgRevertCheckout, version)

	return executeGitCheckout(repoPath, version)
}

func newRepoPathError(version string) *apperror.AppError {
	return apperror.NewWithDetails(
		constants.OpRevertCheckout,
		constants.ErrCodeE2012,
		constants.ErrNoRepoPath,
		constants.DomainCmdRevert,
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"version": version},
	)
}

func executeGitCheckout(repoPath, version string) error {
	cmd := exec.Command(constants.BinGit, constants.GitFlagDir, repoPath,
		constants.GitCmdCheckout, version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// launchRevertHandoff creates a handoff copy and runs the revert-runner.
func launchRevertHandoff() error {
	selfPath, err := os.Executable()
	if err != nil {
		return err
	}

	copyPath := cmdupdate.CreateHandoffCopy(selfPath)
	fmt.Printf(constants.MsgUpdateActive, selfPath, copyPath)

	return launchRevertRunner(copyPath)
}

// launchRevertRunner runs the handoff binary with revert-runner command.
func launchRevertRunner(copyPath string) error {
	copyArgs := []string{constants.CmdRevertRunner}
	if cmdupdate.HasFlag(constants.FlagVerbose) {
		copyArgs = append(copyArgs, constants.FlagArgVerbose)
	}

	cmd := exec.Command(copyPath, copyArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		cmdupdate.HandleHandoffError(err)

		return err
	}

	return nil
}
