package cmdssh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

var SJAddAuthCmd = &cobra.Command{
	Use:   "add-auth [target]",
	Short: "Push local SSH public key to remote host",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJAddAuth(cmd, args, cmd.Context())
	},
}

const msgMissingAuthTarget = `missing target host or alias

Usage:
  gitmap ssh join add-auth <alias|ip|user@ip>
  gitmap sj add-auth <alias|ip|user@ip>

Examples:
  gitmap sj add-auth devbox
  gitmap sj add-auth 192.168.1.14
  gitmap sj add-auth alim@192.168.1.14`

func checkExistingKeyFile(candidate string) (string, bool) {
	if _, err := os.Stat(candidate); err == nil {
		return candidate, true
	}
	return "", false
}

func findDefaultPublicKey(homeDir string) string {
	rsa := filepath.Join(homeDir, ".ssh", "id_rsa.pub")
	if path, isFound := checkExistingKeyFile(rsa); isFound {
		return path
	}

	ed := filepath.Join(homeDir, ".ssh", "id_ed25519.pub")
	if path, isFound := checkExistingKeyFile(ed); isFound {
		return path
	}

	return rsa
}

func resolveKeyPath(keyPath string) (string, error) {
	if keyPath != "" {
		return keyPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return findDefaultPublicKey(homeDir), nil
}

func readKeyFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", apperror.New("getLocalPublicKey", "E_NOT_FOUND", map[string]any{"msg": "key missing", "path": path})
	}
	if err != nil {
		return "", apperror.New("getLocalPublicKey", "E_INTERNAL_ERROR", map[string]any{"err": err.Error(), "path": path})
	}

	keyStr := strings.TrimSpace(string(data))
	if len(keyStr) == 0 {
		return "", apperror.New("getLocalPublicKey", "E_NOT_FOUND", map[string]any{"msg": "key file empty", "path": path})
	}
	return keyStr, nil
}

// getLocalPublicKey reads the given public key path.
func getLocalPublicKey(ctx context.Context, keyPath string, isParsed bool) (string, error) {
	resolvedPath, err := resolveKeyPath(keyPath)
	if err != nil {
		return "", apperror.New("getLocalPublicKey", "E_INTERNAL_ERROR", map[string]any{"err": err.Error()})
	}

	return readKeyFile(resolvedPath)
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
	if err := SpawnSSH(ctx, target, []string{script}); err == nil {
		return nil
	}

	sudoScript := buildSudoAppendScript(pubKey)
	if sudoErr := SpawnSSH(ctx, target, []string{sudoScript}); sudoErr != nil {
		return apperror.New("appendKeyRemote", "E_INTERNAL_ERROR", map[string]any{"err": sudoErr.Error()})
	}

	return nil
}

func validateAuthTarget(args []string) (string, error) {
	if len(args) == 0 {
		return "", apperror.NewValidationError(msgMissingAuthTarget)
	}

	target := strings.TrimSpace(args[0])
	if target == "" {
		return "", apperror.NewValidationError(msgMissingAuthTarget)
	}

	return target, nil
}

func resolveHostAliasTarget(ctx context.Context, target string) string {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return target
	}
	defer dbConn.Close()

	host, err := store.GetHostByAlias(ctx, target, dbConn.SQL())
	if err != nil || host.IP == "" {
		return target
	}

	return fmt.Sprintf("%s@%s", host.Username, host.IP)
}

func pushAuthToTarget(ctx context.Context, target SSHTarget) error {
	pubKey, err := getLocalPublicKey(ctx, "", false)
	if err != nil {
		return err
	}

	if err := appendKeyRemote(ctx, pubKey, target); err != nil {
		return err
	}

	fmt.Printf("Added auth to %s\n", target.String())
	return nil
}

// runSJAddAuth handles the 'gitmap sj add-auth' command.
//
//nolint:revive
func runSJAddAuth(cmd *cobra.Command, args []string, ctx context.Context) error {
	rawTarget, err := validateAuthTarget(args)
	if err != nil {
		return err
	}

	resolvedTarget := resolveHostAliasTarget(ctx, rawTarget)
	target, err := ParseSSHTarget(resolvedTarget, resolveDefaultUsername(), 22)
	if err != nil {
		return err
	}

	return pushAuthToTarget(ctx, *target)
}
