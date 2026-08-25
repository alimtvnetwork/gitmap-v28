package cmd

import (
	"testing"
)

func TestComp003(t *testing.T) {
	testInput := Input003{
		ID: "4e07408562be",
	}

	output, err := HandleComp003(testInput)
	if err != nil {
		t.Fatalf("HandleComp003 returned unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}

	invalidInput := Input003{
		ID: "",
	}

	failedOutput, failedErr := HandleComp003(invalidInput)
	if failedErr == nil {
		t.Fatalf("expected error for empty ID, got nil")
	}

	if failedOutput.Result {
		t.Fatalf("expected failedOutput Result to be false, got true")
	}
}
