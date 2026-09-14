package cmdssh

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func assertContains(t *testing.T, content, substr string) {
	t.Helper()
	hasSub := strings.Contains(content, substr)
	if !hasSub {
		t.Fatalf("expected content to contain %q, got:\n%s", substr, content)
	}
}

func assertNotContains(t *testing.T, content, substr string) {
	t.Helper()
	hasSub := strings.Contains(content, substr)
	if hasSub {
		t.Fatalf("expected content not to contain %q, got:\n%s", substr, content)
	}
}

func TestGenerateNetplanScript_Defaults(t *testing.T) {
	script := GenerateNetplanScript("192.168.0.10", "")
	assertContains(t, script, "192.168.0.10/24")
	assertContains(t, script, "via: 192.168.0.1")
	assertContains(t, script, "addresses: [8.8.8.8, 1.1.1.1]")
	assertContains(t, script, "netplan apply")
}

func TestGenerateNetplanScript_Custom(t *testing.T) {
	script := GenerateNetplanScript("10.0.0.5/16", "10.0.0.254")
	assertContains(t, script, "10.0.0.5/16")
	assertContains(t, script, "via: 10.0.0.254")
	assertContains(t, script, "netplan apply")
}

func TestGenerateBasePackagesScript(t *testing.T) {
	script := GenerateBasePackagesScript()
	assertContains(t, script, "apt-get update -y")
	assertContains(t, script, "curl wget git zsh net-tools htop build-essential")
}

func TestGenerateCreateUserScript_WithPassword(t *testing.T) {
	script := GenerateCreateUserScript("alice", "secret123", "agnoster")
	assertContains(t, script, "useradd -m -s /usr/bin/zsh alice")
	assertContains(t, script, "echo 'alice:secret123' | chpasswd")
	assertContains(t, script, "/etc/sudoers.d/alice")
	assertContains(t, script, "ZSH_THEME=\"agnoster\"")
}

func TestGenerateCreateUserScript_Defaults(t *testing.T) {
	script := GenerateCreateUserScript("kube", "", "")
	assertContains(t, script, "useradd -m -s /usr/bin/zsh kube")
	assertContains(t, script, "/etc/sudoers.d/kube")
	assertContains(t, script, "ZSH_THEME=\"fletcherm\"")
	assertNotContains(t, script, "chpasswd")
}

func TestGenerateThemeScript(t *testing.T) {
	scriptCustom := GenerateThemeScript("candy")
	assertContains(t, scriptCustom, "ZSH_THEME=\"candy\"")
	scriptDefault := GenerateThemeScript("")
	assertContains(t, scriptDefault, "ZSH_THEME=\"fletcherm\"")
}

func TestGeneratePurgeScript(t *testing.T) {
	script := GeneratePurgeScript()
	assertContains(t, script, "apt-get autoremove --purge -y")
	assertContains(t, script, "apt-get clean")
}

func TestParseSetIPArgs_Success(t *testing.T) {
	target, ip, gw, err := parseSetIPArgs([]string{"workers", "192.168.0.105", "--route-ip", "192.168.0.254"})
	if err != nil || target != "workers" || ip != "192.168.0.105" || gw != "192.168.0.254" {
		t.Fatalf("unexpected parse result: %s %s %s %v", target, ip, gw, err)
	}
}

func TestParseSetIPArgs_InlineFlag(t *testing.T) {
	target, ip, gw, err := parseSetIPArgs([]string{"k8s-w1", "192.168.0.105", "--route-ip=192.168.0.1"})
	if err != nil || target != "k8s-w1" || ip != "192.168.0.105" || gw != "192.168.0.1" {
		t.Fatalf("unexpected parse result: %s %s %s %v", target, ip, gw, err)
	}
}

func TestParseSetIPArgs_Error(t *testing.T) {
	_, _, _, err := parseSetIPArgs([]string{"only-one"})
	if err == nil {
		t.Fatalf("expected error for missing args")
	}
}

func TestParseCreateUserArgs_Success(t *testing.T) {
	target, user, pass, theme, err := parseCreateUserArgs([]string{"all", "dev", "pass123", "--theme", "bira"})
	if err != nil || target != "all" || user != "dev" || pass != "pass123" || theme != "bira" {
		t.Fatalf("unexpected parse result: %s %s %s %s %v", target, user, pass, theme, err)
	}
}

func TestParseCreateUserArgs_InlineTheme(t *testing.T) {
	target, user, pass, theme, err := parseCreateUserArgs([]string{"all", "dev", "--theme=bira"})
	if err != nil || target != "all" || user != "dev" || pass != "" || theme != "bira" {
		t.Fatalf("unexpected parse result: %s %s %s %s %v", target, user, pass, theme, err)
	}
}

func TestParseCreateUserArgs_Error(t *testing.T) {
	_, _, _, _, err := parseCreateUserArgs([]string{"all"})
	if err == nil {
		t.Fatalf("expected error for missing username")
	}
}

func TestParseSetThemeArgs(t *testing.T) {
	target, theme, err := parseSetThemeArgs([]string{"all", "robbyrussell"})
	if err != nil || target != "all" || theme != "robbyrussell" {
		t.Fatalf("unexpected result: %s %s %v", target, theme, err)
	}
	_, _, errShort := parseSetThemeArgs([]string{"all"})
	if errShort == nil {
		t.Fatalf("expected error for short args")
	}
}

func TestParseSingleTargetArg(t *testing.T) {
	target, err := parseSingleTargetArg("purge", []string{"workers"})
	if err != nil || target != "workers" {
		t.Fatalf("unexpected result: %s %v", target, err)
	}
	_, errShort := parseSingleTargetArg("purge", []string{})
	if errShort == nil {
		t.Fatalf("expected error for short args")
	}
}

func TestRunClusterNodeCLI_Help(t *testing.T) {
	errHelp := RunClusterNodeCLI([]string{"--help"})
	if errHelp != nil {
		t.Fatalf("expected nil on help, got: %v", errHelp)
	}
	errEmpty := RunClusterNodeCLI([]string{})
	if errEmpty != nil {
		t.Fatalf("expected nil on empty, got: %v", errEmpty)
	}
}

func TestRunClusterNodeCLI_Unknown(t *testing.T) {
	err := RunClusterNodeCLI([]string{"nonexistent-command"})
	if err == nil {
		t.Fatalf("expected error for unknown subcommand")
	}
}

func withMockNodeExecutor(fn func(calls *[]string)) {
	var calls []string
	prevExec := executeNodeCmdFn
	executeNodeCmdFn = func(ctx context.Context, host store.SSHHost, shellCmd string, isSudo bool) ClusterRunResult {
		calls = append(calls, fmt.Sprintf("%s:%t", host.Alias, isSudo))
		return ClusterRunResult{Host: host, ExitCode: 0}
	}
	defer func() { executeNodeCmdFn = prevExec }()
	fn(&calls)
}

func TestRunClusterNodeCLI_MockExecution(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		ctx := context.Background()
		seedClusterTestHosts(t, ctx, db.SQL())
		withMockNodeExecutor(func(calls *[]string) {
			err := RunClusterNodeCLI([]string{"install-base", "workers"})
			if err != nil {
				t.Fatalf("RunClusterNodeCLI failed: %v", err)
			}
			if len(*calls) != 2 {
				t.Fatalf("expected 2 worker calls, got %d", len(*calls))
			}
		})
	})
}
