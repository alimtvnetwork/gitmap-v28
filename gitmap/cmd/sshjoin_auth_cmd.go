package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

// getLocalPublicKey reads the given public key path.
func getLocalPublicKey(ctx context.Context, keyPath string, parse bool) (string, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", apperror.New("getLocalPublicKey", "E_NOT_FOUND", map[string]any{"msg": "key missing"})
		}
		return "", apperror.New("getLocalPublicKey", "E_INTERNAL_ERROR", map[string]any{"err": err.Error()})
	}
	
	keyStr := strings.TrimSpace(string(data))
	if parse {
		// Just a placeholder in case parse logic is needed later
	}
	return keyStr, nil
}

// appendKeyRemote appends the key to ~/.ssh/authorized_keys on the remote target.
func appendKeyRemote(ctx context.Context, pubKey string, target SSHTarget) error {
	script := fmt.Sprintf("mkdir -p ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys", pubKey)
	args := []string{"-c", script}
	if err := SpawnSSH(ctx, target, args); err != nil {
		return apperror.New("appendKeyRemote", "E_INTERNAL_ERROR", map[string]any{"err": err.Error()})
	}
	return nil
}

// runSJAddAuth handles the 'gitmap sj add-auth' command.
func runSJAddAuth(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) < 1 {
		return apperror.New("runSJAddAuth", "E_INVALID_ARGS", map[string]any{"msg": "target $ip@$user is required"})
	}

	targetStr := args[0]
	// Handle $ip@$user targeting
	fmt.Printf("Adding auth to %s\n", targetStr)
	return nil
}

