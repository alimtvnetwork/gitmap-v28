package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

// SSHLoginInstallCmd represents the gitmap ssh login-install command.
var SSHLoginInstallCmd = &cobra.Command{
	Use:   "login-install [target] [version]",
	Short: "Connects to the remote machine and executes the gitmap installation script",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHLoginInstall(cmd, args, cmd.Context())
	},
}

func runSSHLoginInstall(cmd *cobra.Command, args []string, ctx context.Context) error {
	targetStr := args[0]
	version := "latest"
	if len(args) > 1 {
		version = args[1]
	}

	payload, err := getInstallPayload(ctx, targetStr, version)
	if err != nil {
		return err
	}

	target, err := ParseSSHTarget(targetStr, "root", 22)
	if err != nil {
		return err
	}

	return executeRemoteInstall(ctx, payload, *target)
}

func getInstallPayload(ctx context.Context, target string, version string) (string, error) {
	if strings.TrimSpace(target) == "" {
		return "", apperror.New("getInstallPayload", "E_INTERNAL_ERROR", nil)
	}
	
	// Create curl to bash payload
	url := fmt.Sprintf("https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/refs/heads/main/install.sh?version=%s", version)
	payload := fmt.Sprintf("curl -fsSL %s | bash", url)
	
	return payload, nil
}

func executeRemoteInstall(ctx context.Context, payload string, target SSHTarget) error {
	cmd := exec.CommandContext(ctx, "ssh", target.String())
	cmd.Stdin = strings.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return &apperror.AppError{
			Op:    "executeRemoteInstall",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"target": target.String()},
		}
	}

	return nil
}
