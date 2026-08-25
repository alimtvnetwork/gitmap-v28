package cmd

import (
	"testing"
)

func TestComp240Success(t *testing.T) {
	input := Input240{ID: Comp240Uniqueness}
	output, err := HandleComp240(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp240Failure(t *testing.T) {
	input := Input240{ID: ""}
	output, err := HandleComp240(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp240(t *testing.T) {
	t.Run("success", TestComp240Success)
	t.Run("failure", TestComp240Failure)
}
