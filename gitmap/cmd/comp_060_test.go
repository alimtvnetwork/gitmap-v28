package cmd

import (
	"testing"
)

func TestComp060Success(t *testing.T) {
	input := Input060{ID: Comp060Uniqueness}
	output, err := HandleComp060(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp060Failure(t *testing.T) {
	input := Input060{ID: ""}
	output, err := HandleComp060(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp060(t *testing.T) {
	t.Run("success", TestComp060Success)
	t.Run("failure", TestComp060Failure)
}
