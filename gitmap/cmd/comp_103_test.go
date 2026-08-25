package cmd

import (
	"testing"
)

func TestComp103Success(t *testing.T) {
	input := Input103{ID: Comp103Uniqueness}
	output, err := HandleComp103(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp103Failure(t *testing.T) {
	input := Input103{ID: ""}
	output, err := HandleComp103(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp103(t *testing.T) {
	t.Run("success", TestComp103Success)
	t.Run("failure", TestComp103Failure)
}
