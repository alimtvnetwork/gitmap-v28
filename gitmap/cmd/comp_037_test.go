package cmd

import (
	"testing"
)

func TestComp037Success(t *testing.T) {
	input := Input037{ID: Comp037Uniqueness}
	output, err := HandleComp037(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp037Failure(t *testing.T) {
	input := Input037{ID: ""}
	output, err := HandleComp037(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp037(t *testing.T) {
	t.Run("success", TestComp037Success)
	t.Run("failure", TestComp037Failure)
}
