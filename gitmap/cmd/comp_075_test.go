package cmd

import (
	"testing"
)

func TestComp075Success(t *testing.T) {
	input := Input075{ID: Comp075Uniqueness}
	output, err := HandleComp075(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp075Failure(t *testing.T) {
	input := Input075{ID: ""}
	output, err := HandleComp075(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp075(t *testing.T) {
	t.Run("success", TestComp075Success)
	t.Run("failure", TestComp075Failure)
}
