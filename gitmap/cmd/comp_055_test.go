package cmd

import (
	"testing"
)

func TestComp055(t *testing.T) {
	in := Input055{ID: "test-id"}
	out, err := HandleComp055(in)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected Result to be true, got false")
	}

	// Test empty ID failure just in case we hit the apperror path
	inEmpty := Input055{ID: ""}
	_, err = HandleComp055(inEmpty)
	if err == nil {
		t.Errorf("Expected an error for empty ID, got nil")
	}
}
