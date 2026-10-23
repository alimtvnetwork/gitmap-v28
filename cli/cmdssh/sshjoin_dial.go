package cmdssh

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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

func dialNodeWithSigner(target *SSHTarget, signer ssh.Signer) *ssh.Client {
	config := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{ssh.PublicKeys(signer)})
	addr := net.JoinHostPort(target.IP, strconv.Itoa(resolveHealthPort(target.Port)))
	client, err := dialSSHClient(addr, config)
	if err != nil {
		return nil
	}
	return client
}

func buildKeyboardInteractiveAuth(pass string) ssh.AuthMethod {
	return ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
		answers := make([]string, len(questions))
		for i := range answers {
			answers[i] = pass
		}
		return answers, nil
	})
}

// dialNodeWithPassword attempts standard password auth first, falling back to
// keyboard-interactive (PAM) if password auth is rejected.
func dialNodeWithPassword(target *SSHTarget, pass string) (*ssh.Client, error) {
	addr := net.JoinHostPort(target.IP, strconv.Itoa(resolveHealthPort(target.Port)))
	trace := GetActiveSSHTrace()

	configPass := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{ssh.Password(pass)})
	client, errPass := dialSSHClient(addr, configPass)
	if errPass == nil {
		trace.AddStep("Password Auth", fmt.Sprintf("authenticated '%s@%s'", target.Username, addr), "SUCCESS", nil)

		return client, nil
	}
	if isNetworkDialFailure(errPass) {
		trace.AddStep("Password Auth", "network dial failed: "+errPass.Error(), "FAILED", errPass)

		return nil, errPass
	}

	trace.AddStep("Password Auth", fmt.Sprintf("password rejected for '%s@%s'", target.Username, addr), "FAILED", errPass)

	return dialKeyboardInteractiveFallback(target, addr, pass, errPass, trace)
}

func dialKeyboardInteractiveFallback(target *SSHTarget, addr, pass string, errPass error, trace *SSHExecutionTrace) (*ssh.Client, error) {
	configKbd := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{buildKeyboardInteractiveAuth(pass)})
	clientKbd, errKbd := dialSSHClient(addr, configKbd)
	if errKbd == nil {
		trace.AddStep("Keyboard-Interactive Auth", fmt.Sprintf("authenticated '%s@%s' via PAM", target.Username, addr), "SUCCESS", nil)

		return clientKbd, nil
	}

	trace.AddStep("Keyboard-Interactive Auth", fmt.Sprintf("PAM rejected for '%s@%s'", target.Username, addr), "FAILED", errKbd)

	return nil, normalizeSSHAuthFailure(errPass)
}

func isNetworkDialFailure(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "refused") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "no route") ||
		strings.Contains(msg, "unreachable")
}

func isSSHPasswordUnsupported(err error) bool {
	if err == nil {
		return false
	}
	// If server only advertises keyboard-interactive, 'password' was never tried.
	return !strings.Contains(err.Error(), "password")
}

func normalizeSSHAuthFailure(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.Contains(msg, "unexpected message type 51") ||
		strings.Contains(msg, "unable to authenticate") ||
		strings.Contains(msg, "handshake failed") {
		return apperror.WrapSimple(err, "ssh: authentication failed: invalid password or remote server rejected credentials")
	}
	return err
}

func tryConnectDefaultKey(target *SSHTarget) *ssh.Client {
	trace := GetActiveSSHTrace()
	keys := findAllUserSSHKeys()
	for _, keyPath := range keys {
		res := loadPrivateKeySigner(keyPath)
		if !res.hasSigner {
			continue
		}
		if client := dialNodeWithSigner(target, res.signer); client != nil {
			trace.AddStep("Public Key Auth", fmt.Sprintf("authenticated using %s", filepath.Base(keyPath)), "SUCCESS", nil)

			return client
		}
	}
	if len(keys) > 0 {
		trace.AddStep("Public Key Auth", fmt.Sprintf("checked %d default user keys, none accepted", len(keys)), "SKIPPED", nil)
	}

	return nil
}
