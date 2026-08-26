package cmd

import (
	"testing"
)

func TestComp287Success(t *testing.T) {
	input := Input287{ID: Comp287Uniqueness}
	output, err := HandleComp287(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp287Failure(t *testing.T) {
	input := Input287{ID: ""}
	output, err := HandleComp287(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp287(t *testing.T) {
	t.Run("success", TestComp287Success)
	t.Run("failure", TestComp287Failure)
}
