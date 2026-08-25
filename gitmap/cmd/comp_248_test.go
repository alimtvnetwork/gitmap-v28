package cmd

import (
	"testing"
)

func TestComp248Success(t *testing.T) {
	input := Input248{ID: Comp248Uniqueness}
	output, err := HandleComp248(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp248Failure(t *testing.T) {
	input := Input248{ID: ""}
	output, err := HandleComp248(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp248(t *testing.T) {
	t.Run("success", TestComp248Success)
	t.Run("failure", TestComp248Failure)
}
