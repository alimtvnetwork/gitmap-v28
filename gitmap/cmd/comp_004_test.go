package cmd

import (
	"testing"
)

func TestComp004(t *testing.T) {
	testInput := Input004{
		ID: comp004BoundID,
	}

	output, err := HandleComp004(testInput)
	if err != nil {
		t.Fatalf("HandleComp004 returned unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}

	invalidInput := Input004{
		ID: "",
	}

	failedOutput, failedErr := HandleComp004(invalidInput)
	if failedErr == nil {
		t.Fatalf("expected error for empty ID, got nil")
	}

	if failedOutput.Result {
		t.Fatalf("expected failedOutput Result to be false, got true")
	}
}
