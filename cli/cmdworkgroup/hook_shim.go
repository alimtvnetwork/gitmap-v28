package cmdworkgroup

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

// ValidateSSHKey checks whether the specified private key file exists.
func ValidateSSHKey(keyPath string) error {
	if keyPath == "" {
		return apperror.NewSimple("SSH key path cannot be empty", "E9001")
	}
	if !fsutil.FileExists(keyPath) {
		return apperror.NewSimple(fmt.Sprintf("SSH key file does not exist: %s", keyPath), "E9002")
	}

	return nil
}

// SynthesizeGitSSHCommand creates an isolated GIT_SSH_COMMAND string for a key.
func SynthesizeGitSSHCommand(keyPath string) string {
	cleanPath := strings.ReplaceAll(keyPath, "\\", "/")

	return fmt.Sprintf("ssh -i \"%s\" -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new -F /dev/null", cleanPath)
}

func resolveSSHHost(host string) string {
	if strings.TrimSpace(host) == "" {
		return "github.com"
	}

	return strings.TrimSpace(host)
}

func runSSHProbe(keyPath, targetHost string) (string, error) {
	cmd := exec.Command("ssh", "-T", "-i", keyPath,
		"-o", "StrictHostKeyChecking=accept-new",
		"-o", "IdentitiesOnly=yes",
		"-o", "BatchMode=yes",
		"-o", "ConnectTimeout=10",
		"-F", "/dev/null",
		"git@"+targetHost,
	)
	out, err := cmd.CombinedOutput()
	msg := strings.TrimSpace(string(out))
	if strings.Contains(msg, "successfully authenticated") || strings.Contains(msg, "Welcome to GitLab") {
		return msg, nil
	}
	if err != nil && msg != "" {
		return msg, nil
	}

	return msg, err
}

// TestSSHConnection executes an isolated SSH handshake test using the bound key.
func TestSSHConnection(keyPath, host string) (string, error) {
	if err := ValidateSSHKey(keyPath); err != nil {
		return "", err
	}
	targetHost := resolveSSHHost(host)
	out, err := runSSHProbe(keyPath, targetHost)
	if err != nil && out == "" {
		return "", apperror.WrapSimple(err, fmt.Sprintf("SSH handshake failed against %s", targetHost))
	}

	return out, nil
}

// BuildWorkgroupEnv returns an environment slice with GIT_SSH_COMMAND injected.
func BuildWorkgroupEnv(baseEnv []string, keyPath string) []string {
	if keyPath == "" {
		return baseEnv
	}
	sshCmd := SynthesizeGitSSHCommand(keyPath)
	res := append([]string{}, baseEnv...)

	return append(res, "GIT_SSH_COMMAND="+sshCmd)
}
