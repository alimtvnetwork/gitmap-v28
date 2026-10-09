package cmdimport

import (
	"testing"
)

func TestRunImport_BareInvocationReturnsCleanGuidance(t *testing.T) {
	err := RunImport([]string{})
	if err != nil {
		t.Errorf("expected bare RunImport to return nil without throwing error, got: %v", err)
	}
}

func TestRunImport_MissingConfirmReturnsSafetyNotice(t *testing.T) {
	err := RunImport([]string{"nonexistent.json"})
	if err != nil {
		t.Errorf("expected RunImport without --confirm to return nil without throwing error, got: %v", err)
	}
}
