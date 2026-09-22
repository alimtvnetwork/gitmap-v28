//go:build e2e

package cmdagy

import (
	"os"
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
