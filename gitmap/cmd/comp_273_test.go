package cmd

import (
	"testing"
)

func TestComp273Success(t *testing.T) {
	input := Input273{ID: Comp273Uniqueness}
	output, err := HandleComp273(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp273Failure(t *testing.T) {
	input := Input273{ID: ""}
	output, err := HandleComp273(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp273(t *testing.T) {
	t.Run("success", TestComp273Success)
	t.Run("failure", TestComp273Failure)
}
