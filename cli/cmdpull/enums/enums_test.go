package enums

import (
	"testing"
)

func TestPullStepTypeMethods(t *testing.T) {
	if PullStepInvalid.IsValid() {
		t.Errorf("expected PullStepInvalid to be invalid")
	}

	if !PullStepPending.IsValid() {
		t.Errorf("expected PullStepPending to be valid")
	}

	if !PullStepSkipped.IsValid() {
		t.Errorf("expected PullStepSkipped to be valid")
	}

	if PullStepType(99).IsValid() {
		t.Errorf("expected PullStepType(99) to be invalid")
	}

	if str := PullStepUpToDate.String(); str != "UP_TO_DATE" {
		t.Errorf("expected UP_TO_DATE, got %s", str)
	}

	if str := PullStepFastForward.String(); str != "FAST_FORWARD" {
		t.Errorf("expected FAST_FORWARD, got %s", str)
	}

	if str := PullStepType(99).String(); str != "INVALID" {
		t.Errorf("expected INVALID for out-of-range, got %s", str)
	}

	if badge := PullStepUpToDate.Badge(); badge != "✔" {
		t.Errorf("expected ✔, got %s", badge)
	}

	if badge := PullStepType(99).Badge(); badge != "" {
		t.Errorf("expected empty badge for out-of-range, got %s", badge)
	}

	if !PullStepUpToDate.IsComplete() {
		t.Errorf("expected PullStepUpToDate to be complete")
	}

	if !PullStepFastForward.IsComplete() {
		t.Errorf("expected PullStepFastForward to be complete")
	}

	if PullStepFetching.IsComplete() {
		t.Errorf("expected PullStepFetching to not be complete")
	}

	if !PullStepError.HasFailed() {
		t.Errorf("expected PullStepError to have failed")
	}

	if !PullStepConflict.HasFailed() {
		t.Errorf("expected PullStepConflict to have failed")
	}

	if PullStepUpToDate.HasFailed() {
		t.Errorf("expected PullStepUpToDate to not have failed")
	}
}

func TestPullOutcomeStatusTypeMethods(t *testing.T) {
	if !StatusUpToDate.IsValid() {
		t.Errorf("expected StatusUpToDate to be valid")
	}

	if !StatusClean.IsValid() {
		t.Errorf("expected StatusClean to be valid")
	}

	if PullOutcomeStatusType("unknown").IsValid() {
		t.Errorf("expected unknown status to be invalid")
	}

	if str := StatusMerged.String(); str != "auto-merged" {
		t.Errorf("expected auto-merged, got %s", str)
	}

	if !StatusDirty.IsDirty() {
		t.Errorf("expected StatusDirty.IsDirty() to be true")
	}

	if StatusClean.IsDirty() {
		t.Errorf("expected StatusClean.IsDirty() to be false")
	}

	if !StatusClean.IsClean() {
		t.Errorf("expected StatusClean.IsClean() to be true")
	}

	if StatusDirty.IsClean() {
		t.Errorf("expected StatusDirty.IsClean() to be false")
	}
}

func TestConflictSideTypeMethods(t *testing.T) {
	if !ConflictSideTheirs.IsValid() {
		t.Errorf("expected ConflictSideTheirs to be valid")
	}

	if !ConflictSideOurs.IsValid() {
		t.Errorf("expected ConflictSideOurs to be valid")
	}

	if ConflictSideType("other").IsValid() {
		t.Errorf("expected other to be invalid")
	}

	if str := ConflictSideTheirs.String(); str != "theirs" {
		t.Errorf("expected theirs, got %s", str)
	}

	if str := ConflictSideOurs.String(); str != "ours" {
		t.Errorf("expected ours, got %s", str)
	}
}
