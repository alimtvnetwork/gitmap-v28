package cmd

import (
	"testing"
)

func TestComp198Success(t *testing.T) {
	input := Input198{ID: Comp198Uniqueness}
	output, err := HandleComp198(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp198Failure(t *testing.T) {
	input := Input198{ID: ""}
	output, err := HandleComp198(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp198(t *testing.T) {
	t.Run("success", TestComp198Success)
	t.Run("failure", TestComp198Failure)
}
