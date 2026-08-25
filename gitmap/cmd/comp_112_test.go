package cmd

import (
	"testing"
)

func TestComp112Success(t *testing.T) {
	input := Input112{ID: Comp112Uniqueness}
	output, err := HandleComp112(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp112Failure(t *testing.T) {
	input := Input112{ID: ""}
	output, err := HandleComp112(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp112(t *testing.T) {
	t.Run("success", TestComp112Success)
	t.Run("failure", TestComp112Failure)
}
