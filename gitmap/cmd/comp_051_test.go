package cmd

import (
	"testing"
)

func TestComp051Success(t *testing.T) {
	input := Input051{ID: Comp051Uniqueness}
	output, err := HandleComp051(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp051Failure(t *testing.T) {
	input := Input051{ID: ""}
	output, err := HandleComp051(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp051(t *testing.T) {
	t.Run("success", TestComp051Success)
	t.Run("failure", TestComp051Failure)
}
