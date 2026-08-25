package cmd

import (
	"testing"
)

func TestComp274Success(t *testing.T) {
	input := Input274{ID: Comp274Uniqueness}
	output, err := HandleComp274(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp274Failure(t *testing.T) {
	input := Input274{ID: ""}
	output, err := HandleComp274(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp274(t *testing.T) {
	t.Run("success", TestComp274Success)
	t.Run("failure", TestComp274Failure)
}
