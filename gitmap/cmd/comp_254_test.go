package cmd

import (
	"testing"
)

func TestComp254Success(t *testing.T) {
	input := Input254{ID: Comp254Uniqueness}
	output, err := HandleComp254(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp254Failure(t *testing.T) {
	input := Input254{ID: ""}
	output, err := HandleComp254(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp254(t *testing.T) {
	t.Run("success", TestComp254Success)
	t.Run("failure", TestComp254Failure)
}
