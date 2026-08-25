package cmd

import (
	"testing"
)

func TestComp090Success(t *testing.T) {
	input := Input090{ID: Comp090Uniqueness}
	output, err := HandleComp090(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp090Failure(t *testing.T) {
	input := Input090{ID: ""}
	output, err := HandleComp090(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp090(t *testing.T) {
	t.Run("success", TestComp090Success)
	t.Run("failure", TestComp090Failure)
}
