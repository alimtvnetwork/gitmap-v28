package cmd

import (
	"testing"
)

func TestComp253Success(t *testing.T) {
	input := Input253{ID: Comp253Uniqueness}
	output, err := HandleComp253(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp253Failure(t *testing.T) {
	input := Input253{ID: ""}
	output, err := HandleComp253(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp253(t *testing.T) {
	t.Run("success", TestComp253Success)
	t.Run("failure", TestComp253Failure)
}
