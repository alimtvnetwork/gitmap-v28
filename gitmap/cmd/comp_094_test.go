package cmd

import (
	"testing"
)

func TestComp094Success(t *testing.T) {
	input := Input094{ID: Comp094Uniqueness}
	output, err := HandleComp094(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp094Failure(t *testing.T) {
	input := Input094{ID: ""}
	output, err := HandleComp094(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp094(t *testing.T) {
	t.Run("success", TestComp094Success)
	t.Run("failure", TestComp094Failure)
}
