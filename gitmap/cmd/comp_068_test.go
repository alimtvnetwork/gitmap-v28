package cmd

import (
	"testing"
)

func TestComp068(t *testing.T) {
	input := Input068{ID: comp068BoundID}
	output, err := HandleComp068(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}

	failInput := Input068{ID: ""}
	failOutput, failErr := HandleComp068(failInput)
	if failErr == nil {
		t.Fatalf("expected error for empty ID, got nil")
	}

	if failOutput.Result {
		t.Fatalf("expected Result to be false on failure, got true")
	}
}
