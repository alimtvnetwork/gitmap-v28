package cmd

import (
	"testing"
)

func TestComp296Success(t *testing.T) {
	input := Input296{ID: Comp296Uniqueness}
	output, err := HandleComp296(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp296Failure(t *testing.T) {
	input := Input296{ID: ""}
	output, err := HandleComp296(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp296(t *testing.T) {
	t.Run("success", TestComp296Success)
	t.Run("failure", TestComp296Failure)
}
