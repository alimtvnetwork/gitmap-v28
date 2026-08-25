package cmd

import (
	"testing"
)

func TestComp070Success(t *testing.T) {
	input := Input070{ID: comp070Uniqueness}
	output, err := HandleComp070(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp070Failure(t *testing.T) {
	input := Input070{ID: ""}
	output, err := HandleComp070(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp070(t *testing.T) {
	t.Run("success", TestComp070Success)
	t.Run("failure", TestComp070Failure)
}
