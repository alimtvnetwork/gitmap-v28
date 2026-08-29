package cluster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ExecRestart triggers a machine restart.

func ExecRestart(ctx context.Context, node ClusterNode, forceLifecycle bool, providedPassword string) (string, string, int, error) {
	if err := checkLifecycleGuards(node, forceLifecycle, providedPassword); err != nil {
		return "", "", constants.ExitCodeError, err
	}
	return runCmd(buildRestartCmd(ctx))
}

func buildRestartCmd(ctx context.Context) *exec.Cmd {
	if runtime.GOOS == constants.PlatformWindows {
		return exec.CommandContext(ctx, constants.LifecycleCmdShutdown, constants.ArgRestart, constants.ArgTimeout, constants.ArgZero)
	}
	return exec.CommandContext(ctx, constants.LifecycleCmdReboot)
}

// ExecShutdown triggers a machine shutdown.

func ExecShutdown(ctx context.Context, node ClusterNode, forceLifecycle bool, providedPassword string) (string, string, int, error) {
	if err := checkLifecycleGuards(node, forceLifecycle, providedPassword); err != nil {
		return "", "", constants.ExitCodeError, err
	}
	return runCmd(buildShutdownCmd(ctx))
}

func buildShutdownCmd(ctx context.Context) *exec.Cmd {
	if runtime.GOOS == constants.PlatformWindows {
		return exec.CommandContext(ctx, constants.LifecycleCmdShutdown, constants.ArgShutdownWin, constants.ArgTimeout, constants.ArgZero)
	}
	return exec.CommandContext(ctx, constants.LifecycleCmdShutdown, constants.ArgHalt, constants.ArgNow)
}

// ExecLogoff logs off the current user.

func ExecLogoff(ctx context.Context, node ClusterNode, forceLifecycle bool, providedPassword string) (string, string, int, error) {
	if err := checkLifecycleGuards(node, forceLifecycle, providedPassword); err != nil {
		return "", "", constants.ExitCodeError, err
	}
	return runCmd(buildLogoffCmd(ctx))
}

func buildLogoffCmd(ctx context.Context) *exec.Cmd {
	if runtime.GOOS == constants.PlatformWindows {
		return exec.CommandContext(ctx, constants.LifecycleCmdLogoff)
	}
	return exec.CommandContext(ctx, constants.UnixShell, constants.UnixShellArg, constants.LifecycleCmdUnixLogoffArgs)
}

func checkLifecycleGuards(node ClusterNode, forceLifecycle bool, providedPassword string) error {
	if node.NodeRole == constants.NodeRoleServer || node.IsServer {
		return errors.New(constants.ErrClusterServerProtected)
	}
	if !forceLifecycle {
		return errors.New(constants.ErrClusterLifecycleRequiresForce)
	}
	return checkPasswordAuth(node.PasswordHash, providedPassword)
}

func checkPasswordAuth(hash, password string) error {
	if hash == "" {
		return nil
	}
	if password == "" {
		return errors.New(constants.ErrClusterPasswordRequired)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New(constants.ErrClusterInvalidPassword)
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
	if exitErr, ok := err.(*exec.ExitError); ok {
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

func countdownTick(ctx context.Context, tickChan <-chan time.Time, action string, count, remaining int) error {
	fmt.Printf(constants.MsgClusterCountdown+"\n", action, count, remaining)
	select {
	case <-ctx.Done():
		fmt.Println(constants.MsgClusterAbortedByUser)
		return ctx.Err()
	case <-tickChan:
		return nil
	}
}
