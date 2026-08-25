package cmd

import (
	"testing"
)

func TestComp084Success(t *testing.T) {
	input := Input084{ID: Comp084Uniqueness}
	output, err := HandleComp084(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp084Failure(t *testing.T) {
	input := Input084{ID: ""}
	output, err := HandleComp084(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp084(t *testing.T) {
	t.Run("success", TestComp084Success)
	t.Run("failure", TestComp084Failure)
}
