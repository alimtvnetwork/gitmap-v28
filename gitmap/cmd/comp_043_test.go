package cmd

import (
	"testing"
)

func TestComp043Success(t *testing.T) {
	input := Input043{ID: Comp043Uniqueness}
	output, err := HandleComp043(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp043Failure(t *testing.T) {
	input := Input043{ID: ""}
	output, err := HandleComp043(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp043(t *testing.T) {
	t.Run("success", TestComp043Success)
	t.Run("failure", TestComp043Failure)
}
