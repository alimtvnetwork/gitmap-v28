package cmd

import (
	"testing"
)

func TestComp289Success(t *testing.T) {
	input := Input289{ID: Comp289Uniqueness}
	output, err := HandleComp289(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp289Failure(t *testing.T) {
	input := Input289{ID: ""}
	output, err := HandleComp289(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp289(t *testing.T) {
	t.Run("success", TestComp289Success)
	t.Run("failure", TestComp289Failure)
}
