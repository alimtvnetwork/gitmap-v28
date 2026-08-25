package cmd

import (
	"testing"
)

func TestComp270Success(t *testing.T) {
	input := Input270{ID: Comp270Uniqueness}
	output, err := HandleComp270(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp270Failure(t *testing.T) {
	input := Input270{ID: ""}
	output, err := HandleComp270(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp270(t *testing.T) {
	t.Run("success", TestComp270Success)
	t.Run("failure", TestComp270Failure)
}
