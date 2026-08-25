package cmd

import (
	"testing"
)

func TestComp186Success(t *testing.T) {
	input := Input186{ID: Comp186Uniqueness}
	output, err := HandleComp186(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp186Failure(t *testing.T) {
	input := Input186{ID: ""}
	output, err := HandleComp186(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp186(t *testing.T) {
	t.Run("success", TestComp186Success)
	t.Run("failure", TestComp186Failure)
}
