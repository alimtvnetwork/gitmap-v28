package cmd

import (
	"testing"
)

func TestComp129Success(t *testing.T) {
	input := Input129{ID: Comp129Uniqueness}
	output, err := HandleComp129(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp129Failure(t *testing.T) {
	input := Input129{ID: ""}
	output, err := HandleComp129(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp129(t *testing.T) {
	t.Run("success", TestComp129Success)
	t.Run("failure", TestComp129Failure)
}
