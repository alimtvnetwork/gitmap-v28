package cmd

import (
	"testing"
)

func TestComp281Success(t *testing.T) {
	input := Input281{ID: Comp281Uniqueness}
	output, err := HandleComp281(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp281Failure(t *testing.T) {
	input := Input281{ID: ""}
	output, err := HandleComp281(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp281(t *testing.T) {
	t.Run("success", TestComp281Success)
	t.Run("failure", TestComp281Failure)
}
