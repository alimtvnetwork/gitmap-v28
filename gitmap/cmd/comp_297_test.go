package cmd

import (
	"testing"
)

func TestComp297Success(t *testing.T) {
	input := Input297{ID: Comp297Uniqueness}
	output, err := HandleComp297(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp297Failure(t *testing.T) {
	input := Input297{ID: ""}
	output, err := HandleComp297(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp297(t *testing.T) {
	t.Run("success", TestComp297Success)
	t.Run("failure", TestComp297Failure)
}
