package cmd

import (
	"testing"
)

func TestComp268Success(t *testing.T) {
	input := Input268{ID: Comp268Uniqueness}
	output, err := HandleComp268(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp268Failure(t *testing.T) {
	input := Input268{ID: ""}
	output, err := HandleComp268(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp268(t *testing.T) {
	t.Run("success", TestComp268Success)
	t.Run("failure", TestComp268Failure)
}
