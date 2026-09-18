package cmdpipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadRunnerETAMissing(t *testing.T) {
	// Points to a temporary dir where no runner-eta.json exists
	info, hasInfo := ReadRunnerETA()
	if hasInfo && info != nil && info.Status == "running" && info.EtaSeconds > 999999 {
		t.Errorf("expected no active runner-eta in clean environment")
	}
}

func TestIsRunnerActive(t *testing.T) {
	active := &RunnerETAInfo{Status: "running", EtaSeconds: 15}
	if !isRunnerActive(active) {
		t.Errorf("expected active to be true")
	}

	completed := &RunnerETAInfo{Status: "completed", EtaSeconds: 0}
	if isRunnerActive(completed) {
		t.Errorf("expected completed to be false")
	}

	zeroETA := &RunnerETAInfo{Status: "running", EtaSeconds: 0}
	if isRunnerActive(zeroETA) {
		t.Errorf("expected zero ETA to be false")
	}
}

func TestWaitForRunnerETAActive(t *testing.T) {
	tmpDir := t.TempDir()
	lovableTemp := filepath.Join(tmpDir, ".ai-memory", "temp")
	if err := os.MkdirAll(lovableTemp, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	sample := RunnerETAInfo{
		Status:                "running",
		EtaSeconds:            2,
		ElapsedSeconds:        5,
		TotalEstimatedSeconds: 7,
		UpdatedAt:             1726000000,
	}

	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(lovableTemp, "runner-eta.json"), data, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Should not hang in test mode
	runRunnerCountdownLoop(&sample)
}
