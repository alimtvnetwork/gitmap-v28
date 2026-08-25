package cmd

import (
	"testing"
)

func TestComp065Success(t *testing.T) {
	input := Input065{ID: Comp065Uniqueness}
	output, err := HandleComp065(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp065Failure(t *testing.T) {
	input := Input065{ID: ""}
	output, err := HandleComp065(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp065(t *testing.T) {
	t.Run("success", TestComp065Success)
	t.Run("failure", TestComp065Failure)
}
