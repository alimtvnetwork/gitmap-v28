package cmd

import (
	"testing"
)

func TestComp283Success(t *testing.T) {
	input := Input283{ID: Comp283Uniqueness}
	output, err := HandleComp283(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp283Failure(t *testing.T) {
	input := Input283{ID: ""}
	output, err := HandleComp283(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp283(t *testing.T) {
	t.Run("success", TestComp283Success)
	t.Run("failure", TestComp283Failure)
}
