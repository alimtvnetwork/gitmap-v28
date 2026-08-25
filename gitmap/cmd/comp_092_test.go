package cmd

import (
	"testing"
)

func TestComp092Success(t *testing.T) {
	input := Input092{ID: Comp092Uniqueness}
	output, err := HandleComp092(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp092Failure(t *testing.T) {
	input := Input092{ID: ""}
	output, err := HandleComp092(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp092(t *testing.T) {
	t.Run("success", TestComp092Success)
	t.Run("failure", TestComp092Failure)
}
