package cmd

import (
	"testing"
)

func TestComp083(t *testing.T) {
	in := Input083{ID: "test"}
	out, err := HandleComp083(in)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected Result true, got false")
	}

	// Test error case
	inErr := Input083{ID: ""}
	_, err = HandleComp083(inErr)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
