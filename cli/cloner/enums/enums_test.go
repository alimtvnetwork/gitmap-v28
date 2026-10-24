package enums

import "testing"

func TestAuditActionType(t *testing.T) {
	if !AuditActionClone.IsValid() {
		t.Fatal("expected AuditActionClone to be valid")
	}
	if AuditActionClone.Marker() != "+" {
		t.Fatalf("expected '+', got %q", AuditActionClone.Marker())
	}
	if AuditActionInvalid.IsValid() {
		t.Fatal("expected AuditActionInvalid to be not valid")
	}
}

func TestGitUnmergedStatusType(t *testing.T) {
	if !StatusBothModified.IsValid() {
		t.Fatal("expected StatusBothModified to be valid")
	}
	if StatusBothModified.String() != "UU" {
		t.Fatalf("expected UU, got %s", StatusBothModified.String())
	}
}

func TestCloneDepthType(t *testing.T) {
	if !CloneDepthShallow.IsValid() {
		t.Fatal("expected shallow to be valid")
	}
	if !CloneDepthShallow.IsShallow() {
		t.Fatal("expected isShallow true")
	}
	if CloneDepthFull.IsShallow() {
		t.Fatal("expected isShallow false for full")
	}
}
