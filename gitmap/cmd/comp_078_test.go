package cmd

import (
	"testing"
)

func TestComp078Success(t *testing.T) {
	input := Input078{ID: Comp078Uniqueness}
	output, err := HandleComp078(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp078Failure(t *testing.T) {
	input := Input078{ID: ""}
	output, err := HandleComp078(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp078(t *testing.T) {
	t.Run("success", TestComp078Success)
	t.Run("failure", TestComp078Failure)
}
