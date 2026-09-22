package cmdssh

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func dialNodeWithSigner(target *SSHTarget, signer ssh.Signer) *ssh.Client {
	config := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{ssh.PublicKeys(signer)})
	addr := net.JoinHostPort(target.IP, strconv.Itoa(resolveHealthPort(target.Port)))
	client, err := ssh.Dial("tcp", addr, config)
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
	client, errPass := ssh.Dial("tcp", addr, configPass)
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
	clientKbd, errKbd := ssh.Dial("tcp", addr, configKbd)
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
		return fmt.Errorf("ssh: authentication failed: invalid password or remote server rejected credentials: %w", err)
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
