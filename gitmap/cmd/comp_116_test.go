package cmd

import (
	"testing"
)

func TestComp116Success(t *testing.T) {
	input := Input116{ID: Comp116Uniqueness}
	output, err := HandleComp116(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp116Failure(t *testing.T) {
	input := Input116{ID: ""}
	output, err := HandleComp116(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp116(t *testing.T) {
	t.Run("success", TestComp116Success)
	t.Run("failure", TestComp116Failure)
}
