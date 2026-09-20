package cmdagy

import (
	"testing"
)

func TestExecuteAgyPing(t *testing.T) {
	report := ExecuteAgyPing("")
	hasTimestamp := len(report.Timestamp) > 0
	if hasTimestamp == false {
		t.Errorf("expected non-empty timestamp, got %s", report.Timestamp)
	}

	hasWorkspace := len(report.Workspace.TargetWorkspace) > 0
	if hasWorkspace == false {
		t.Errorf("expected target workspace, got %s", report.Workspace.TargetWorkspace)
	}
}

func TestRenderPingReport(t *testing.T) {
	report := ExecuteAgyPing(".")
	renderPingReport(report)
}

func TestRenderPingJSON(t *testing.T) {
	report := ExecuteAgyPing(".")
	err := renderPingJSON(report)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("expected nil error from renderPingJSON, got %v", err)
	}
}

func TestCheckIDEExecutable(t *testing.T) {
	chk := checkIDEExecutable()
	hasResult := chk.IsFound || len(chk.Error) > 0
	if hasResult == false {
		t.Errorf("expected either IsFound or Error set, got %+v", chk)
	}
}

func TestCheckIDEProcess(t *testing.T) {
	chk := checkIDEProcess()
	hasResult := chk.IsRunning || len(chk.Error) > 0
	if hasResult == false {
		t.Errorf("expected either IsRunning or Error set, got %+v", chk)
	}
}

func TestCheckFilesystemHealth(t *testing.T) {
	chk := checkFilesystemHealth()
	hasBrain := len(chk.BrainDir) > 0
	if hasBrain == false {
		t.Errorf("expected BrainDir, got %+v", chk)
	}
}

func TestCheckPromptQueue(t *testing.T) {
	chk := checkPromptQueue()
	hasValidQueueState := chk.HasQueue || len(chk.Error) > 0
	if hasValidQueueState == false {
		t.Errorf("expected queue state, got %+v", chk)
	}
}

func TestComputePingHealth(t *testing.T) {
	healthyReport := AgyPingReport{
		Filesystem: AgyFilesystemHealth{IsAccessible: true},
		Executable: AgyPingExecutableCheck{IsFound: true},
	}
	isHealthy := computePingHealth(healthyReport)
	if isHealthy == false {
		t.Errorf("expected healthy report to be true")
	}
}
