package crypto

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// ConnectWithPassword establishes an SSH connection using a password.
func ConnectWithPassword(ip, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
}

// ConnectWithKey establishes an SSH connection using a private key file.
func ConnectWithKey(ip, user, keyPath string) (*ssh.Client, error) {
	signer, err := parseKeyFile(keyPath)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
}

func parseKeyFile(keyPath string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	return signer, nil
}

// RunCommand executes a command over the provided SSH client.
func RunCommand(client *ssh.Client, cmd, shellType string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	wrappedCmd := wrapCommandForShell(cmd, shellType)
	out, err := session.CombinedOutput(wrappedCmd)
	if err != nil {
		return string(out), fmt.Errorf("command execution failed: %w (output: %s)", err, string(out))
	}

	return string(out), nil
}

func wrapCommandForShell(cmd, shellType string) string {
	if shellType == "cmd" {
		return fmt.Sprintf("cmd.exe /c \"%s\"", cmd)
	}
	isPowerShell := shellType == "ps" || shellType == "pwsh" || shellType == "powershell"
	if isPowerShell {
		return fmt.Sprintf("powershell -NoProfile -Command \"%s\"", cmd)
	}
	return cmd
}
