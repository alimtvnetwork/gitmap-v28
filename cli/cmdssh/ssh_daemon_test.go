package cmdssh

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestSSHEnableCommandGenerationWindows(t *testing.T) {
	capCheck := buildWindowsCheckCapabilityCmd()
	if !strings.Contains(capCheck, "Get-WindowsCapability") {
		t.Errorf("expected Get-WindowsCapability in check cmd, got %s", capCheck)
	}

	capAdd := buildWindowsAddCapabilityCmd()
	if !strings.Contains(capAdd, "Add-WindowsCapability") {
		t.Errorf("expected Add-WindowsCapability in add cmd, got %s", capAdd)
	}

	fwRule := buildWindowsFirewallRuleCmd(2222)
	if !strings.Contains(fwRule, "New-NetFirewallRule") {
		t.Errorf("expected New-NetFirewallRule in rule, got %s", fwRule)
	}

	if !strings.Contains(fwRule, "-LocalPort 2222") {
		t.Errorf("expected -LocalPort 2222 in rule, got %s", fwRule)
	}
}

func isFirstArgMismatch(args []string, expected string) bool {
	if len(args) == 0 {
		return true
	}

	return args[0] != expected
}

func assertStepMatch(t *testing.T, step LinuxStep, expectedBin string, expectedFirstArg string) {
	t.Helper()
	if step.Binary != expectedBin {
		t.Errorf("expected binary %q, got %q", expectedBin, step.Binary)
	}

	if isFirstArgMismatch(step.Args, expectedFirstArg) {
		t.Errorf("expected first arg %q, got %v", expectedFirstArg, step.Args)
	}
}

func TestSSHEnableCommandGenerationLinux(t *testing.T) {
	debSteps := buildDebianEnableSteps(2222)
	if len(debSteps) != 3 {
		t.Fatalf("expected 3 debian steps, got %d", len(debSteps))
	}

	assertStepMatch(t, debSteps[0], "apt-get", "install")
	assertStepMatch(t, debSteps[1], "systemctl", "enable")
	assertStepMatch(t, debSteps[2], "ufw", "allow")

	rhelSteps := buildRhelEnableSteps("dnf", 2222)
	if len(rhelSteps) != 4 {
		t.Fatalf("expected 4 rhel steps, got %d", len(rhelSteps))
	}

	assertStepMatch(t, rhelSteps[0], "dnf", "install")
	assertStepMatch(t, rhelSteps[1], "systemctl", "enable")
	assertStepMatch(t, rhelSteps[2], "firewall-cmd", "--permanent")
	assertStepMatch(t, rhelSteps[3], "firewall-cmd", "--reload")
}

func TestSSHEnableCommandGeneration(t *testing.T) {
	TestSSHEnableCommandGenerationWindows(t)
	TestSSHEnableCommandGenerationLinux(t)
}

func assertValidPortArg(t *testing.T, args []string, expectedPort int) {
	t.Helper()
	port, isHelp, err := parsePortArgument(args)
	if err != nil {
		t.Fatalf("unexpected error for %v: %v", args, err)
	}

	if isHelp {
		t.Errorf("expected isHelp false for %v", args)
	}

	if port != expectedPort {
		t.Errorf("expected port %d, got %d for %v", expectedPort, port, args)
	}
}

func assertInvalidPortArg(t *testing.T, args []string, desc string) {
	t.Helper()
	_, _, err := parsePortArgument(args)
	if err == nil {
		t.Errorf("expected error for %s (%v)", desc, args)
	}
}

func TestSSHPortArgumentParsingAndBounds(t *testing.T) {
	assertValidPortArg(t, []string{"2222"}, 2222)
	assertValidPortArg(t, []string{"--port", "3333"}, 3333)
	assertValidPortArg(t, []string{"-p", "4444"}, 4444)
	assertValidPortArg(t, []string{"--port=5555"}, 5555)
	assertValidPortArg(t, []string{"1"}, 1)
	assertValidPortArg(t, []string{"65535"}, 65535)

	assertInvalidPortArg(t, []string{"0"}, "zero port")
	assertInvalidPortArg(t, []string{"-1"}, "negative port")
	assertInvalidPortArg(t, []string{"65536"}, "port 65536")
	assertInvalidPortArg(t, []string{"99999"}, "port 99999")
	assertInvalidPortArg(t, []string{"abc"}, "non-integer port")
	assertInvalidPortArg(t, []string{"--port"}, "missing port value")
}

func TestSSHPortValidation(t *testing.T) {
	TestSSHPortArgumentParsingAndBounds(t)
}

func TestSSHDConfigContentUpdate(t *testing.T) {
	cfgComment := "# Configuration\n#Port 22\nPermitRootLogin yes\n"
	updated1 := updateSSHDConfigContent(cfgComment, 2222)
	if !strings.Contains(updated1, "Port 2222") {
		t.Errorf("expected 'Port 2222' in updated config:\n%s", updated1)
	}

	cfgSpace := "# Configuration\n# Port 22\nPermitRootLogin yes\n"
	updated2 := updateSSHDConfigContent(cfgSpace, 3333)
	if !strings.Contains(updated2, "Port 3333") {
		t.Errorf("expected 'Port 3333' in updated config:\n%s", updated2)
	}

	cfgActive := "Port 22\nPermitRootLogin yes\n"
	updated3 := updateSSHDConfigContent(cfgActive, 8022)
	if !strings.Contains(updated3, "Port 8022") {
		t.Errorf("expected 'Port 8022' in updated config:\n%s", updated3)
	}

	cfgNone := "PermitRootLogin yes\n"
	updated4 := updateSSHDConfigContent(cfgNone, 4444)
	if !strings.Contains(updated4, "Port 4444") {
		t.Errorf("expected 'Port 4444' in appended config:\n%s", updated4)
	}
}

func TestSSHDConfigUpdate(t *testing.T) {
	TestSSHDConfigContentUpdate(t)
}

func assertHelpDetected(t *testing.T, flag string) {
	t.Helper()
	_, isPortHelp, _ := parsePortArgument([]string{flag})
	if !isPortHelp {
		t.Errorf("expected parsePortArgument to detect help for %s", flag)
	}

	enableOpts, _ := parseSSHEnableFlags([]string{flag})
	if !enableOpts.isHelp {
		t.Errorf("expected parseSSHEnableFlags to detect help for %s", flag)
	}
}

func TestSSHHelpArgumentDetection(t *testing.T) {
	assertHelpDetected(t, "--help")
	assertHelpDetected(t, "-h")
	assertHelpDetected(t, "help")
}

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

	defaultOpts, defErr := parseSSHEnableFlags([]string{})
	if defErr != nil {
		t.Fatalf("unexpected error for default args: %v", defErr)
	}

	if defaultOpts.port != 22 {
		t.Errorf("expected default port 22, got %d", defaultOpts.port)
	}
}

func assertEnableFlagError(t *testing.T, args []string, desc string) {
	t.Helper()
	_, err := parseSSHEnableFlags(args)
	if err == nil {
		t.Errorf("expected error for %s", desc)
	}
}

func TestSSHEnableFlagValidationErrors(t *testing.T) {
	assertEnableFlagError(t, []string{"--port", "0"}, "port 0")
	assertEnableFlagError(t, []string{"--port", "70000"}, "port 70000")
	assertEnableFlagError(t, []string{"--port"}, "missing port argument")
	assertEnableFlagError(t, []string{"--port", "abc"}, "non-integer port")
}

func TestSSHTroubleshootDiagnosis(t *testing.T) {
	if diag := classifyTroubleshootDiagnosis(false, true, false, true); diag != DiagFirewallInboundDrop {
		t.Errorf("expected %s, got %s", DiagFirewallInboundDrop, diag)
	}

	if diag := classifyTroubleshootDiagnosis(false, true, false, false); diag != DiagHostOffline {
		t.Errorf("expected %s, got %s", DiagHostOffline, diag)
	}

	if diag := classifyTroubleshootDiagnosis(false, false, true, true); diag != DiagDaemonNotRunning {
		t.Errorf("expected %s, got %s", DiagDaemonNotRunning, diag)
	}

	if diag := classifyTroubleshootDiagnosis(true, false, false, true); diag != DiagSSHPortOpen {
		t.Errorf("expected %s, got %s", DiagSSHPortOpen, diag)
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

func assertDispatchMatched(t *testing.T, sub string) {
	t.Helper()
	res := dispatchDaemonSSH(context.Background(), sub, []string{"--help"})
	if !res.IsMatched() {
		t.Errorf("expected %s to be matched in dispatchDaemonSSH", sub)
	}
}

func TestSSHDispatchRouting(t *testing.T) {
	subs := []string{"enable", "enable-server", "sshd", "enable-sshd", "port", "ports", "set-port", "troubleshoot", "doctor", "diagnose"}
	for _, sub := range subs {
		assertDispatchMatched(t, sub)
	}

	resUnknown := dispatchDaemonSSH(context.Background(), "nonexistent-daemon-op", []string{})
	if resUnknown.IsMatched() {
		t.Errorf("expected unknown op to not be matched")
	}
}

func assertDetectedConnError(t *testing.T, err error, stderr string, expected bool) {
	t.Helper()
	detected := isConnectionOrDaemonError(err, stderr)
	if detected != expected {
		t.Errorf("expected detection %v, got %v for err=%v stderr=%q", expected, detected, err, stderr)
	}
}

func TestSSHClientErrorInterception(t *testing.T) {
	errTimeout := errors.New("ssh: connect to host 192.168.1.1 port 22: Connection timed out")
	assertDetectedConnError(t, errTimeout, "", true)

	errRefused := errors.New("ssh: connect to host 192.168.1.1 port 22: Connection refused")
	assertDetectedConnError(t, errRefused, "", true)

	errExit255 := errors.New("Process exited with status 255")
	assertDetectedConnError(t, errExit255, "", true)

	errStderr := errors.New("command failed")
	assertDetectedConnError(t, errStderr, "ssh: connect to host: connection refused", true)

	errRandom := errors.New("permission denied (publickey)")
	assertDetectedConnError(t, errRandom, "", false)
	assertDetectedConnError(t, nil, "", false)
}
