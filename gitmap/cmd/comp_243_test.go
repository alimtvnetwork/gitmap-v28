package cmd

import (
	"testing"
)

func TestComp243Success(t *testing.T) {
	input := Input243{ID: Comp243Uniqueness}
	output, err := HandleComp243(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp243Failure(t *testing.T) {
	input := Input243{ID: ""}
	output, err := HandleComp243(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp243(t *testing.T) {
	t.Run("success", TestComp243Success)
	t.Run("failure", TestComp243Failure)
}
