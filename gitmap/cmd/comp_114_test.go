package cmd

import (
	"testing"
)

func TestComp114Success(t *testing.T) {
	input := Input114{ID: Comp114Uniqueness}
	output, err := HandleComp114(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp114Failure(t *testing.T) {
	input := Input114{ID: ""}
	output, err := HandleComp114(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp114(t *testing.T) {
	t.Run("success", TestComp114Success)
	t.Run("failure", TestComp114Failure)
}
