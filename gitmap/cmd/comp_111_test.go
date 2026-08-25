package cmd

import (
	"testing"
)

func TestComp111Success(t *testing.T) {
	input := Input111{ID: Comp111Uniqueness}
	output, err := HandleComp111(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp111Failure(t *testing.T) {
	input := Input111{ID: ""}
	output, err := HandleComp111(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp111(t *testing.T) {
	t.Run("success", TestComp111Success)
	t.Run("failure", TestComp111Failure)
}
