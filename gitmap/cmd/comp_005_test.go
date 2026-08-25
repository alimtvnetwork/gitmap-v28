package cmd

import (
	"testing"
)

func TestComp005(t *testing.T) {
	testInput := Input005{
		ID: comp005BoundID,
	}

	output, err := HandleComp005(testInput)
	if err != nil {
		t.Fatalf("HandleComp005 returned unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}

	invalidInput := Input005{
		ID: "",
	}

	failedOutput, failedErr := HandleComp005(invalidInput)
	if failedErr == nil {
		t.Fatalf("expected error for empty ID, got nil")
	}

	if failedOutput.Result {
		t.Fatalf("expected failedOutput Result to be false, got true")
	}
}
