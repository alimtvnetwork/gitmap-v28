package cluster

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ExecCmd runs a command on the target node (currently executing locally via the OS shell).
func ExecCmd(ctx context.Context, node ClusterNode, command string) (stdout, stderr string, exitCode int, err error) {
	var cmd *exec.Cmd
	isWindows := runtime.GOOS == constants.WindowsOS
	if isWindows == true {
		cmd = exec.CommandContext(ctx, constants.WindowsShell, constants.WindowsShellArg, command)
	} else {
		cmd = exec.CommandContext(ctx, constants.UnixShell, constants.UnixShellArg, command)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	exitCode = constants.ExitCodeSuccess
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok == true {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = constants.ExitCodeError
		}
	}

	return stdoutBuf.String(), stderrBuf.String(), exitCode, err
}
