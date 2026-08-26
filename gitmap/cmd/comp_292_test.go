package cmd

import (
	"testing"
)

func TestComp292Success(t *testing.T) {
	input := Input292{ID: Comp292Uniqueness}
	output, err := HandleComp292(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp292Failure(t *testing.T) {
	input := Input292{ID: ""}
	output, err := HandleComp292(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp292(t *testing.T) {
	t.Run("success", TestComp292Success)
	t.Run("failure", TestComp292Failure)
}
