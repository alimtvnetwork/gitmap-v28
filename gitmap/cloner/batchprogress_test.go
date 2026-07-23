package cloner

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestFailWithError_RecordsFailure(t *testing.T) {
	p := NewBatchProgress(3, "test", true)
	p.FailWithError("repo-a", "connection refused")

	if p.Failed() != 1 {
		t.Errorf("Failed() = %d, want 1", p.Failed())
	}
	if !p.HasFailures() {
		t.Error("HasFailures() = false, want true")
	}
	if len(p.Failures()) != 1 {
		t.Fatalf("len(Failures()) = %d, want 1", len(p.Failures()))
	}
	if p.Failures()[0].Name != "repo-a" {
		t.Errorf("Name = %q, want %q", p.Failures()[0].Name, "repo-a")
	}
	if p.Failures()[0].Error != "connection refused" {
		t.Errorf("Error = %q, want %q", p.Failures()[0].Error, "connection refused")
	}
}

func TestFailWithError_MultipleFailures(t *testing.T) {
	p := NewBatchProgress(5, "pull", true)
	p.FailWithError("repo-a", "timeout")
	p.FailWithError("repo-b", "auth failed")
	p.FailWithError("repo-c", "not found")

	if p.Failed() != 3 {
		t.Errorf("Failed() = %d, want 3", p.Failed())
	}
	if len(p.Failures()) != 3 {
		t.Errorf("len(Failures()) = %d, want 3", len(p.Failures()))
	}
}

func TestStopOnFail_StopsAfterFirstFailure(t *testing.T) {
	p := NewBatchProgress(5, "exec", true)
	p.SetStopOnFail(true)

	if p.Stopped() {
		t.Error("Stopped() = true before any failure")
	}

	p.FailWithError("repo-x", "exit status 1")

	if !p.Stopped() {
		t.Error("Stopped() = false after FailWithError with stop-on-fail")
	}
}

func TestStopOnFail_Disabled(t *testing.T) {
	p := NewBatchProgress(5, "exec", true)

	p.FailWithError("repo-x", "exit status 1")

	if p.Stopped() {
		t.Error("Stopped() = true without SetStopOnFail")
	}
}

func TestExitCodeForBatch_ZeroOnSuccess(t *testing.T) {
	p := NewBatchProgress(2, "pull", true)
	p.Succeed()
	p.Succeed()

	if code := p.ExitCodeForBatch(); code != 0 {
		t.Errorf("ExitCodeForBatch() = %d, want 0", code)
	}
}

func TestExitCodeForBatch_PartialFailure(t *testing.T) {
	p := NewBatchProgress(3, "pull", true)
	p.Succeed()
	p.FailWithError("repo-b", "error")
	p.Succeed()

	if code := p.ExitCodeForBatch(); code != constants.ExitPartialFailure {
		t.Errorf("ExitCodeForBatch() = %d, want %d", code, constants.ExitPartialFailure)
	}
}

func TestExitCodeForBatch_AllFailed(t *testing.T) {
	p := NewBatchProgress(2, "exec", true)
	p.FailWithError("a", "err1")
	p.FailWithError("b", "err2")

	if code := p.ExitCodeForBatch(); code != constants.ExitPartialFailure {
		t.Errorf("ExitCodeForBatch() = %d, want %d", code, constants.ExitPartialFailure)
	}
}

func TestNoFailures_Defaults(t *testing.T) {
	p := NewBatchProgress(1, "status", true)

	if p.HasFailures() {
		t.Error("HasFailures() = true on fresh progress")
	}
	if p.Failed() != 0 {
		t.Errorf("Failed() = %d, want 0", p.Failed())
	}
	if p.Stopped() {
		t.Error("Stopped() = true on fresh progress")
	}
	if code := p.ExitCodeForBatch(); code != 0 {
		t.Errorf("ExitCodeForBatch() = %d, want 0", code)
	}
}

func TestMixedOperations_Counters(t *testing.T) {
	p := NewBatchProgress(6, "pull", true)
	p.BeginItem("a")
	p.Succeed()
	p.BeginItem("b")
	p.FailWithError("b", "err")
	p.BeginItem("c")
	p.Skip()
	p.BeginItem("d")
	p.Succeed()
	p.BeginItem("e")
	p.Fail()
	p.BeginItem("f")
	p.Succeed()

	if p.Succeeded() != 3 {
		t.Errorf("Succeeded() = %d, want 3", p.Succeeded())
	}
	if p.Failed() != 2 {
		t.Errorf("Failed() = %d, want 2", p.Failed())
	}
	if p.Skipped() != 1 {
		t.Errorf("Skipped() = %d, want 1", p.Skipped())
	}
	if len(p.Failures()) != 1 {
		t.Errorf("len(Failures()) = %d, want 1 (only FailWithError records)", len(p.Failures()))
	}
}

func TestPrintFailureReport_NoFailures(t *testing.T) {
	p := NewBatchProgress(1, "test", true)
	// Should not panic when called with no failures
	p.PrintFailureReport()
}

func TestPrintFailureReport_WithFailures(t *testing.T) {
	p := NewBatchProgress(2, "exec", false)
	p.FailWithError("repo-x", "permission denied")
	p.FailWithError("repo-y", "not a git repository")
	// Should not panic; output goes to stderr
	p.PrintFailureReport()
}
