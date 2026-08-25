package cmd

import (
	"testing"
)

func TestComp106Success(t *testing.T) {
	input := Input106{ID: Comp106Uniqueness}
	output, err := HandleComp106(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp106Failure(t *testing.T) {
	input := Input106{ID: ""}
	output, err := HandleComp106(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp106(t *testing.T) {
	t.Run("success", TestComp106Success)
	t.Run("failure", TestComp106Failure)
}
