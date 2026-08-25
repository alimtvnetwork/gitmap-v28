package cmd

import (
	"testing"
)

func TestComp241Success(t *testing.T) {
	input := Input241{ID: Comp241Uniqueness}
	output, err := HandleComp241(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp241Failure(t *testing.T) {
	input := Input241{ID: ""}
	output, err := HandleComp241(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp241(t *testing.T) {
	t.Run("success", TestComp241Success)
	t.Run("failure", TestComp241Failure)
}
