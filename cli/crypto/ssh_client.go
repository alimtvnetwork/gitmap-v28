package crypto

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func resolveTCPAddress(ip string) string {
	if strings.Contains(ip, ":") {
		return ip
	}
	return fmt.Sprintf("%s:22", ip)
}

// ConnectWithPassword establishes an SSH connection using a password.
func ConnectWithPassword(ip, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         4 * time.Second,
	}

	return ssh.Dial("tcp", resolveTCPAddress(ip), config)
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
		Timeout:         4 * time.Second,
	}

	return ssh.Dial("tcp", resolveTCPAddress(ip), config)
}

// ConnectWithFallback attempts connection with key first, falling back to password.
func ConnectWithFallback(ip, user, keyPath, password string) (*ssh.Client, error) {
	if client, isKeyOk := tryConnectWithKeyCandidate(ip, user, keyPath); isKeyOk {
		return client, nil
	}
	if password != "" {
		return ConnectWithPassword(ip, user, password)
	}
	return nil, fmt.Errorf("no valid credentials provided for %s@%s", user, ip)
}

func tryConnectWithKeyCandidate(ip, user, keyPath string) (*ssh.Client, bool) {
	if keyPath == "" {
		return nil, false
	}
	client, err := ConnectWithKey(ip, user, keyPath)
	return client, err == nil
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

// RunCommandWithInput executes a command over SSH piping input to standard input.
func RunCommandWithInput(client *ssh.Client, cmd, shellType string, in io.Reader) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	if in != nil {
		session.Stdin = in
	}
	wrappedCmd := wrapCommandForShell(cmd, shellType)
	out, err := session.CombinedOutput(wrappedCmd)
	if err != nil {
		return string(out), fmt.Errorf("command execution failed: %w (output: %s)", err, string(out))
	}
	return string(out), nil
}

func wrapCommandForShell(cmd, shellType string) string {
	trimmed := strings.TrimSpace(cmd)
	switch {
	case shellType == "cmd":
		return wrapCmdCommand(cmd, trimmed)
	case shellType == "ps" || shellType == "pwsh" || shellType == "powershell":
		return wrapPowerShellCommand(cmd, trimmed)
	case shellType == "bash":
		return wrapBashCommand(cmd, trimmed)
	case shellType == "sh":
		return wrapShCommand(cmd, trimmed)
	default:
		return cmd
	}
}

func wrapCmdCommand(cmd, trimmed string) string {
	if strings.HasPrefix(trimmed, "cmd.exe") || strings.HasPrefix(trimmed, "cmd ") {
		return cmd
	}
	return fmt.Sprintf("cmd.exe /c \"%s\"", cmd)
}

func wrapPowerShellCommand(cmd, trimmed string) string {
	if strings.HasPrefix(trimmed, "powershell") || strings.HasPrefix(trimmed, "pwsh") {
		return cmd
	}
	return fmt.Sprintf("powershell -NoProfile -Command \"%s\"", cmd)
}

func wrapBashCommand(cmd, trimmed string) string {
	if strings.HasPrefix(trimmed, "bash ") || strings.HasPrefix(trimmed, "sh ") {
		return cmd
	}
	return fmt.Sprintf("bash -c %q", cmd)
}

func wrapShCommand(cmd, trimmed string) string {
	if strings.HasPrefix(trimmed, "sh ") {
		return cmd
	}
	return fmt.Sprintf("sh -c %q", cmd)
}
