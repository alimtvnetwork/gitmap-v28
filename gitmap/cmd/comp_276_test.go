package cmd

import (
	"testing"
)

func TestComp276Success(t *testing.T) {
	input := Input276{ID: Comp276Uniqueness}
	output, err := HandleComp276(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp276Failure(t *testing.T) {
	input := Input276{ID: ""}
	output, err := HandleComp276(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp276(t *testing.T) {
	t.Run("success", TestComp276Success)
	t.Run("failure", TestComp276Failure)
}
