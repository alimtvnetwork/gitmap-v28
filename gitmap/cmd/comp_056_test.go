package cmd

import (
	"testing"
)

func TestComp056Success(t *testing.T) {
	input := Input056{ID: Comp056Uniqueness}
	output, err := HandleComp056(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp056Failure(t *testing.T) {
	input := Input056{ID: ""}
	output, err := HandleComp056(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp056(t *testing.T) {
	t.Run("success", TestComp056Success)
	t.Run("failure", TestComp056Failure)
}
