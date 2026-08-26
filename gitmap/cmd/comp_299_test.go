package cmd

import (
	"testing"
)

func TestComp299Success(t *testing.T) {
	input := Input299{ID: Comp299Uniqueness}
	output, err := HandleComp299(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp299Failure(t *testing.T) {
	input := Input299{ID: ""}
	output, err := HandleComp299(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp299(t *testing.T) {
	t.Run("success", TestComp299Success)
	t.Run("failure", TestComp299Failure)
}
