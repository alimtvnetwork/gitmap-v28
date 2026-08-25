package cmd

import (
	"testing"
)

func TestComp244Success(t *testing.T) {
	input := Input244{ID: Comp244Uniqueness}
	output, err := HandleComp244(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp244Failure(t *testing.T) {
	input := Input244{ID: ""}
	output, err := HandleComp244(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp244(t *testing.T) {
	t.Run("success", TestComp244Success)
	t.Run("failure", TestComp244Failure)
}
