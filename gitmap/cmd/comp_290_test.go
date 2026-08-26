package cmd

import (
	"testing"
)

func TestComp290Success(t *testing.T) {
	input := Input290{ID: Comp290Uniqueness}
	output, err := HandleComp290(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp290Failure(t *testing.T) {
	input := Input290{ID: ""}
	output, err := HandleComp290(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp290(t *testing.T) {
	t.Run("success", TestComp290Success)
	t.Run("failure", TestComp290Failure)
}
