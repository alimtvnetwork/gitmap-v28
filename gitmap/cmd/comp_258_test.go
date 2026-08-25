package cmd

import (
	"testing"
)

func TestComp258Success(t *testing.T) {
	input := Input258{ID: Comp258Uniqueness}
	output, err := HandleComp258(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp258Failure(t *testing.T) {
	input := Input258{ID: ""}
	output, err := HandleComp258(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp258(t *testing.T) {
	t.Run("success", TestComp258Success)
	t.Run("failure", TestComp258Failure)
}
