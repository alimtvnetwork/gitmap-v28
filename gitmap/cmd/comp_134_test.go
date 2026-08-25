package cmd

import (
	"testing"
)

func TestComp134Success(t *testing.T) {
	input := Input134{ID: Comp134Uniqueness}
	output, err := HandleComp134(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp134Failure(t *testing.T) {
	input := Input134{ID: ""}
	output, err := HandleComp134(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp134(t *testing.T) {
	t.Run("success", TestComp134Success)
	t.Run("failure", TestComp134Failure)
}
