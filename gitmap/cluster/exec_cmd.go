package cluster

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

//nolint:unused
type execRunnerFunc func(cmd *exec.Cmd) error

//nolint:unused
type lookPathFunc func(file string) (string, error)

var (
	runCmdFunc      = defaultExecRunner
	lookPathFuncVar = defaultLookPath
)

func defaultExecRunner(cmd *exec.Cmd) error {
	return cmd.Run()
}

func defaultLookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// ExecCmd runs a command on the target node (currently executing locally via the OS shell).

func ExecCmd(
	ctx context.Context,
	node ClusterNode,
	command string,
) (stdout, stderr string, exitCode int, err error) {
	var cmd *exec.Cmd
	isWindows := runtime.GOOS == constants.WindowsOS
	if isWindows {
		cmd = exec.CommandContext(ctx, constants.WindowsShell, constants.WindowsShellArg, command)
	} else {
		cmd = exec.CommandContext(ctx, constants.UnixShell, constants.UnixShellArg, command)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = runCmdFunc(cmd)

	exitCode = constants.ExitCodeSuccess
	if err != nil {
		exitCode = constants.ExitCodeError
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	return stdoutBuf.String(), stderrBuf.String(), exitCode, err
}
