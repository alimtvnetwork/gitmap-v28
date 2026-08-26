package cmd

import (
	"testing"
)

func TestComp285Success(t *testing.T) {
	input := Input285{ID: Comp285Uniqueness}
	output, err := HandleComp285(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp285Failure(t *testing.T) {
	input := Input285{ID: ""}
	output, err := HandleComp285(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp285(t *testing.T) {
	t.Run("success", TestComp285Success)
	t.Run("failure", TestComp285Failure)
}
