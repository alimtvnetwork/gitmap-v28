package cmd

import (
	"testing"
)

func TestComp110Success(t *testing.T) {
	input := Input110{ID: Comp110Uniqueness}
	output, err := HandleComp110(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp110Failure(t *testing.T) {
	input := Input110{ID: ""}
	output, err := HandleComp110(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp110(t *testing.T) {
	t.Run("success", TestComp110Success)
	t.Run("failure", TestComp110Failure)
}
