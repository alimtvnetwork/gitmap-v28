package cmd

import (
	"testing"
)

func TestComp295Success(t *testing.T) {
	input := Input295{ID: Comp295Uniqueness}
	output, err := HandleComp295(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp295Failure(t *testing.T) {
	input := Input295{ID: ""}
	output, err := HandleComp295(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp295(t *testing.T) {
	t.Run("success", TestComp295Success)
	t.Run("failure", TestComp295Failure)
}
