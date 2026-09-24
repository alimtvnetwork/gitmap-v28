package cmdssh

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestParseCommonIPTokens_ShorthandAndAliases(t *testing.T) {
	tokens := []string{"192.168.1.3(w1)", "7(w2)", "12(w3)"}
	targets, err := ParseCommonIPTokens("administrator", 22, tokens)
	if err != nil {
		t.Fatalf("unexpected error parsing tokens: %v", err)
	}

	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(targets))
	}

	expected := []struct {
		ip    string
		alias string
	}{
		{"192.168.1.3", "w1"},
		{"192.168.1.7", "w2"},
		{"192.168.1.12", "w3"},
	}

	for i, exp := range expected {
		if targets[i].FullIP != exp.ip {
			t.Errorf("[%d] expected IP %s, got %s", i, exp.ip, targets[i].FullIP)
		}
		if targets[i].Alias != exp.alias {
			t.Errorf("[%d] expected alias %s, got %s", i, exp.alias, targets[i].Alias)
		}
		if targets[i].Username != "administrator" {
			t.Errorf("[%d] expected username administrator, got %s", i, targets[i].Username)
		}
	}
}

func TestRunSSHJoinCommonCLI_DryRun(t *testing.T) {
	args := []string{"administrator", "192.168.1.3(w1),7(w2)", "--dry-run", "--pass", "secret123"}
	opts, err := parseSJCFlags(args)
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}

	var buf bytes.Buffer
	res, appErr := ExecuteCommonJoin(context.Background(), &buf, opts)
	if appErr != nil {
		t.Fatalf("unexpected error running dry-run: %v", appErr)
	}

	if res.TotalCount != 2 {
		t.Fatalf("expected 2 total targets, got %d", res.TotalCount)
	}

	out := buf.String()
	if !strings.Contains(out, "192.168.1.3") || !strings.Contains(out, "192.168.1.7") {
		t.Errorf("expected output to contain expanded IPs, got: %s", out)
	}
	if !strings.Contains(out, "SUCCESS") {
		t.Errorf("expected dry-run to mark targets as success, got: %s", out)
	}
}

func TestResolveCommonPassword_DryRun(t *testing.T) {
	pwd, err := resolveCommonPassword(context.Background(), "", true)
	if err != nil {
		t.Fatalf("unexpected error in dry-run password resolution: %v", err)
	}
	if pwd != "" {
		t.Errorf("expected empty password for dry run, got: %s", pwd)
	}
}
