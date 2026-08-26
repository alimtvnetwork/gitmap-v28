package cmd

import (
	"testing"
)

func TestComp278Success(t *testing.T) {
	input := Input278{ID: Comp278Uniqueness}
	output, err := HandleComp278(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp278Failure(t *testing.T) {
	input := Input278{ID: ""}
	output, err := HandleComp278(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp278(t *testing.T) {
	t.Run("success", TestComp278Success)
	t.Run("failure", TestComp278Failure)
}
