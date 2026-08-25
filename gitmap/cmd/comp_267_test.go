package cmd

import (
	"testing"
)

func TestComp267Success(t *testing.T) {
	input := Input267{ID: Comp267Uniqueness}
	output, err := HandleComp267(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp267Failure(t *testing.T) {
	input := Input267{ID: ""}
	output, err := HandleComp267(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp267(t *testing.T) {
	t.Run("success", TestComp267Success)
	t.Run("failure", TestComp267Failure)
}
