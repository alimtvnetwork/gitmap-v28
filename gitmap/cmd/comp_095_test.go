package cmd

import (
	"testing"
)

func TestComp095Success(t *testing.T) {
	input := Input095{ID: Comp095Uniqueness}
	output, err := HandleComp095(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp095Failure(t *testing.T) {
	input := Input095{ID: ""}
	output, err := HandleComp095(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp095(t *testing.T) {
	t.Run("success", TestComp095Success)
	t.Run("failure", TestComp095Failure)
}
