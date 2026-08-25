package cmd

import (
	"testing"
)

func TestComp046Success(t *testing.T) {
	input := Input046{ID: Comp046Uniqueness}
	output, err := HandleComp046(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp046Failure(t *testing.T) {
	input := Input046{ID: ""}
	output, err := HandleComp046(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp046(t *testing.T) {
	t.Run("success", TestComp046Success)
	t.Run("failure", TestComp046Failure)
}
