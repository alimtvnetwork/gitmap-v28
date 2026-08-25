package cmd

import (
	"testing"
)

func TestComp255Success(t *testing.T) {
	input := Input255{ID: Comp255Uniqueness}
	output, err := HandleComp255(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp255Failure(t *testing.T) {
	input := Input255{ID: ""}
	output, err := HandleComp255(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp255(t *testing.T) {
	t.Run("success", TestComp255Success)
	t.Run("failure", TestComp255Failure)
}
