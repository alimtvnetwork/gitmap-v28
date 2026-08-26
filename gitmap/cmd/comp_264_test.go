package cmd

import (
	"testing"
)

func TestComp264Success(t *testing.T) {
	input := Input264{ID: Comp264Uniqueness}
	output, err := HandleComp264(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp264Failure(t *testing.T) {
	input := Input264{ID: ""}
	output, err := HandleComp264(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp264(t *testing.T) {
	t.Run("success", TestComp264Success)
	t.Run("failure", TestComp264Failure)
}
