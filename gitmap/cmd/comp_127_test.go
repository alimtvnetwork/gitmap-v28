package cmd

import (
	"testing"
)

func TestComp127Success(t *testing.T) {
	input := Input127{ID: Comp127Uniqueness}
	output, err := HandleComp127(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp127Failure(t *testing.T) {
	input := Input127{ID: ""}
	output, err := HandleComp127(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp127(t *testing.T) {
	t.Run("success", TestComp127Success)
	t.Run("failure", TestComp127Failure)
}
