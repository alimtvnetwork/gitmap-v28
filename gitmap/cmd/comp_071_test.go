package cmd

import (
	"testing"
)

func TestComp071Success(t *testing.T) {
	input := Input071{ID: Comp071Uniqueness}
	output, err := HandleComp071(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp071Failure(t *testing.T) {
	input := Input071{ID: ""}
	output, err := HandleComp071(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp071(t *testing.T) {
	t.Run("success", TestComp071Success)
	t.Run("failure", TestComp071Failure)
}
