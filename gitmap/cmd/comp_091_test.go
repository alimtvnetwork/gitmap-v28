package cmd

import (
	"testing"
)

func TestComp091Success(t *testing.T) {
	input := Input091{ID: Comp091Uniqueness}
	output, err := HandleComp091(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp091Failure(t *testing.T) {
	input := Input091{ID: ""}
	output, err := HandleComp091(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp091(t *testing.T) {
	t.Run("success", TestComp091Success)
	t.Run("failure", TestComp091Failure)
}
