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

	if isHelp || opts.TargetPort != 0 || opts.CommonOnly || opts.FirewallOnly || opts.JSONOutput {
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

func TestParseServiceNames(t *testing.T) {
	testCases := []struct {
		arg          string
		expectedPort int
	}{
		{"ssh", 22},
		{"sshd", 22},
		{"SSH", 22},
		{"rdp", 3389},
		{"RDP", 3389},
		{"winrm", 5985},
		{"WinRM", 5985},
		{"winrm-https", 5986},
		{"winrm-ssl", 5986},
		{"http", 80},
		{"web", 80},
		{"https", 443},
		{"ssl", 443},
		{"dns", 53},
		{"smb", 445},
		{"mysql", 3306},
		{"postgres", 5432},
		{"postgresql", 5432},
		{"redis", 6379},
	}

	for _, tc := range testCases {
		opts, isHelp, err := parsePortsFlags([]string{tc.arg})

		if err != nil {
			t.Fatalf("unexpected error parsing %s: %v", tc.arg, err)
		}

		if isHelp {
			t.Fatalf("expected isHelp to be false for %s", tc.arg)
		}

		if opts.TargetPort != tc.expectedPort {
			t.Fatalf("expected port %d for %s, got %d", tc.expectedPort, tc.arg, opts.TargetPort)
		}
	}
}

func TestParseFirewallFlags(t *testing.T) {
	firewallArgs := []string{"firewall", "fw", "--firewall"}

	for _, arg := range firewallArgs {
		opts, isHelp, err := parsePortsFlags([]string{arg})

		if err != nil {
			t.Fatalf("unexpected error parsing firewall flag %s: %v", arg, err)
		}

		if isHelp {
			t.Fatalf("expected isHelp false for %s", arg)
		}

		if !opts.CommonOnly {
			t.Fatalf("expected CommonOnly to be true for %s", arg)
		}

		if !opts.FirewallOnly {
			t.Fatalf("expected FirewallOnly to be true for %s", arg)
		}
	}
}

func TestResolveServiceNamePort(t *testing.T) {
	port, isFound := resolveServiceNamePort("ssh")

	if !isFound || port != 22 {
		t.Fatalf("expected 22, true for ssh, got %d, %v", port, isFound)
	}

	portSSHUpper, isFoundUpper := resolveServiceNamePort("SSH")

	if !isFoundUpper || portSSHUpper != 22 {
		t.Fatalf("expected 22, true for SSH, got %d, %v", portSSHUpper, isFoundUpper)
	}

	portRDP, isFoundRDP := resolveServiceNamePort("Rdp")

	if !isFoundRDP || portRDP != 3389 {
		t.Fatalf("expected 3389, true for Rdp, got %d, %v", portRDP, isFoundRDP)
	}

	_, isFoundUnknown := resolveServiceNamePort("nonexistent-service")

	if isFoundUnknown {
		t.Fatalf("expected false for unknown service")
	}
}

