package cmd

import (
	"testing"
)

func TestComp187Success(t *testing.T) {
	input := Input187{ID: Comp187Uniqueness}
	output, err := HandleComp187(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp187Failure(t *testing.T) {
	input := Input187{ID: "invalid"}
	output, err := HandleComp187(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp187(t *testing.T) {
	t.Run("success", TestComp187Success)
	t.Run("failure", TestComp187Failure)
}
