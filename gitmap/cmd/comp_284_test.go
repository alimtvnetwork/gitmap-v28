package cmd

import (
	"testing"
)

func TestComp284Success(t *testing.T) {
	input := Input284{ID: Comp284Uniqueness}
	output, err := HandleComp284(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp284Failure(t *testing.T) {
	input := Input284{ID: ""}
	output, err := HandleComp284(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp284(t *testing.T) {
	t.Run("success", TestComp284Success)
	t.Run("failure", TestComp284Failure)
}
