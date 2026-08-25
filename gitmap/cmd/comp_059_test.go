package cmd

import (
	"testing"
)

func TestComp059Success(t *testing.T) {
	input := Input059{ID: Comp059Uniqueness}
	output, err := HandleComp059(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp059Failure(t *testing.T) {
	input := Input059{ID: ""}
	output, err := HandleComp059(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp059(t *testing.T) {
	t.Run("success", TestComp059Success)
	t.Run("failure", TestComp059Failure)
}
