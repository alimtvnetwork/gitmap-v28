package cmd

import (
	"testing"
)

func TestComp077Success(t *testing.T) {
	input := Input077{ID: Comp077Uniqueness}
	output, err := HandleComp077(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp077Failure(t *testing.T) {
	input := Input077{ID: ""}
	output, err := HandleComp077(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp077(t *testing.T) {
	t.Run("success", TestComp077Success)
	t.Run("failure", TestComp077Failure)
}
