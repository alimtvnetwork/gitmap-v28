package cmd

import (
	"testing"
)

func TestComp282Success(t *testing.T) {
	input := Input282{ID: Comp282Uniqueness}
	output, err := HandleComp282(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp282Failure(t *testing.T) {
	input := Input282{ID: ""}
	output, err := HandleComp282(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp282(t *testing.T) {
	t.Run("success", TestComp282Success)
	t.Run("failure", TestComp282Failure)
}