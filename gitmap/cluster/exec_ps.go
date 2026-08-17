package cluster

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
)

const (
	pwshCmd       = "pwsh"
	powershellCmd = "powershell"
	cmdFlag       = "-Command"
	nonIntFlag    = "-NonInteractive"
	winOS         = "windows"
)

// ExecPS executes a PowerShell command on the specified node.
func ExecPS(ctx context.Context, node ClusterNode, command string) (stdout, stderr string, exitCode int, err error) {
	var cmd *exec.Cmd

	isWin := runtime.GOOS == winOS

	cmdPath, errLook := "", error(nil)
	if isWin {
		cmdPath, errLook = exec.LookPath(pwshCmd)
	}
	if isWin && errLook != nil {
		cmdPath, errLook = exec.LookPath(powershellCmd)
	}
	if isWin && errLook != nil {
		return "", "", 1, errLook
	}

	if !isWin {
		cmdPath, errLook = exec.LookPath(pwshCmd)
	}
	if !isWin && errLook != nil {
		return "", "pwsh not found, skipping", 0, nil
	}

	if isWin {
		cmd = exec.CommandContext(ctx, cmdPath, nonIntFlag, cmdFlag, command)
	}
	if !isWin {
		cmd = exec.CommandContext(ctx, cmdPath, cmdFlag, command)
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	stdout = outBuf.String()
	stderr = errBuf.String()

	exitCode = 0
	if err != nil {
		exitCode = 1
	}
	exitErr, isExitErr := err.(*exec.ExitError)
	if isExitErr {
		exitCode = exitErr.ExitCode()
	}

	return stdout, stderr, exitCode, err
}
