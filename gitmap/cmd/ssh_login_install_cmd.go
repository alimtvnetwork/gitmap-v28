package cmd

import (
	"context"
	"fmt"
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
		target := args[0]
		version := "latest"
		if len(args) > 1 {
			version = args[1]
		}
		
		payload, err := getInstallPayload(cmd.Context(), target, version)
		if err != nil {
			return err
		}
		
		fmt.Println("Would execute:", payload)
		// Rest of the implementation would pipe this to SSH
		return nil
	},
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
