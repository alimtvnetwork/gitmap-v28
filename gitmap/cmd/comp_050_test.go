package cmd

import (
	"testing"
)

func TestComp050Success(t *testing.T) {
	input := Input050{ID: Comp050Uniqueness}
	output, err := HandleComp050(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp050Failure(t *testing.T) {
	input := Input050{ID: ""}
	output, err := HandleComp050(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp050(t *testing.T) {
	t.Run("success", TestComp050Success)
	t.Run("failure", TestComp050Failure)
}
