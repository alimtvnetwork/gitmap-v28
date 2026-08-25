package cmd

import (
	"testing"
)

func TestComp096Success(t *testing.T) {
	input := Input096{ID: Comp096Uniqueness}
	output, err := HandleComp096(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp096Failure(t *testing.T) {
	input := Input096{ID: ""}
	output, err := HandleComp096(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp096(t *testing.T) {
	t.Run("success", TestComp096Success)
	t.Run("failure", TestComp096Failure)
}
