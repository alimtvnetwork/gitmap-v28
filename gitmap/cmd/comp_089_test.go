package cmd

import (
	"testing"
)

func TestComp089Success(t *testing.T) {
	input := Input089{ID: Comp089Uniqueness}
	output, err := HandleComp089(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp089Failure(t *testing.T) {
	input := Input089{ID: ""}
	output, err := HandleComp089(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp089(t *testing.T) {
	t.Run("success", TestComp089Success)
	t.Run("failure", TestComp089Failure)
}
