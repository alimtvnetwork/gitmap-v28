package cmd

import (
	"testing"
)

func TestComp082(t *testing.T) {
	in := Input082{ID: "test"}
	out, err := HandleComp082(in)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected Result true, got false")
	}

	// Test error case
	inErr := Input082{ID: ""}
	_, err = HandleComp082(inErr)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
