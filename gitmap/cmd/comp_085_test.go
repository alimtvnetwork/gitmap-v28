package cmd

import (
	"testing"
)

func TestComp085Success(t *testing.T) {
	input := Input085{ID: Comp085Uniqueness}
	output, err := HandleComp085(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp085Failure(t *testing.T) {
	input := Input085{ID: ""}
	output, err := HandleComp085(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp085(t *testing.T) {
	t.Run("success", TestComp085Success)
	t.Run("failure", TestComp085Failure)
}
