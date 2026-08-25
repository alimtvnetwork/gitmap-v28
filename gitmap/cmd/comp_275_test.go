package cmd

import (
	"testing"
)

func TestComp275Success(t *testing.T) {
	input := Input275{ID: Comp275Uniqueness}
	output, err := HandleComp275(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp275Failure(t *testing.T) {
	input := Input275{ID: ""}
	output, err := HandleComp275(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp275(t *testing.T) {
	t.Run("success", TestComp275Success)
	t.Run("failure", TestComp275Failure)
}
