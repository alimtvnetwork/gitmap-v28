package cmd

import (
	"testing"
)

func TestComp298Success(t *testing.T) {
	input := Input298{ID: Comp298Uniqueness}
	output, err := HandleComp298(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp298Failure(t *testing.T) {
	input := Input298{ID: ""}
	output, err := HandleComp298(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp298(t *testing.T) {
	t.Run("success", TestComp298Success)
	t.Run("failure", TestComp298Failure)
}
