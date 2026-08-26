package cmd

import (
	"testing"
)

func TestComp277Success(t *testing.T) {
	input := Input277{ID: Comp277Uniqueness}
	output, err := HandleComp277(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp277Failure(t *testing.T) {
	input := Input277{ID: ""}
	output, err := HandleComp277(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp277(t *testing.T) {
	t.Run("success", TestComp277Success)
	t.Run("failure", TestComp277Failure)
}
