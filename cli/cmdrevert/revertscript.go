package cmdrevert

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
)

// RunRevertRunner is the hidden command that performs the revert build.
func RunRevertRunner() error {
	repoPath := getRepoPath()
	if len(repoPath) == 0 {
		return newRunnerRepoPathError()
	}

	cmdupdate.InitRunnerVerbose()
	fmt.Print(constants.MsgRevertStarting)
	if err := executeRevert(repoPath); err != nil {
		return err
	}

	cmdupdate.ScheduleDeployedCleanupHandoff()

	return nil
}

func newRunnerRepoPathError() *apperror.AppError {
	return apperror.NewWithDetails(
		constants.OpRevertRunner,
		constants.ErrCodeE2013,
		constants.ErrNoRepoPath,
		constants.DomainCmdRevertScript,
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		nil,
	)
}

// executeRevert writes a temp PS1 script and runs it.
func executeRevert(repoPath string) error {
	scriptPath, err := writeRevertScript(repoPath)
	if err != nil {
		return err
	}

	defer os.Remove(scriptPath)
	logRevertExec(scriptPath)

	return runRevertPS(scriptPath)
}

func logRevertExec(scriptPath string) {
	log := verbose.Get()
	if log == nil {
		return
	}

	log.Log(constants.RevertScriptLogExec, scriptPath)
}

// writeRevertScript creates a temporary PowerShell script for revert build.
func writeRevertScript(repoPath string) (string, error) {
	runPS1 := filepath.Join(repoPath, constants.ScriptRunPS1)
	script := buildRevertScript(repoPath, runPS1)

	return cmdupdate.WriteScriptToTemp(script)
}

// buildRevertScript generates the PowerShell script content for revert.
func buildRevertScript(repoPath, runPS1 string) string {
	return fmt.Sprintf(constants.RevertPSHeader, repoPath) +
		fmt.Sprintf(constants.RevertPSBuild, runPS1) +
		constants.RevertPSPostActions
}

// runRevertPS executes the PowerShell script with output piped to terminal.
func runRevertPS(scriptPath string) error {
	cmd := exec.Command(constants.BinPowershell, constants.PSArgExecPolicy, constants.PSValBypass,
		constants.PSArgNoProfile, constants.PSArgNoLogo, constants.PSArgFile, scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()

	logRevertResult(err)
	if err != nil {
		return apperror.WrapSimple(err, constants.OpRunRevertPS)
	}

	fmt.Print(constants.MsgRevertDone)

	return nil
}

// logRevertResult logs the revert script exit status if verbose is active.
func logRevertResult(err error) {
	log := verbose.Get()
	if log == nil {
		return
	}

	log.Log(constants.RevertScriptLogExit, err)
}
