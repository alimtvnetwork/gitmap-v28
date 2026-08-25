package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

var SJAddAuthCmd = &cobra.Command{
	Use:   "add-auth [target]",
	Short: "Push local SSH public key to remote host",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJAddAuth(cmd, args, cmd.Context())
	},
}

func resolveKeyPath(keyPath string) (string, error) {
	if keyPath != "" {
		return keyPath, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".ssh", "id_rsa.pub"), nil
}

// getLocalPublicKey reads the given public key path.
func getLocalPublicKey(ctx context.Context, keyPath string, parse bool) (string, error) {
	resolvedPath, err := resolveKeyPath(keyPath)
	if err != nil {
		return "", apperror.New("getLocalPublicKey", "E_INTERNAL_ERROR", map[string]any{"err": err.Error()})
	}
	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", apperror.New("getLocalPublicKey", "E_NOT_FOUND", map[string]any{"msg": "key missing", "path": resolvedPath})
		}
		return "", apperror.New("getLocalPublicKey", "E_INTERNAL_ERROR", map[string]any{"err": err.Error(), "path": resolvedPath})
	}

	keyStr := strings.TrimSpace(string(data))
	if len(keyStr) == 0 {
		return "", apperror.New("getLocalPublicKey", "E_NOT_FOUND", map[string]any{"msg": "key file empty", "path": resolvedPath})
	}

	return keyStr, nil
}

func buildAppendScript(pubKey string) string {
	return fmt.Sprintf("mkdir -p ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys", pubKey)
}

func buildSudoAppendScript(pubKey string) string {
	return fmt.Sprintf("sudo sh -c \"mkdir -p ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys\"", pubKey)
}

// appendKeyRemote appends the key to ~/.ssh/authorized_keys on the remote target.
func appendKeyRemote(ctx context.Context, pubKey string, target SSHTarget) error {
	script := buildAppendScript(pubKey)
	err := SpawnSSH(ctx, target, []string{script})
	if err == nil {
		return nil
	}

	sudoScript := buildSudoAppendScript(pubKey)
	sudoErr := SpawnSSH(ctx, target, []string{sudoScript})
	if sudoErr != nil {
		return apperror.New("appendKeyRemote", "E_INTERNAL_ERROR", map[string]any{"err": sudoErr.Error()})
	}

	return nil
}

// runSJAddAuth handles the 'gitmap sj add-auth' command.
func runSJAddAuth(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) < 1 {
		return apperror.New("runSJAddAuth", "E_INVALID_ARGS", map[string]any{"msg": "target $ip@$user is required"})
	}

	targetStr := args[0]
	target, err := ParseSSHTarget(targetStr, "root", 22)
	if err != nil {
		return err
	}

	pubKey, err := getLocalPublicKey(ctx, "", false)
	if err != nil {
		return err
	}

	if err := appendKeyRemote(ctx, pubKey, *target); err != nil {
		return err
	}

	fmt.Printf("Added auth to %s\n", target.String())
	return nil
}

