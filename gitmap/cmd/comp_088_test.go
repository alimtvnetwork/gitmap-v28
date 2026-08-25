package cmd

import (
	"testing"
)

func TestComp088Success(t *testing.T) {
	input := Input088{ID: Comp088Uniqueness}
	output, err := HandleComp088(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp088Failure(t *testing.T) {
	input := Input088{ID: ""}
	output, err := HandleComp088(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp088(t *testing.T) {
	t.Run("success", TestComp088Success)
	t.Run("failure", TestComp088Failure)
}
