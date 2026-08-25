package cmd

import (
	"testing"
)

func TestComp238Success(t *testing.T) {
	input := Input238{ID: Comp238Uniqueness}
	output, err := HandleComp238(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp238Failure(t *testing.T) {
	input := Input238{ID: ""}
	output, err := HandleComp238(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp238(t *testing.T) {
	t.Run("success", TestComp238Success)
	t.Run("failure", TestComp238Failure)
}
