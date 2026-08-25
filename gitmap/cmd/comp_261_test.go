package cmd

import (
	"testing"
)

func TestComp261Success(t *testing.T) {
	input := Input261{ID: Comp261Uniqueness}
	output, err := HandleComp261(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp261Failure(t *testing.T) {
	input := Input261{ID: ""}
	output, err := HandleComp261(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp261(t *testing.T) {
	t.Run("success", TestComp261Success)
	t.Run("failure", TestComp261Failure)
}
