package cmdssh

import (
	"context"
	"strings"
	"testing"
)

func TestParseCommonIPTokens_ShorthandOctets(t *testing.T) {
	rawArgs := []string{"192.168.1.3(w1),7(w2),12(w3)"}
	tokens := ExtractRawTokens(rawArgs)

	targets, err := ParseCommonIPTokens("administrator", 22, tokens)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
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

func TestParseCommonIPTokens_SubnetSwitching(t *testing.T) {
	rawArgs := []string{"192.168.1.3(w1)", "7(w2)", "10.0.0.5(u1)", "8(u2)"}
	tokens := ExtractRawTokens(rawArgs)

	targets, err := ParseCommonIPTokens("admin", 2222, tokens)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if len(targets) != 4 {
		t.Fatalf("expected 4 targets, got %d", len(targets))
	}

	expectedIPs := []string{"192.168.1.3", "192.168.1.7", "10.0.0.5", "10.0.0.8"}
	for i, expIP := range expectedIPs {
		if targets[i].FullIP != expIP {
			t.Errorf("[%d] expected IP %s, got %s", i, expIP, targets[i].FullIP)
		}
		if targets[i].Port != 2222 {
			t.Errorf("[%d] expected port 2222, got %d", i, targets[i].Port)
		}
	}
}

func TestParseCommonIPTokens_MissingAlias(t *testing.T) {
	tokens := []string{"192.168.1.50"}
	targets, err := ParseCommonIPTokens("root", 22, tokens)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if targets[0].Alias != "node-192-168-1-50" {
		t.Errorf("expected default alias node-192-168-1-50, got %s", targets[0].Alias)
	}
}

func TestParseCommonIPTokens_InvalidFirstToken(t *testing.T) {
	tokens := []string{"7(w1)"}
	_, err := ParseCommonIPTokens("admin", 22, tokens)
	if err == nil {
		t.Fatal("expected error when first token is octet without prefix, got nil")
	}
}

func TestParseSJCFlags(t *testing.T) {
	args := []string{"administrator", "192.168.1.3(w1),7(w2)", "--pass", "secret123", "--port", "2222", "--json", "-d"}
	opts, err := parseSJCFlags(args)
	if err != nil {
		t.Fatalf("unexpected flag parse error: %v", err)
	}

	if opts.Username != "administrator" {
		t.Errorf("expected username administrator, got %s", opts.Username)
	}
	if opts.Password != "secret123" {
		t.Errorf("expected password secret123, got %s", opts.Password)
	}
	if opts.Port != 2222 {
		t.Errorf("expected port 2222, got %d", opts.Port)
	}
	if !opts.IsJSON {
		t.Error("expected IsJSON true")
	}
	if !opts.DryRun {
		t.Error("expected DryRun true")
	}
}

func TestExecuteCommonJoin_DryRun(t *testing.T) {
	var buf strings.Builder
	opts := SSHCommonJoinOptions{
		Username: "administrator",
		RawIPs:   []string{"192.168.1.3(w1),7(w2),12(w3)"},
		Password: "test-password",
		DryRun:   true,
		IsJSON:   true,
	}

	res, err := ExecuteCommonJoin(context.Background(), &buf, opts)
	if err != nil {
		t.Fatalf("unexpected ExecuteCommonJoin error: %v", err)
	}

	if res.TotalCount != 3 {
		t.Fatalf("expected 3 total targets, got %d", res.TotalCount)
	}
	if res.SuccessCount != 3 {
		t.Fatalf("expected 3 successful targets in dry-run, got %d", res.SuccessCount)
	}
	if res.FailureCount != 0 {
		t.Fatalf("expected 0 failures in dry-run, got %d", res.FailureCount)
	}
}
