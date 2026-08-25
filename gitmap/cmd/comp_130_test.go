package cmd

import (
	"testing"
)

func TestComp130Success(t *testing.T) {
	input := Input130{ID: Comp130Uniqueness}
	output, err := HandleComp130(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp130Failure(t *testing.T) {
	input := Input130{ID: ""}
	output, err := HandleComp130(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp130(t *testing.T) {
	t.Run("success", TestComp130Success)
	t.Run("failure", TestComp130Failure)
}
