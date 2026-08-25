package cmd

import (
	"testing"
)

func TestComp081(t *testing.T) {
	in := Input081{ID: "test"}
	out, err := HandleComp081(in)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected Result true, got false")
	}

	// Test error case
	inErr := Input081{ID: ""}
	_, err = HandleComp081(inErr)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
