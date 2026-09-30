package cmdssh

import (
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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

func TestLookupVaultPassword(t *testing.T) {
	withMockSSHDB(t, func(mockDB *store.DB) {
		ctx := t.Context()
		pass := "mySecretPass123"
		encPass, _ := EncryptSSHPassword(pass)
		host := store.SSHHost{
			IP:                "192.168.1.55",
			Username:          "admin",
			Alias:             "vaultnode",
			EncryptedPassword: encPass,
		}
		if err := store.InsertSSHHost(ctx, host, mockDB.SQL()); err != nil {
			t.Fatalf("InsertSSHHost failed: %v", err)
		}

		resolvedByAlias, okAlias := lookupVaultPassword("vaultnode", "")
		if !okAlias || resolvedByAlias != pass {
			t.Errorf("expected resolved password %q by alias, got %q (ok=%v)", pass, resolvedByAlias, okAlias)
		}

		resolvedByIP, okIP := lookupVaultPassword("", "192.168.1.55")
		if !okIP || resolvedByIP != pass {
			t.Errorf("expected resolved password %q by IP, got %q (ok=%v)", pass, resolvedByIP, okIP)
		}

		_, okNotFound := lookupVaultPassword("unknownnode", "192.168.1.99")
		if okNotFound {
			t.Errorf("expected false for unknown node")
		}
	})
}
