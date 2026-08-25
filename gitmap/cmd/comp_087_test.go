package cmd

import (
	"testing"
)

func TestComp087Success(t *testing.T) {
	input := Input087{ID: Comp087Uniqueness}
	output, err := HandleComp087(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp087Failure(t *testing.T) {
	input := Input087{ID: ""}
	output, err := HandleComp087(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp087(t *testing.T) {
	t.Run("success", TestComp087Success)
	t.Run("failure", TestComp087Failure)
}
