package cmd

import (
	"testing"
)

func TestComp133Success(t *testing.T) {
	input := Input133{ID: Comp133Uniqueness}
	output, err := HandleComp133(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp133Failure(t *testing.T) {
	input := Input133{ID: ""}
	output, err := HandleComp133(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp133(t *testing.T) {
	t.Run("success", TestComp133Success)
	t.Run("failure", TestComp133Failure)
}
