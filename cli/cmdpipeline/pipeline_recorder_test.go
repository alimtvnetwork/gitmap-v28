package cmdpipeline

import (
	"testing"
	"time"
)

func TestBuildPipelineRunRecord_PopulatesDurationAndSuccess(t *testing.T) {
	now := time.Now().UTC()
	created := now.Add(-125 * time.Second).Format(time.RFC3339)
	updated := now.Format(time.RFC3339)

	payload := PipelineStatusPayload{
		Repo:       "alimtvnetwork/gitmap-v28",
		EtaSeconds: 90,
	}

	runSuccess := ghRunItem{
		DatabaseId: 1001,
		Name:       "CI",
		Status:     "completed",
		Conclusion: "success",
		HeadBranch: "main",
		HeadSha:    "abc1234",
		Url:        "https://github.com/runs/1001",
		CreatedAt:  created,
		UpdatedAt:  updated,
	}

	recSuccess := buildPipelineRunRecord(payload, runSuccess)
	if recSuccess.DurationSeconds != 125 {
		t.Fatalf("expected DurationSeconds 125, got %d", recSuccess.DurationSeconds)
	}
	if !recSuccess.IsSuccess {
		t.Fatalf("expected IsSuccess true for successful run")
	}
	if recSuccess.EtaSeconds != 0 {
		t.Fatalf("expected EtaSeconds 0 for completed run, got %d", recSuccess.EtaSeconds)
	}

	runRunning := ghRunItem{
		DatabaseId: 1002,
		Name:       "CI",
		Status:     "in_progress",
		Conclusion: "",
		CreatedAt:  created,
		UpdatedAt:  created,
	}

	recRunning := buildPipelineRunRecord(payload, runRunning)
	if recRunning.IsSuccess {
		t.Fatalf("expected IsSuccess false for in_progress run")
	}
	if recRunning.EtaSeconds != 90 {
		t.Fatalf("expected EtaSeconds 90 for running run, got %d", recRunning.EtaSeconds)
	}
}

func TestFormatEtaDisplay(t *testing.T) {
	if formatEtaDisplay(0) != "-" {
		t.Errorf("expected '-', got %s", formatEtaDisplay(0))
	}
	if formatEtaDisplay(45) != "~45s" {
		t.Errorf("expected '~45s', got %s", formatEtaDisplay(45))
	}
	if formatEtaDisplay(180) != "~180s (3m)" {
		t.Errorf("expected '~180s (3m)', got %s", formatEtaDisplay(180))
	}
	if formatEtaDisplay(125) != "~125s (2m 5s)" {
		t.Errorf("expected '~125s (2m 5s)', got %s", formatEtaDisplay(125))
	}
}
