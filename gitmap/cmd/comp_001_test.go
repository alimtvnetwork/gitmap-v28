package cmd

import (
	"testing"
)

func TestComp001Success(t *testing.T) {
	input := Input001{ID: Comp001Uniqueness}
	output, err := HandleComp001(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp001Failure(t *testing.T) {
	input := Input001{ID: ""}
	output, err := HandleComp001(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp001(t *testing.T) {
	t.Run("success", TestComp001Success)
	t.Run("failure", TestComp001Failure)
}
