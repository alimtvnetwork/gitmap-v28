package cmd

import (
	"testing"
)

func TestComp044Success(t *testing.T) {
	input := Input044{ID: Comp044Uniqueness}
	output, err := HandleComp044(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp044Failure(t *testing.T) {
	input := Input044{ID: ""}
	output, err := HandleComp044(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp044(t *testing.T) {
	t.Run("success", TestComp044Success)
	t.Run("failure", TestComp044Failure)
}
