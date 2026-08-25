package cmd

import (
	"testing"
)

func TestComp080Success(t *testing.T) {
	input := Input080{ID: Comp080Uniqueness}
	output, err := HandleComp080(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp080Failure(t *testing.T) {
	input := Input080{ID: ""}
	output, err := HandleComp080(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp080(t *testing.T) {
	t.Run("success", TestComp080Success)
	t.Run("failure", TestComp080Failure)
}
