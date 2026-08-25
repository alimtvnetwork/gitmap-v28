package cmd

import (
	"testing"
)

func TestComp063Success(t *testing.T) {
	input := Input063{ID: comp063Uniqueness}
	output, err := HandleComp063(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp063Failure(t *testing.T) {
	input := Input063{ID: ""}
	output, err := HandleComp063(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp063(t *testing.T) {
	t.Run("success", TestComp063Success)
	t.Run("failure", TestComp063Failure)
}
