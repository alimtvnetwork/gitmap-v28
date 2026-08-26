package cmd

import (
	"testing"
)

func TestComp293Success(t *testing.T) {
	input := Input293{ID: Comp293Uniqueness}
	output, err := HandleComp293(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp293Failure(t *testing.T) {
	input := Input293{ID: ""}
	output, err := HandleComp293(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp293(t *testing.T) {
	t.Run("success", TestComp293Success)
	t.Run("failure", TestComp293Failure)
}
