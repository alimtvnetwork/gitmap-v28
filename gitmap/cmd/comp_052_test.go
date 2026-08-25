package cmd

import (
	"testing"
)

func TestComp052Success(t *testing.T) {
	input := Input052{ID: Comp052Uniqueness}
	output, err := HandleComp052(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp052Failure(t *testing.T) {
	input := Input052{ID: ""}
	output, err := HandleComp052(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp052(t *testing.T) {
	t.Run("success", TestComp052Success)
	t.Run("failure", TestComp052Failure)
}
