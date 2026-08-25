package cmd

import (
	"testing"
)

func TestComp047Success(t *testing.T) {
	input := Input047{ID: Comp047Uniqueness}
	output, err := HandleComp047(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp047Failure(t *testing.T) {
	input := Input047{ID: ""}
	output, err := HandleComp047(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp047(t *testing.T) {
	t.Run("success", TestComp047Success)
	t.Run("failure", TestComp047Failure)
}
