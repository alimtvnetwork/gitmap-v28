package cmd

import (
	"testing"
)

func TestComp042Success(t *testing.T) {
	input := Input042{ID: Comp042Uniqueness}
	output, err := HandleComp042(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp042Failure(t *testing.T) {
	input := Input042{ID: ""}
	output, err := HandleComp042(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp042(t *testing.T) {
	t.Run("success", TestComp042Success)
	t.Run("failure", TestComp042Failure)
}
