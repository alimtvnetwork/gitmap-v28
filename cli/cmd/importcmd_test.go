package cmd

import (
	"testing"
)

func TestRunImport_BareInvocationReturnsCleanGuidance(t *testing.T) {
	err := runImport([]string{})
	if err != nil {
		t.Errorf("expected bare runImport to return nil without throwing error, got: %v", err)
	}
}

func TestRunImport_MissingConfirmReturnsSafetyNotice(t *testing.T) {
	err := runImport([]string{"nonexistent.json"})
	if err != nil {
		t.Errorf("expected runImport without --confirm to return nil without throwing error, got: %v", err)
	}
}
