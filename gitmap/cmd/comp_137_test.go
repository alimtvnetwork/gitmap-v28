package cmd

import (
	"testing"
)

func TestComp137Success(t *testing.T) {
	input := Input137{ID: Comp137Uniqueness}
	output, err := HandleComp137(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp137Failure(t *testing.T) {
	input := Input137{ID: ""}
	output, err := HandleComp137(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp137(t *testing.T) {
	t.Run("success", TestComp137Success)
	t.Run("failure", TestComp137Failure)
}
