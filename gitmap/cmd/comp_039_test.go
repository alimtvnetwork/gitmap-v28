package cmd

import (
	"testing"
)

func TestComp039Success(t *testing.T) {
	input := Input039{ID: Comp039Uniqueness}
	output, err := HandleComp039(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp039Failure(t *testing.T) {
	input := Input039{ID: ""}
	output, err := HandleComp039(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp039(t *testing.T) {
	t.Run("success", TestComp039Success)
	t.Run("failure", TestComp039Failure)
}
