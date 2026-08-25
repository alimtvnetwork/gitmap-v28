package cmd

import (
	"testing"
)

func TestComp045Success(t *testing.T) {
	input := Input045{ID: Comp045Uniqueness}
	output, err := HandleComp045(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp045Failure(t *testing.T) {
	input := Input045{ID: ""}
	output, err := HandleComp045(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp045(t *testing.T) {
	t.Run("success", TestComp045Success)
	t.Run("failure", TestComp045Failure)
}
