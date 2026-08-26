package cmd

import (
	"testing"
)

func TestComp300Success(t *testing.T) {
	input := Input300{ID: Comp300Uniqueness}
	output, err := HandleComp300(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp300Failure(t *testing.T) {
	input := Input300{ID: ""}
	output, err := HandleComp300(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp300(t *testing.T) {
	t.Run("success", TestComp300Success)
	t.Run("failure", TestComp300Failure)
}
