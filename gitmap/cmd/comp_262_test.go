package cmd

import (
	"testing"
)

func TestComp262Success(t *testing.T) {
	input := Input262{ID: Comp262Uniqueness}
	output, err := HandleComp262(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp262Failure(t *testing.T) {
	input := Input262{ID: ""}
	output, err := HandleComp262(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp262(t *testing.T) {
	t.Run("success", TestComp262Success)
	t.Run("failure", TestComp262Failure)
}
