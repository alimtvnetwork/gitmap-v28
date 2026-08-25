package cmd

import (
	"testing"
)

func TestComp124Success(t *testing.T) {
	input := Input124{ID: Comp124Uniqueness}
	output, err := HandleComp124(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp124Failure(t *testing.T) {
	input := Input124{ID: ""}
	output, err := HandleComp124(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp124(t *testing.T) {
	t.Run("success", TestComp124Success)
	t.Run("failure", TestComp124Failure)
}
