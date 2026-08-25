package cmd

import (
	"testing"
)

func TestComp266Success(t *testing.T) {
	input := Input266{ID: Comp266Uniqueness}
	output, err := HandleComp266(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp266Failure(t *testing.T) {
	input := Input266{ID: ""}
	output, err := HandleComp266(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp266(t *testing.T) {
	t.Run("success", TestComp266Success)
	t.Run("failure", TestComp266Failure)
}
