package cmdssh

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestSSHEnableFlagParsing(t *testing.T) {
	opts, err := parseSSHEnableFlags([]string{"--port", "2222", "--force"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.port != 2222 {
		t.Errorf("expected port 2222, got %d", opts.port)
	}

	if !opts.isForce {
		t.Errorf("expected isForce to be true")
	}

	defaultOpts, err := parseSSHEnableFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error for default args: %v", err)
	}

	if defaultOpts.port != 22 {
		t.Errorf("expected default port 22, got %d", defaultOpts.port)
	}

	eqOpts, err := parseSSHEnableFlags([]string{"--port=8022", "-f"})
	if err != nil {
		t.Fatalf("unexpected error for --port=8022: %v", err)
	}

	if eqOpts.port != 8022 || !eqOpts.isForce {
		t.Errorf("expected port 8022 and isForce true, got %d, %v", eqOpts.port, eqOpts.isForce)
	}
}

func TestSSHEnableFlagValidationErrors(t *testing.T) {
	_, errZero := parseSSHEnableFlags([]string{"--port", "0"})
	if errZero == nil {
		t.Errorf("expected error for port 0")
	}

	_, errHigh := parseSSHEnableFlags([]string{"--port", "70000"})
	if errHigh == nil {
		t.Errorf("expected error for port 70000")
	}

	_, errMissing := parseSSHEnableFlags([]string{"--port"})
	if errMissing == nil {
		t.Errorf("expected error for missing port argument")
	}

	_, errInvalid := parseSSHEnableFlags([]string{"--port", "abc"})
	if errInvalid == nil {
		t.Errorf("expected error for non-integer port")
	}
}

func TestSSHPortValidation(t *testing.T) {
	port, isHelp, err := parsePortArgument([]string{"2222"})
	if err != nil || isHelp || port != 2222 {
		t.Errorf("expected 2222, false, nil; got %d, %v, %v", port, isHelp, err)
	}

	flagPort, isHelp2, err2 := parsePortArgument([]string{"--port", "3333"})
	if err2 != nil || isHelp2 || flagPort != 3333 {
		t.Errorf("expected 3333, false, nil; got %d, %v, %v", flagPort, isHelp2, err2)
	}

	_, isHelp3, _ := parsePortArgument([]string{"--help"})
	if !isHelp3 {
		t.Errorf("expected isHelp to be true")
	}

	_, _, errOutRange := parsePortArgument([]string{"99999"})
	if errOutRange == nil {
		t.Errorf("expected error for port 99999")
	}
}

func TestSSHDConfigUpdate(t *testing.T) {
	original := "# Configuration file\n#Port 22\nPermitRootLogin yes\n"
	updated := updateSSHDConfigContent(original, 2222)
	if !strings.Contains(updated, "Port 2222") {
		t.Errorf("expected updated config to contain 'Port 2222', got:\n%s", updated)
	}

	originalExisting := "Port 22\nPermitRootLogin yes\n"
	updatedExisting := updateSSHDConfigContent(originalExisting, 8022)
	if !strings.Contains(updatedExisting, "Port 8022") {
		t.Errorf("expected updated config to contain 'Port 8022', got:\n%s", updatedExisting)
	}

	originalNone := "PermitRootLogin yes\nPasswordAuthentication yes\n"
	updatedNone := updateSSHDConfigContent(originalNone, 4444)
	if !strings.Contains(updatedNone, "Port 4444") {
		t.Errorf("expected appended config to contain 'Port 4444', got:\n%s", updatedNone)
	}
}

func TestSSHTroubleshootDiagnosis(t *testing.T) {
	diagDrop := classifyTroubleshootDiagnosis(false, true, false, true)
	if diagDrop != DiagFirewallInboundDrop {
		t.Errorf("expected %s, got %s", DiagFirewallInboundDrop, diagDrop)
	}

	diagOffline := classifyTroubleshootDiagnosis(false, true, false, false)
	if diagOffline != DiagHostOffline {
		t.Errorf("expected %s, got %s", DiagHostOffline, diagOffline)
	}

	diagRefused := classifyTroubleshootDiagnosis(false, false, true, true)
	if diagRefused != DiagDaemonNotRunning {
		t.Errorf("expected %s, got %s", DiagDaemonNotRunning, diagRefused)
	}

	diagOpen := classifyTroubleshootDiagnosis(true, false, false, true)
	if diagOpen != DiagSSHPortOpen {
		t.Errorf("expected %s, got %s", DiagSSHPortOpen, diagOpen)
	}
}

func TestTroubleshootTargetParsing(t *testing.T) {
	ip1, port1 := parseTroubleshootTarget("192.168.1.50")
	if ip1 != "192.168.1.50" || port1 != 22 {
		t.Errorf("expected 192.168.1.50:22, got %s:%d", ip1, port1)
	}

	ip2, port2 := parseTroubleshootTarget("10.0.0.5:2222")
	if ip2 != "10.0.0.5" || port2 != 2222 {
		t.Errorf("expected 10.0.0.5:2222, got %s:%d", ip2, port2)
	}

	ip3, port3 := parseTroubleshootTarget("admin@172.16.0.10:8022")
	if ip3 != "172.16.0.10" || port3 != 8022 {
		t.Errorf("expected 172.16.0.10:8022, got %s:%d", ip3, port3)
	}
}

func TestSSHDispatchRouting(t *testing.T) {
	ctx := context.Background()

	enableSubs := []string{"enable", "enable-server", "sshd", "enable-sshd"}
	for _, sub := range enableSubs {
		res := dispatchDaemonSSH(ctx, sub, []string{"--help"})
		if !res.IsMatched() {
			t.Errorf("expected %s to be matched in dispatchDaemonSSH", sub)
		}
	}

	portSubs := []string{"port", "ports", "set-port"}
	for _, sub := range portSubs {
		res := dispatchDaemonSSH(ctx, sub, []string{"--help"})
		if !res.IsMatched() {
			t.Errorf("expected %s to be matched in dispatchDaemonSSH", sub)
		}
	}

	troubleshootSubs := []string{"troubleshoot", "doctor", "diagnose"}
	for _, sub := range troubleshootSubs {
		res := dispatchDaemonSSH(ctx, sub, []string{"--help"})
		if !res.IsMatched() {
			t.Errorf("expected %s to be matched in dispatchDaemonSSH", sub)
		}
	}

	resUnknown := dispatchDaemonSSH(ctx, "nonexistent-daemon-op", []string{})
	if resUnknown.IsMatched() {
		t.Errorf("expected unknown op to not be matched")
	}
}

func TestSSHClientErrorInterception(t *testing.T) {
	errTimeout := errors.New("ssh: connect to host 192.168.1.1 port 22: Connection timed out")
	if !isConnectionOrDaemonError(errTimeout, "") {
		t.Errorf("expected timeout error to be detected")
	}

	errRefused := errors.New("ssh: connect to host 192.168.1.1 port 22: Connection refused")
	if !isConnectionOrDaemonError(errRefused, "") {
		t.Errorf("expected connection refused error to be detected")
	}

	errExit255 := errors.New("Process exited with status 255")
	if !isConnectionOrDaemonError(errExit255, "") {
		t.Errorf("expected exit status 255 to be detected")
	}

	errStderr := errors.New("command failed")
	if !isConnectionOrDaemonError(errStderr, "ssh: connect to host: connection refused") {
		t.Errorf("expected stderr refused to be detected")
	}

	errRandom := errors.New("permission denied (publickey)")
	if isConnectionOrDaemonError(errRandom, "") {
		t.Errorf("permission denied should not be classified as daemon/connection error")
	}

	if isConnectionOrDaemonError(nil, "") {
		t.Errorf("nil error should return false")
	}
}
