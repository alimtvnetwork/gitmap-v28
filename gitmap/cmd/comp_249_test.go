package cmd

import (
	"testing"
)

func TestComp249Success(t *testing.T) {
	input := Input249{ID: Comp249Uniqueness}
	output, err := HandleComp249(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp249Failure(t *testing.T) {
	input := Input249{ID: ""}
	output, err := HandleComp249(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp249(t *testing.T) {
	t.Run("success", TestComp249Success)
	t.Run("failure", TestComp249Failure)
}
