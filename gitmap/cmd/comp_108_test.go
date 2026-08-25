package cmd

import (
	"testing"
)

func TestComp108Success(t *testing.T) {
	input := Input108{ID: Comp108Uniqueness}
	output, err := HandleComp108(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp108Failure(t *testing.T) {
	input := Input108{ID: ""}
	output, err := HandleComp108(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp108(t *testing.T) {
	t.Run("success", TestComp108Success)
	t.Run("failure", TestComp108Failure)
}
