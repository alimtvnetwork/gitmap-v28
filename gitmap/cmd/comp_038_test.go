package cmd

import (
	"testing"
)

func TestComp038Success(t *testing.T) {
	input := Input038{ID: Comp038Uniqueness}
	output, err := HandleComp038(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp038Failure(t *testing.T) {
	input := Input038{ID: ""}
	output, err := HandleComp038(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp038(t *testing.T) {
	t.Run("success", TestComp038Success)
	t.Run("failure", TestComp038Failure)
}
