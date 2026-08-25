package cmd

import (
	"testing"
)

func TestComp252Success(t *testing.T) {
	input := Input252{ID: Comp252Uniqueness}
	output, err := HandleComp252(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp252Failure(t *testing.T) {
	input := Input252{ID: ""}
	output, err := HandleComp252(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp252(t *testing.T) {
	t.Run("success", TestComp252Success)
	t.Run("failure", TestComp252Failure)
}
