package cmdssh

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeSSHAuthFailure_Type51(t *testing.T) {
	err := errors.New("ssh: handshake failed: ssh: unexpected message type 51 (expected 60)")
	normalized := normalizeSSHAuthFailure(err)
	if normalized == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(normalized.Error(), "authentication failed") {
		t.Errorf("expected normalized auth failure, got %v", normalized)
	}
}

func TestNormalizeSSHAuthFailure_UnableToAuth(t *testing.T) {
	err := errors.New("ssh: unable to authenticate, attempted methods [none password], no supported methods remain")
	normalized := normalizeSSHAuthFailure(err)
	if normalized == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(normalized.Error(), "authentication failed") {
		t.Errorf("expected normalized auth failure, got %v", normalized)
	}
}

func TestNormalizeSSHAuthFailure_PreservesNetworkError(t *testing.T) {
	err := errors.New("dial tcp 192.168.1.3:22: connect: connection refused")
	normalized := normalizeSSHAuthFailure(err)
	if normalized == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(normalized.Error(), "connection refused") {
		t.Errorf("expected network error to be preserved, got %v", normalized)
	}
}

func TestIsNetworkDialFailure(t *testing.T) {
	cases := []struct {
		err      error
		expected bool
	}{
		{errors.New("dial tcp 192.168.1.3:22: connect: connection refused"), true},
		{errors.New("dial tcp 192.168.1.3:22: i/o timeout"), true},
		{errors.New("dial tcp 192.168.1.3:22: network is unreachable"), true},
		{errors.New("ssh: unable to authenticate"), false},
		{nil, false},
	}

	for _, c := range cases {
		actual := isNetworkDialFailure(c.err)
		if actual != c.expected {
			t.Errorf("isNetworkDialFailure(%v) = %v, want %v", c.err, actual, c.expected)
		}
	}
}

func TestIsSSHPasswordUnsupported(t *testing.T) {
	errNoPassword := errors.New("ssh: unable to authenticate, attempted methods [none], no supported methods remain")
	if !isSSHPasswordUnsupported(errNoPassword) {
		t.Errorf("expected true when password was not in attempted methods")
	}

	errWithPassword := errors.New("ssh: unable to authenticate, attempted methods [none password], no supported methods remain")
	if isSSHPasswordUnsupported(errWithPassword) {
		t.Errorf("expected false when password was attempted")
	}
}
