package cmd

import (
	"testing"
)

func TestComp131Success(t *testing.T) {
	input := Input131{ID: Comp131Uniqueness}
	output, err := HandleComp131(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp131Failure(t *testing.T) {
	input := Input131{ID: ""}
	output, err := HandleComp131(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp131(t *testing.T) {
	t.Run("success", TestComp131Success)
	t.Run("failure", TestComp131Failure)
}
