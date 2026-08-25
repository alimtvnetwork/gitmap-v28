package cmd

import (
	"testing"
)

func TestComp066Success(t *testing.T) {
	input := Input066{ID: comp066Uniqueness}
	output, err := HandleComp066(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp066Failure(t *testing.T) {
	input := Input066{ID: ""}
	output, err := HandleComp066(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp066(t *testing.T) {
	t.Run("success", TestComp066Success)
	t.Run("failure", TestComp066Failure)
}
