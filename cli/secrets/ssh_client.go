package secrets

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolveTCPAddress(ip string) string {
	if strings.Contains(ip, ":") {
		return ip
	}
	return fmt.Sprintf("%s:22", ip)
}

// defaultSSHHandshakeTimeout bounds an SSH dial when the client config
// carries no timeout of its own.
const defaultSSHHandshakeTimeout = 10 * time.Second

// dialSSHClient dials an SSH client under a hard deadline spanning the TCP
// connect and the SSH handshake. ssh.Dial applies ClientConfig.Timeout only
// to the TCP connect; the version exchange inside the handshake reads with
// no deadline, so a peer that accepts TCP but never sends an SSH banner
// (silent firewall, hung sshd, transparent proxy) blocks the caller forever.
func dialSSHClient(addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultSSHHandshakeTimeout
	}

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		_ = conn.Close()

		return nil, err
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		_ = conn.Close()

		return nil, err
	}

	// The handshake is complete; session traffic must not inherit the
	// handshake deadline.
	_ = conn.SetDeadline(time.Time{})

	return ssh.NewClient(sshConn, chans, reqs), nil
}

// ConnectWithPassword establishes an SSH connection using a password.
func ConnectWithPassword(ip, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         4 * time.Second,
	}

	return dialSSHClient(resolveTCPAddress(ip), config)
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

	return dialSSHClient(resolveTCPAddress(ip), config)
}

// ConnectWithFallback attempts connection with key first, falling back to password.
func ConnectWithFallback(ip, user, keyPath, password string) (*ssh.Client, error) {
	if client, isKeyOk := tryConnectWithKeyCandidate(ip, user, keyPath); isKeyOk {
		return client, nil
	}
	if password != "" {
		return ConnectWithPassword(ip, user, password)
	}
	return nil, apperror.NewValidationError(fmt.Sprintf("no valid credentials provided for %s@%s", user, ip))
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
		return nil, apperror.WrapSimple(err, "read private key")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, apperror.WrapSimple(err, "parse private key")
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
		return string(out), apperror.WrapSimple(err, "command execution failed").WithContext("output", string(out))
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
		return string(out), apperror.WrapSimple(err, "command execution failed").WithContext("output", string(out))
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
