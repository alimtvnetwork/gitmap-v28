package cmd

import (
	"testing"
)

func TestComp074Success(t *testing.T) {
	input := Input074{ID: Comp074Uniqueness}
	output, err := HandleComp074(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp074Failure(t *testing.T) {
	input := Input074{ID: ""}
	output, err := HandleComp074(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp074(t *testing.T) {
	t.Run("success", TestComp074Success)
	t.Run("failure", TestComp074Failure)
}
