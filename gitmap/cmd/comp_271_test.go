package cmd

import (
	"testing"
)

func TestComp271Success(t *testing.T) {
	input := Input271{ID: Comp271Uniqueness}
	output, err := HandleComp271(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp271Failure(t *testing.T) {
	input := Input271{ID: ""}
	output, err := HandleComp271(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp271(t *testing.T) {
	t.Run("success", TestComp271Success)
	t.Run("failure", TestComp271Failure)
}
