package cmdssh

import (
	"fmt"
	"net"
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
// keyboard-interactive only when password auth is unsupported by the remote server.
func dialNodeWithPassword(target *SSHTarget, pass string) (*ssh.Client, error) {
	addr := net.JoinHostPort(target.IP, strconv.Itoa(resolveHealthPort(target.Port)))

	configPass := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{ssh.Password(pass)})
	client, errPass := ssh.Dial("tcp", addr, configPass)
	if errPass == nil {
		return client, nil
	}
	if isNetworkDialFailure(errPass) {
		return nil, errPass
	}

	if !isSSHPasswordUnsupported(errPass) {
		return nil, normalizeSSHAuthFailure(errPass)
	}

	configKbd := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{buildKeyboardInteractiveAuth(pass)})
	clientKbd, errKbd := ssh.Dial("tcp", addr, configKbd)
	if errKbd != nil {
		return nil, normalizeSSHAuthFailure(errKbd)
	}

	return clientKbd, nil
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
		return fmt.Errorf("ssh: authentication failed: invalid password or remote server rejected credentials")
	}
	return err
}

func tryConnectDefaultKey(target *SSHTarget) *ssh.Client {
	for _, keyPath := range findAllUserSSHKeys() {
		res := loadPrivateKeySigner(keyPath)
		if !res.hasSigner {
			continue
		}
		if client := dialNodeWithSigner(target, res.signer); client != nil {
			return client
		}
	}
	return nil
}
