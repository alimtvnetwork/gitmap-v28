package cmd

import (
	"testing"
)

func TestComp294Success(t *testing.T) {
	input := Input294{ID: Comp294Uniqueness}
	output, err := HandleComp294(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp294Failure(t *testing.T) {
	input := Input294{ID: ""}
	output, err := HandleComp294(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp294(t *testing.T) {
	t.Run("success", TestComp294Success)
	t.Run("failure", TestComp294Failure)
}
