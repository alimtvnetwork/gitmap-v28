package cmd

import (
	"testing"
)

func TestComp132Success(t *testing.T) {
	input := Input132{ID: Comp132Uniqueness}
	output, err := HandleComp132(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp132Failure(t *testing.T) {
	input := Input132{ID: ""}
	output, err := HandleComp132(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp132(t *testing.T) {
	t.Run("success", TestComp132Success)
	t.Run("failure", TestComp132Failure)
}