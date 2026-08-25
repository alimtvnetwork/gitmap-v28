package cmd

import (
	"testing"
)

func TestComp242Success(t *testing.T) {
	input := Input242{ID: Comp242Uniqueness}
	output, err := HandleComp242(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp242Failure(t *testing.T) {
	input := Input242{ID: ""}
	output, err := HandleComp242(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp242(t *testing.T) {
	t.Run("success", TestComp242Success)
	t.Run("failure", TestComp242Failure)
}
