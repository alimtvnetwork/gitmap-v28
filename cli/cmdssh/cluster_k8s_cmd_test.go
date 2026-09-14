package cmdssh

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func withMockK8sExecutor(fn func(calls *[]string)) {
	var calls []string
	prevExec := executeNodeCmdFn
	executeNodeCmdFn = func(ctx context.Context, host store.SSHHost, shellCmd string, isSudo bool) ClusterRunResult {
		calls = append(calls, fmt.Sprintf("%s:%s", host.Alias, shellCmd))
		return ClusterRunResult{
			Host:     host,
			ExitCode: 0,
			Stdout:   "kubeadm join 192.168.1.10:6443 --token abc.123 --discovery-token-ca-cert-hash sha256:456",
		}
	}
	defer func() { executeNodeCmdFn = prevExec }()
	fn(&calls)
}

func TestExtractKubeadmJoinCommand_Success(t *testing.T) {
	out := "kubeadm join 192.168.0.10:6443 --token abcdef.123456 --discovery-token-ca-cert-hash sha256:fedcba"
	cmd, err := extractKubeadmJoinCommand(out)
	if err != nil || !strings.Contains(cmd, "kubeadm join 192.168.0.10:6443") {
		t.Fatalf("unexpected extract result: %q, err: %v", cmd, err)
	}
}

func TestExtractKubeadmJoinCommand_MultiLine(t *testing.T) {
	out := "info: token created\nkubeadm join 10.0.0.1:6443 --token xyz.987 \\\n\t--discovery-token-ca-cert-hash sha256:112233"
	cmd, err := extractKubeadmJoinCommand(out)
	if err != nil || !strings.Contains(cmd, "10.0.0.1:6443") || !strings.Contains(cmd, "sha256:112233") {
		t.Fatalf("unexpected extract result: %q, err: %v", cmd, err)
	}
}

func TestExtractKubeadmJoinCommand_Error(t *testing.T) {
	_, err := extractKubeadmJoinCommand("Error: preflight check failed")
	if err == nil {
		t.Fatalf("expected error for missing join command")
	}
}

func TestParseSingleK8sTarget_Success(t *testing.T) {
	target, err := parseSingleK8sTarget("prereq", []string{"workers"})
	if err != nil || target != "workers" {
		t.Fatalf("unexpected single target result: %s, err: %v", target, err)
	}
}

func TestParseSingleK8sTarget_Missing(t *testing.T) {
	_, err := parseSingleK8sTarget("prereq", []string{})
	if err == nil {
		t.Fatalf("expected validation error for empty args")
	}
}

func TestParseTargetAndFlag_Inline(t *testing.T) {
	args := []string{"control", "--pod-cidr=10.244.0.0/16"}
	target, flagVal, err := parseTargetAndFlag(args, "pod-cidr", "", "usage")
	if err != nil || target != "control" || flagVal != "10.244.0.0/16" {
		t.Fatalf("unexpected target and flag result: %s, %s, err: %v", target, flagVal, err)
	}
}

func TestParseTargetAndFlag_Short(t *testing.T) {
	args := []string{"k8s-m1", "-p", "calico"}
	target, flagVal, err := parseTargetAndFlag(args, "plugin", "p", "usage")
	if err != nil || target != "k8s-m1" || flagVal != "calico" {
		t.Fatalf("unexpected short flag result: %s, %s, err: %v", target, flagVal, err)
	}
}

func TestParseTargetAndFlag_MissingTarget(t *testing.T) {
	_, _, err := parseTargetAndFlag([]string{"-p", "calico"}, "plugin", "p", "usage")
	if err == nil {
		t.Fatalf("expected error for missing target positional")
	}
}

func TestParseTargetAndTwoFlags_Success(t *testing.T) {
	args := []string{"workers", "--version=1.31", "--hostname", "worker1"}
	target, v1, v2, err := parseTargetAndTwoFlags(args, "version", "v", "hostname", "", "usage")
	if err != nil || target != "workers" || v1 != "1.31" || v2 != "worker1" {
		t.Fatalf("unexpected two flags result: %s, %s, %s, err: %v", target, v1, v2, err)
	}
}

func TestParseTargetAndTwoFlags_MissingTarget(t *testing.T) {
	_, _, _, err := parseTargetAndTwoFlags([]string{"--version", "1.31"}, "version", "v", "hostname", "", "usage")
	if err == nil {
		t.Fatalf("expected error for missing target")
	}
}

func TestRunClusterK8sCLI_Help(t *testing.T) {
	if err := RunClusterK8sCLI([]string{"--help"}); err != nil {
		t.Fatalf("expected nil for help, got: %v", err)
	}
	if err := RunClusterK8sCLI([]string{}); err != nil {
		t.Fatalf("expected nil for empty args, got: %v", err)
	}
}

func TestRunClusterK8sCLI_UnknownCommand(t *testing.T) {
	err := RunClusterK8sCLI([]string{"unknown-subcmd"})
	if err == nil {
		t.Fatalf("expected error for unknown subcommand")
	}
}

func TestRunClusterK8sCLI_PrereqMock(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		seedClusterTestHosts(t, context.Background(), db.SQL())
		withMockK8sExecutor(func(calls *[]string) {
			err := RunClusterK8sCLI([]string{"prereq", "workers"})
			if err != nil || len(*calls) != 2 {
				t.Fatalf("expected 2 calls, got %d, err: %v", len(*calls), err)
			}
		})
	})
}

func TestRunClusterK8sCLI_AutoJoinPipeline(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		seedClusterTestHosts(t, context.Background(), db.SQL())
		withMockK8sExecutor(func(calls *[]string) {
			err := RunClusterK8sCLI([]string{"join", "workers"})
			if err != nil || len(*calls) != 3 {
				t.Fatalf("expected 3 calls (1 control + 2 workers), got %d, err: %v", len(*calls), err)
			}
			hasTokenCall := strings.Contains((*calls)[0], "kubeadm token create")
			if !hasTokenCall {
				t.Fatalf("expected first call on master for token, got: %s", (*calls)[0])
			}
		})
	})
}

func TestRunClusterK8sCLI_ManualJoin(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		seedClusterTestHosts(t, context.Background(), db.SQL())
		withMockK8sExecutor(func(calls *[]string) {
			err := RunClusterK8sCLI([]string{"join", "workers", "--command", "kubeadm join 1.2.3.4:6443 --token tok"})
			if err != nil || len(*calls) != 2 {
				t.Fatalf("expected 2 calls on workers directly, got %d, err: %v", len(*calls), err)
			}
		})
	})
}

func TestRunClusterK8sCLI_TargetNotFound(t *testing.T) {
	withMockClusterStoreDB(t, func(db *store.DB) {
		seedClusterTestHosts(t, context.Background(), db.SQL())
		err := RunClusterK8sCLI([]string{"status", "nonexistent-target"})
		if err == nil {
			t.Fatalf("expected not found error for nonexistent target")
		}
	})
}
