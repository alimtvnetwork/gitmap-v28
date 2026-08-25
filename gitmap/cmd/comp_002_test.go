package cmd

import (
	"testing"
)

func TestComp002(t *testing.T) {
	input := Input002{ID: comp002BoundID}
	output, err := HandleComp002(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}

	failInput := Input002{ID: ""}
	failOutput, failErr := HandleComp002(failInput)
	if failErr == nil {
		t.Fatalf("expected error for empty ID, got nil")
	}

	if failOutput.Result {
		t.Fatalf("expected Result to be false on failure, got true")
	}
}
