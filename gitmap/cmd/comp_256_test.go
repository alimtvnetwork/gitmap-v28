package cmd

import (
	"testing"
)

func TestComp256Success(t *testing.T) {
	input := Input256{ID: Comp256Uniqueness}
	output, err := HandleComp256(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp256Failure(t *testing.T) {
	input := Input256{ID: ""}
	output, err := HandleComp256(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp256(t *testing.T) {
	t.Run("success", TestComp256Success)
	t.Run("failure", TestComp256Failure)
}
