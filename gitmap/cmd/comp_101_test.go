package cmd

import (
	"testing"
)

func TestComp101Success(t *testing.T) {
	input := Input101{ID: Comp101Uniqueness}
	output, err := HandleComp101(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp101Failure(t *testing.T) {
	input := Input101{ID: ""}
	output, err := HandleComp101(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp101(t *testing.T) {
	t.Run("success", TestComp101Success)
	t.Run("failure", TestComp101Failure)
}
