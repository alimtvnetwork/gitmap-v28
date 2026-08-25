package cmd

import (
	"testing"
)

func TestComp105Success(t *testing.T) {
	input := Input105{ID: Comp105Uniqueness}
	output, err := HandleComp105(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp105Failure(t *testing.T) {
	input := Input105{ID: ""}
	output, err := HandleComp105(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp105(t *testing.T) {
	t.Run("success", TestComp105Success)
	t.Run("failure", TestComp105Failure)
}
