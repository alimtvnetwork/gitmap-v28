package cmdlogin

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestParseLoginArgs_Flags(t *testing.T) {
	opts, rest, err := parseLoginArgs([]string{"--browser", "--no-verify", "--status"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.isBrowser || !opts.isStatus || opts.isVerify || len(rest) != 0 {
		t.Errorf("unexpected options parsed: %+v, rest: %v", opts, rest)
	}
}

func TestParseLoginArgs_TokenEquals(t *testing.T) {
	opts, rest, err := parseLoginArgs([]string{"--token=my-secret-token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.hasToken || opts.tokenValue != "my-secret-token" || len(rest) != 0 {
		t.Errorf("unexpected options parsed: %+v", opts)
	}
}

func TestParseLoginArgs_TokenSeparateArg(t *testing.T) {
	opts, rest, err := parseLoginArgs([]string{"--token", "my-secret-token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.hasToken || opts.tokenValue != "my-secret-token" || len(rest) != 0 {
		t.Errorf("unexpected options parsed: %+v", opts)
	}
}

func TestParseLoginArgs_MissingTokenValue(t *testing.T) {
	_, _, err := parseLoginArgs([]string{"--token"})
	if err == nil {
		t.Fatalf("expected error for missing token value")
	}
}

func TestParseSSHArgs(t *testing.T) {
	opts := parseSSHArgs([]string{"ssh", "work", "--yes", "--dir", "/tmp/keys"})
	if opts.alias != "work" || !opts.isOverwrite || opts.customDir != "/tmp/keys" {
		t.Errorf("unexpected ssh options: %+v", opts)
	}

	optsPrefix := parseSSHArgs([]string{"--dir=/tmp/custom", "-f"})
	if optsPrefix.customDir != "/tmp/custom" || !optsPrefix.isOverwrite {
		t.Errorf("unexpected ssh options with prefix: %+v", optsPrefix)
	}
}

func TestMaskLoginToken(t *testing.T) {
	short := maskLoginToken("short")
	if short != constants.TokenMaskFull {
		t.Errorf("expected %q, got %q", constants.TokenMaskFull, short)
	}

	long := maskLoginToken("ghp_1234567890abcdef")
	if !strings.HasPrefix(long, "ghp_") || !strings.HasSuffix(long, "cdef") {
		t.Errorf("unexpected masked token: %q", long)
	}
}

func TestResolveSSHKeyFileName(t *testing.T) {
	defaultName := resolveSSHKeyFileName("")
	if defaultName != constants.DefaultSSHKeyBaseName {
		t.Errorf("expected %q, got %q", constants.DefaultSSHKeyBaseName, defaultName)
	}

	aliasName := resolveSSHKeyFileName("work")
	if aliasName != "id_ed25519_work" {
		t.Errorf("expected 'id_ed25519_work', got %q", aliasName)
	}
}

func TestResolveSSHProfileName(t *testing.T) {
	cfg := model.GitProfileConfig{Active: "personal"}
	if name := resolveSSHProfileName("work", cfg); name != "work" {
		t.Errorf("expected 'work', got %q", name)
	}

	if name := resolveSSHProfileName("", cfg); name != "personal" {
		t.Errorf("expected 'personal', got %q", name)
	}

	emptyCfg := model.GitProfileConfig{}
	if name := resolveSSHProfileName("", emptyCfg); name != constants.DefaultProfileName {
		t.Errorf("expected %q, got %q", constants.DefaultProfileName, name)
	}
}
