package cluster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"golang.org/x/crypto/bcrypt"
)

// ExecRestart triggers a machine restart.
func ExecRestart(ctx context.Context, node ClusterNode, forceLifecycle bool, providedPassword string) (string, string, int, error) {
	if err := checkLifecycleGuards(node, forceLifecycle, providedPassword); err != nil {
		return "", "", 1, err
	}
	isWin := runtime.GOOS == "windows"
	var cmd *exec.Cmd
	if isWin {
		cmd = exec.CommandContext(ctx, "shutdown", "/r", "/t", "0")
	} else {
		cmd = exec.CommandContext(ctx, "reboot")
	}
	return runCmd(cmd)
}

// ExecShutdown triggers a machine shutdown.
func ExecShutdown(ctx context.Context, node ClusterNode, forceLifecycle bool, providedPassword string) (string, string, int, error) {
	if err := checkLifecycleGuards(node, forceLifecycle, providedPassword); err != nil {
		return "", "", 1, err
	}
	isWin := runtime.GOOS == "windows"
	var cmd *exec.Cmd
	if isWin {
		cmd = exec.CommandContext(ctx, "shutdown", "/s", "/t", "0")
	} else {
		cmd = exec.CommandContext(ctx, "shutdown", "-h", "now")
	}
	return runCmd(cmd)
}

// ExecLogoff logs off the current user.
func ExecLogoff(ctx context.Context, node ClusterNode, forceLifecycle bool, providedPassword string) (string, string, int, error) {
	if err := checkLifecycleGuards(node, forceLifecycle, providedPassword); err != nil {
		return "", "", 1, err
	}
	isWin := runtime.GOOS == "windows"
	var cmd *exec.Cmd
	if isWin {
		cmd = exec.CommandContext(ctx, "logoff")
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", "pkill -KILL -u whoami")
	}
	return runCmd(cmd)
}

func checkLifecycleGuards(node ClusterNode, forceLifecycle bool, providedPassword string) error {
	if node.NodeRole == "server" || node.IsServer {
		return errors.New(constants.ErrClusterServerProtected)
	}
	if !forceLifecycle {
		return errors.New(constants.ErrClusterLifecycleRequiresForce)
	}
	if node.PasswordHash == "" {
		return nil
	}
	if providedPassword == "" {
		return errors.New(constants.ErrClusterPasswordRequired)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(node.PasswordHash), []byte(providedPassword)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

func runCmd(cmd *exec.Cmd) (string, string, int, error) {
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := runCmdFunc(cmd)
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}
	return outBuf.String(), errBuf.String(), exitCode, err
}

func PrintCountdown(ctx context.Context, nodes []string, action string, seconds int) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := seconds; i > 0; i-- {
		fmt.Printf(constants.MsgClusterCountdown+"\n", action, len(nodes), i)
		select {
		case <-ctx.Done():
			fmt.Println("Aborted by user.")
			return ctx.Err()
		case <-ticker.C:
		}
	}
	return nil
}
