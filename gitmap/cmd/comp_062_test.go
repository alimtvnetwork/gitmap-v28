package cmd

import (
	"testing"
)

func TestComp062Success(t *testing.T) {
	input := Input062{ID: Comp062Uniqueness}
	output, err := HandleComp062(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp062Failure(t *testing.T) {
	input := Input062{ID: ""}
	output, err := HandleComp062(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp062(t *testing.T) {
	t.Run("success", TestComp062Success)
	t.Run("failure", TestComp062Failure)
}
