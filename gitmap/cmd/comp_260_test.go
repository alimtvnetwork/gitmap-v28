package cmd

import (
	"testing"
)

func TestComp260Success(t *testing.T) {
	input := Input260{ID: Comp260Uniqueness}
	output, err := HandleComp260(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp260Failure(t *testing.T) {
	input := Input260{ID: ""}
	output, err := HandleComp260(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp260(t *testing.T) {
	t.Run("success", TestComp260Success)
	t.Run("failure", TestComp260Failure)
}
