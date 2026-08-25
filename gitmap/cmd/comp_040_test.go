package cmd

import (
	"testing"
)

func TestComp040Success(t *testing.T) {
	input := Input040{ID: Comp040Uniqueness}
	output, err := HandleComp040(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp040Failure(t *testing.T) {
	input := Input040{ID: ""}
	output, err := HandleComp040(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp040(t *testing.T) {
	t.Run("success", TestComp040Success)
	t.Run("failure", TestComp040Failure)
}
