package cmd

import (
	"testing"
)

func TestComp246Success(t *testing.T) {
	input := Input246{ID: Comp246Uniqueness}
	output, err := HandleComp246(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp246Failure(t *testing.T) {
	input := Input246{ID: ""}
	output, err := HandleComp246(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp246(t *testing.T) {
	t.Run("success", TestComp246Success)
	t.Run("failure", TestComp246Failure)
}
