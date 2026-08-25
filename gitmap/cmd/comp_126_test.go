package cmd

import (
	"testing"
)

func TestComp126Success(t *testing.T) {
	input := Input126{ID: Comp126Uniqueness}
	output, err := HandleComp126(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp126Failure(t *testing.T) {
	input := Input126{ID: ""}
	output, err := HandleComp126(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp126(t *testing.T) {
	t.Run("success", TestComp126Success)
	t.Run("failure", TestComp126Failure)
}
