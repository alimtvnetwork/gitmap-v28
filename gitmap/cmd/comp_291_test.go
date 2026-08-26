package cmd

import (
	"testing"
)

func TestComp291Success(t *testing.T) {
	input := Input291{ID: Comp291Uniqueness}
	output, err := HandleComp291(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp291Failure(t *testing.T) {
	input := Input291{ID: ""}
	output, err := HandleComp291(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp291(t *testing.T) {
	t.Run("success", TestComp291Success)
	t.Run("failure", TestComp291Failure)
}
