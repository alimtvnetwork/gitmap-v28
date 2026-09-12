package cluster

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ExecRestart triggers a machine restart.

func ExecRestart(params LifecycleExecParams) (string, string, int, *apperror.AppError) {
	if err := checkLifecycleGuards(params.Node, params.IsForceLifecycle, params.ProvidedPassword); err != nil {
		return "", "", constants.ExitCodeError, err
	}

	stdout, stderr, code, rawErr := runCmd(buildRestartCmd(params.Context))
	if rawErr != nil {
		return stdout, stderr, code, apperror.WrapSimple(rawErr, "ExecRestart")
	}

	return stdout, stderr, code, nil
}

func buildRestartCmd(ctx context.Context) *exec.Cmd {
	if runtime.GOOS == constants.PlatformWindows {
		return exec.CommandContext(ctx, constants.LifecycleCmdShutdown, constants.ArgRestart, constants.ArgTimeout, constants.ArgZero)
	}

	return exec.CommandContext(ctx, constants.LifecycleCmdReboot)
}

// ExecShutdown triggers a machine shutdown.

func ExecShutdown(params LifecycleExecParams) (string, string, int, *apperror.AppError) {
	if err := checkLifecycleGuards(params.Node, params.IsForceLifecycle, params.ProvidedPassword); err != nil {
		return "", "", constants.ExitCodeError, err
	}

	stdout, stderr, code, rawErr := runCmd(buildShutdownCmd(params.Context))
	if rawErr != nil {
		return stdout, stderr, code, apperror.WrapSimple(rawErr, "ExecShutdown")
	}

	return stdout, stderr, code, nil
}

func buildShutdownCmd(ctx context.Context) *exec.Cmd {
	if runtime.GOOS == constants.PlatformWindows {
		return exec.CommandContext(ctx, constants.LifecycleCmdShutdown, constants.ArgShutdownWin, constants.ArgTimeout, constants.ArgZero)
	}

	return exec.CommandContext(ctx, constants.LifecycleCmdShutdown, constants.ArgHalt, constants.ArgNow)
}

// ExecLogoff logs off the current user.

func ExecLogoff(params LifecycleExecParams) (string, string, int, *apperror.AppError) {
	if err := checkLifecycleGuards(params.Node, params.IsForceLifecycle, params.ProvidedPassword); err != nil {
		return "", "", constants.ExitCodeError, err
	}

	stdout, stderr, code, rawErr := runCmd(buildLogoffCmd(params.Context))
	if rawErr != nil {
		return stdout, stderr, code, apperror.WrapSimple(rawErr, "ExecLogoff")
	}

	return stdout, stderr, code, nil
}

func buildLogoffCmd(ctx context.Context) *exec.Cmd {
	if runtime.GOOS == constants.PlatformWindows {
		return exec.CommandContext(ctx, constants.LifecycleCmdLogoff)
	}

	return exec.CommandContext(ctx, constants.UnixShell, constants.UnixShellArg, constants.LifecycleCmdUnixLogoffArgs)
}

func checkLifecycleGuards(node ClusterNode, isForceLifecycle bool, providedPassword string) *apperror.AppError {
	if node.NodeRole == constants.NodeRoleServer || node.IsServer {
		return apperror.NewSimple(constants.ErrClusterServerProtected, "E8001")
	}

	if !isForceLifecycle {
		return apperror.NewSimple(constants.ErrClusterLifecycleRequiresForce, "E8002")
	}

	return checkPasswordAuth(node.PasswordHash, providedPassword)
}

func checkPasswordAuth(hash, password string) *apperror.AppError {
	if hash == "" {
		return nil
	}

	if password == "" {
		return apperror.NewSimple(constants.ErrClusterPasswordRequired, "E8003")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return apperror.NewSimple(constants.ErrClusterInvalidPassword, "E8004")
	}

	return nil
}

func runCmd(cmd *exec.Cmd) (string, string, int, error) {
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := runCmdFunc(cmd)
	exitCode := extractExitCode(err)

	return outBuf.String(), errBuf.String(), exitCode, err
}

func extractExitCode(err error) int {
	if err == nil {
		return constants.ExitCodeSuccess
	}

	exitErr, isExitErr := err.(*exec.ExitError)
	if isExitErr {
		return exitErr.ExitCode()
	}

	return constants.ExitCodeError
}

func PrintCountdown(ctx context.Context, nodes []string, action string, seconds int) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for i := seconds; i > 0; i-- {
		if err := countdownTick(ctx, ticker.C, action, len(nodes), i); err != nil {
			return err
		}
	}

	return nil
}

func countdownTick(
	ctx context.Context,
	tickChan <-chan time.Time,
	action string,
	count,
	remaining int,
) error {
	fmt.Printf(constants.MsgClusterCountdown+"\n", action, count, remaining)
	select {
	case <-ctx.Done():
		fmt.Println(constants.MsgClusterAbortedByUser)

		return ctx.Err()
	case <-tickChan:
		return nil
	}
}
