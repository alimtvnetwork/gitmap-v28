//go:build e2e

package cmdagy

import (
	"os"
	"strings"
	"testing"
)

func skipIfInCI(t *testing.T) {
	isCI := os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
	if isCI {
		t.Skip("skipping adhoc E2E test in CI/CD environment")
	}
}

func TestE2E_AntigravityIDEDetection(t *testing.T) {
	skipIfInCI(t)
	procRes := DetectRunningAntigravityIDE()
	if procRes.IsFailure() {
		t.Logf("Antigravity IDE is not running: %v", procRes.Err)
		return
	}
	t.Logf("Antigravity IDE detected: PID=%d, Name=%s", procRes.Value.PID, procRes.Value.Name)
}

func TestE2E_AntigravityCLIDetection(t *testing.T) {
	skipIfInCI(t)
	cliRes := ResolveAntigravityCLI()
	if cliRes.IsFailure() {
		t.Logf("Antigravity CLI (agy) not found in PATH or standard locations")
		return
	}
	t.Logf("Antigravity CLI (agy) found at: %s", cliRes.Value)
}

func TestE2E_AgentAPIResolution(t *testing.T) {
	skipIfInCI(t)
	binPath, args, isResolved := ResolveAgentAPI()
	if !isResolved {
		t.Logf("agentapi not resolved on this machine")
		return
	}
	t.Logf("agentapi resolved: bin=%s, args=%v", binPath, args)
}

func TestE2E_WindowFocus(t *testing.T) {
	skipIfInCI(t)
	procRes := DetectRunningAntigravityIDE()
	if procRes.IsFailure() {
		t.Skip("Antigravity IDE not running, skipping window focus")
	}
	FocusAntigravityWindow(procRes.Value.PID)
	t.Logf("FocusAntigravityWindow invoked for PID %d", procRes.Value.PID)
}

func TestE2E_LSEnvDiscovery(t *testing.T) {
	skipIfInCI(t)
	addr, token, ok := ResolveAntigravityLSEnv()
	if !ok {
		t.Fatalf("failed to discover Antigravity Language Server address and token")
	}
	t.Logf("Discovered Language Server: addr=%s, tokenLen=%d", addr, len(token))
}

func TestE2E_AgentAPIExecution(t *testing.T) {
	skipIfInCI(t)
	addr, _, ok := ResolveAntigravityLSEnv()
	if !ok {
		t.Skip("Antigravity Language Server not running, skipping execution test")
	}
	t.Logf("Testing agentapi with addr=%s...", addr)
	rawRes := executeAgentAPICmd([]string{"get-conversation-metadata", "probe-check"})
	if rawRes.IsFailure() && strings.Contains(rawRes.Err.Error(), "ANTIGRAVITY_LS_ADDRESS is not set") {
		t.Fatalf("agentapi failed with ANTIGRAVITY_LS_ADDRESS is not set: %v", rawRes.Err)
	}
	t.Logf("agentapi executed successfully (no LS_ADDRESS error)")
}
