package cmd

import (
	"testing"
)

func TestComp128Success(t *testing.T) {
	input := Input128{ID: Comp128Uniqueness}
	output, err := HandleComp128(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp128Failure(t *testing.T) {
	input := Input128{ID: ""}
	output, err := HandleComp128(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp128(t *testing.T) {
	t.Run("success", TestComp128Success)
	t.Run("failure", TestComp128Failure)
}
