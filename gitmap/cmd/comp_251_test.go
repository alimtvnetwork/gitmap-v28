package cmd

import (
	"testing"
)

func TestComp251Success(t *testing.T) {
	input := Input251{ID: Comp251Uniqueness}
	output, err := HandleComp251(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp251Failure(t *testing.T) {
	input := Input251{ID: ""}
	output, err := HandleComp251(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp251(t *testing.T) {
	t.Run("success", TestComp251Success)
	t.Run("failure", TestComp251Failure)
}
