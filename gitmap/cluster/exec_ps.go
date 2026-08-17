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
	if isWin == true {
		cmdPath, errLook := exec.LookPath(pwshCmd)
		if errLook != nil {
			cmdPath, errLook = exec.LookPath(powershellCmd)
			if errLook != nil {
				return "", "", 1, errLook
			}
		}
		cmd = exec.CommandContext(ctx, cmdPath, nonIntFlag, cmdFlag, command)
	} else {
		cmdPath, errLook := exec.LookPath(pwshCmd)
		if errLook != nil {
			return "", "pwsh not found, skipping", 0, nil
		}
		cmd = exec.CommandContext(ctx, cmdPath, cmdFlag, command)
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	
	stdout = outBuf.String()
	stderr = errBuf.String()
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	} else {
		exitCode = 0
	}

	return stdout, stderr, exitCode, err
}
