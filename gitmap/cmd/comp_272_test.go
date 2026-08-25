package cmd

import (
	"testing"
)

func TestComp272Success(t *testing.T) {
	input := Input272{ID: Comp272Uniqueness}
	output, err := HandleComp272(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp272Failure(t *testing.T) {
	input := Input272{ID: ""}
	output, err := HandleComp272(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp272(t *testing.T) {
	t.Run("success", TestComp272Success)
	t.Run("failure", TestComp272Failure)
}
