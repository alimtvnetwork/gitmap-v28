package cmdports

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParsePortsFlags(t *testing.T) {
	opts, isHelp, err := parsePortsFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if isHelp || opts.TargetPort != 0 || opts.CommonOnly || opts.JSONOutput {
		t.Fatalf("expected zeroed options, got %+v", opts)
	}

	optsPort, _, err := parsePortsFlags([]string{"-p", "22"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if optsPort.TargetPort != 22 {
		t.Fatalf("expected TargetPort 22, got %d", optsPort.TargetPort)
	}

	optsEq, _, err := parsePortsFlags([]string{"--port=443"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if optsEq.TargetPort != 443 {
		t.Fatalf("expected TargetPort 443, got %d", optsEq.TargetPort)
	}

	optsCommon, _, err := parsePortsFlags([]string{"--common"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !optsCommon.CommonOnly {
		t.Fatalf("expected CommonOnly true")
	}

	optsJSON, _, err := parsePortsFlags([]string{"--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !optsJSON.JSONOutput {
		t.Fatalf("expected JSONOutput true")
	}

	_, isHelpFlag, err := parsePortsFlags([]string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !isHelpFlag {
		t.Fatalf("expected isHelp true")
	}

	_, _, errInvalidPort := parsePortsFlags([]string{"-p", "invalid"})
	if errInvalidPort == nil {
		t.Fatalf("expected error for invalid port")
	}

	_, _, errRangePort := parsePortsFlags([]string{"-p", "99999"})
	if errRangePort == nil {
		t.Fatalf("expected error for out of range port")
	}

	_, _, errUnknown := parsePortsFlags([]string{"--unknown-flag"})
	if errUnknown == nil {
		t.Fatalf("expected error for unknown flag")
	}
}

func TestCommonPortsCatalog(t *testing.T) {
	expected := []int{22, 80, 443, 3389, 5985, 5986, 8080}
	if len(commonPortsList) != len(expected) {
		t.Fatalf("expected %d common ports, got %d", len(expected), len(commonPortsList))
	}

	for i, p := range expected {
		if commonPortsList[i] != p {
			t.Errorf("expected port %d at index %d, got %d", p, i, commonPortsList[i])
		}
	}

	recSSHListening := resolvePortRecommendation(22, "LISTENING", "Allow")
	if !strings.Contains(recSSHListening, "SSH") {
		t.Errorf("expected recommendation to mention SSH, got: %s", recSSHListening)
	}

	recSSHClosed := resolvePortRecommendation(22, "CLOSED", "No Rule")
	if !strings.Contains(recSSHClosed, "gitmap ssh enable") {
		t.Errorf("expected recommendation to mention gitmap ssh enable, got: %s", recSSHClosed)
	}

	recRDP := resolvePortRecommendation(3389, "CLOSED", "No Rule")
	if !strings.Contains(recRDP, "Remote Desktop") {
		t.Errorf("expected recommendation to mention Remote Desktop, got: %s", recRDP)
	}
}

func TestRenderPortsTable(t *testing.T) {
	entries := []PortEntry{
		{
			Port:           22,
			Protocol:       "TCP",
			ProcessName:    "sshd.exe",
			PID:            1234,
			State:          "LISTENING",
			FirewallStatus: "Allow",
			Recommendation: "SSH service active and accessible",
		},
		{
			Port:           80,
			Protocol:       "TCP",
			ProcessName:    "-",
			PID:            0,
			State:          "CLOSED",
			FirewallStatus: "No Rule",
			Recommendation: "Start HTTP web service (IIS/Nginx/Apache)",
		},
	}

	cfg := buildPortsTableConfig(entries)
	if len(cfg.Columns) != 7 {
		t.Fatalf("expected 7 columns, got %d", len(cfg.Columns))
	}

	if len(cfg.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(cfg.Rows))
	}

	// Verify table config formatting doesn't panic
	renderPortsTable(entries, PortsOptions{CommonOnly: true})
}

func TestOutputPortsJSON(t *testing.T) {
	entries := []PortEntry{
		{
			Port:           22,
			Protocol:       "TCP",
			ProcessName:    "sshd.exe",
			PID:            1234,
			State:          "LISTENING",
			FirewallStatus: "Allow",
			Recommendation: "SSH service active and accessible",
		},
	}

	data, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	var parsed []PortEntry
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(parsed) != 1 || parsed[0].Port != 22 {
		t.Fatalf("JSON roundtrip mismatch: %+v", parsed)
	}
}
