package cmd

import (
	"testing"
)

func TestComp123Success(t *testing.T) {
	input := Input123{ID: Comp123Uniqueness}
	output, err := HandleComp123(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp123Failure(t *testing.T) {
	input := Input123{ID: ""}
	output, err := HandleComp123(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp123(t *testing.T) {
	t.Run("success", TestComp123Success)
	t.Run("failure", TestComp123Failure)
}
