package cmd

import (
	"testing"
)

func TestComp048Success(t *testing.T) {
	input := Input048{ID: Comp048Uniqueness}
	output, err := HandleComp048(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp048Failure(t *testing.T) {
	input := Input048{ID: ""}
	output, err := HandleComp048(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp048(t *testing.T) {
	t.Run("success", TestComp048Success)
	t.Run("failure", TestComp048Failure)
}
