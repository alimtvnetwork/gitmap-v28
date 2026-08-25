package cmd

import (
	"testing"
)

func TestComp263Success(t *testing.T) {
	input := Input263{ID: Comp263Uniqueness}
	output, err := HandleComp263(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp263Failure(t *testing.T) {
	input := Input263{ID: ""}
	output, err := HandleComp263(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp263(t *testing.T) {
	t.Run("success", TestComp263Success)
	t.Run("failure", TestComp263Failure)
}
