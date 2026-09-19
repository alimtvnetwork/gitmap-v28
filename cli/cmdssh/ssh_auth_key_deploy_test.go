package cmdssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAuthKeyArgs_Defaults(t *testing.T) {
	target, keyPath, isUnix, err := parseAuthKeyArgs([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "all" || keyPath != "" || isUnix {
		t.Fatalf("expected all, empty keyPath, false unix, got target=%s, keyPath=%s, unix=%v", target, keyPath, isUnix)
	}
}

func TestParseAuthKeyArgs_WithDeployAndTarget(t *testing.T) {
	target, keyPath, _, err := parseAuthKeyArgs([]string{"deploy", "node1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "node1" || keyPath != "" {
		t.Fatalf("expected node1, got target=%s, keyPath=%s", target, keyPath)
	}
}

func TestParseAuthKeyArgs_WithIdentityFlag(t *testing.T) {
	target, keyPath, _, err := parseAuthKeyArgs([]string{"-i", "/custom/key.pub", "node1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "node1" || keyPath != "/custom/key.pub" {
		t.Fatalf("expected node1 and /custom/key.pub, got target=%s, keyPath=%s", target, keyPath)
	}
}

func TestParseAuthKeyArgs_WithIdentityEqualsFlag(t *testing.T) {
	target, keyPath, _, err := parseAuthKeyArgs([]string{"--identity=/custom/key.pub", "--all"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "all" || keyPath != "/custom/key.pub" {
		t.Fatalf("expected all and /custom/key.pub, got target=%s, keyPath=%s", target, keyPath)
	}
}

func TestParseAuthKeyArgs_MissingIdentityValue(t *testing.T) {
	_, _, _, err := parseAuthKeyArgs([]string{"-i"})
	hasError := err != nil
	if !hasError {
		t.Fatal("expected error for missing identity value, got nil")
	}
}

func TestParseAuthKeyArgs_FixAuthMultipleTargets(t *testing.T) {
	target, _, _, err := parseAuthKeyArgs([]string{"fix-auth", "machineid,", "ip,", "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	targets := ParseMultiIPList(target)
	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d: %v", len(targets), targets)
	}
	if targets[0] != "machineid" || targets[1] != "ip" || targets[2] != "id" {
		t.Fatalf("unexpected targets: %v", targets)
	}
}

func TestParseAuthKeyArgs_UnixFlag(t *testing.T) {
	target, _, isUnix, err := parseAuthKeyArgs([]string{"node1", "--unix"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "node1" || !isUnix {
		t.Fatalf("expected node1 and isUnix=true, got %s / %v", target, isUnix)
	}
}

func TestBuildUnixAuthKeyScript(t *testing.T) {
	key := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGitMapTestKey user@host"
	script := buildUnixAuthKeyScript(key)
	if !strings.Contains(script, "chmod 700 ~/.ssh") {
		t.Fatalf("expected chmod 700, got: %s", script)
	}
	if !strings.Contains(script, "chmod 600 ~/.ssh/authorized_keys") {
		t.Fatalf("expected chmod 600, got: %s", script)
	}
	if !strings.Contains(script, "grep -qF '"+key+"'") {
		t.Fatalf("expected grep -qF deduplication, got: %s", script)
	}
}

func TestBuildWindowsAuthKeyScript(t *testing.T) {
	key := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGitMapTestKey user@host"
	script := buildWindowsAuthKeyScript(key)
	if !strings.Contains(script, "authorized_keys") {
		t.Fatalf("expected authorized_keys reference, got: %s", script)
	}
	if !strings.Contains(script, key) {
		t.Fatalf("expected key in powershell script, got: %s", script)
	}
}

func TestLoadSpecifiedPublicKey_FileExists(t *testing.T) {
	tmpDir := t.TempDir()
	pubFile := filepath.Join(tmpDir, "test.pub")
	keyData := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITestPublicKey"
	_ = os.WriteFile(pubFile, []byte(keyData), 0600)

	content, resolved, err := loadSpecifiedPublicKey(pubFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != keyData || resolved != pubFile {
		t.Fatalf("mismatched content (%s) or path (%s)", content, resolved)
	}
}

func TestLoadSpecifiedPublicKey_AutoAppendsPub(t *testing.T) {
	tmpDir := t.TempDir()
	baseFile := filepath.Join(tmpDir, "id_custom")
	pubFile := baseFile + ".pub"
	keyData := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCTestKey"
	_ = os.WriteFile(pubFile, []byte(keyData), 0600)

	content, resolved, err := loadSpecifiedPublicKey(baseFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != keyData || resolved != pubFile {
		t.Fatalf("expected auto-append .pub to resolve, got %s / %s", content, resolved)
	}
}

func TestStripSubcommandKeyword(t *testing.T) {
	res := stripSubcommandKeyword([]string{"deploy", "all"}, "deploy")
	if len(res) != 1 || res[0] != "all" {
		t.Fatalf("expected [all], got %v", res)
	}
	res2 := stripSubcommandKeyword([]string{"node1"}, "deploy")
	if len(res2) != 1 || res2[0] != "node1" {
		t.Fatalf("expected [node1], got %v", res2)
	}
}

func TestIsIdentityFlag(t *testing.T) {
	if !isIdentityFlag("-i") || !isIdentityFlag("--identity") || !isIdentityFlag("--identity-file") {
		t.Fatal("expected identity flags to match")
	}
	if !isIdentityFlag("-i=key.pub") || !isIdentityFlag("--identity=key.pub") {
		t.Fatal("expected identity flag with equals to match")
	}
	if isIdentityFlag("--other") || isIdentityFlag("target") {
		t.Fatal("non-identity flags should not match")
	}
}

func TestResolveTargetFromArg(t *testing.T) {
	if resolveTargetFromArg("--all", "node1") != "all" {
		t.Fatal("expected all for --all")
	}
	if resolveTargetFromArg("-a", "node1") != "all" {
		t.Fatal("expected all for -a")
	}
	if resolveTargetFromArg("node2", "all") != "node2" {
		t.Fatal("expected node2 for positional target")
	}
}

func TestStripSubcommandKeyword_FixAuth(t *testing.T) {
	res := stripSubcommandKeyword([]string{"fix-auth", "m1"}, "deploy")
	if len(res) != 1 || res[0] != "m1" {
		t.Fatalf("expected [m1], got %v", res)
	}
}

func TestResolveTargetFromArg_Multiple(t *testing.T) {
	target := resolveTargetFromArg("m1,", "all")
	target = resolveTargetFromArg("ip,", target)
	target = resolveTargetFromArg("id", target)
	targets := ParseMultiIPList(target)
	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d: %v", len(targets), targets)
	}
}

func TestResolveTargetOS(t *testing.T) {
	if resolveTargetOS("windows", false) != "windows" {
		t.Fatal("expected windows when not forced")
	}
	if resolveTargetOS("windows", true) != "linux" {
		t.Fatal("expected linux when forced")
	}
	if resolveTargetOS("", false) != "linux" {
		t.Fatal("expected linux default")
	}
}
