package cmd

import (
	"testing"
)

func TestComp107Success(t *testing.T) {
	input := Input107{ID: Comp107Uniqueness}
	output, err := HandleComp107(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp107Failure(t *testing.T) {
	input := Input107{ID: ""}
	output, err := HandleComp107(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp107(t *testing.T) {
	t.Run("success", TestComp107Success)
	t.Run("failure", TestComp107Failure)
}
