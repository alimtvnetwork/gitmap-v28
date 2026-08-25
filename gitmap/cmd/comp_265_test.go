package cmd

import (
	"testing"
)

func TestComp265Success(t *testing.T) {
	input := Input265{ID: Comp265Uniqueness}
	output, err := HandleComp265(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp265Failure(t *testing.T) {
	input := Input265{ID: ""}
	output, err := HandleComp265(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp265(t *testing.T) {
	t.Run("success", TestComp265Success)
	t.Run("failure", TestComp265Failure)
}
