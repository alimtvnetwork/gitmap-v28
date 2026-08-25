package cmd

import (
	"testing"
)

func TestComp098Success(t *testing.T) {
	input := Input098{ID: Comp098Uniqueness}
	output, err := HandleComp098(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp098Failure(t *testing.T) {
	input := Input098{ID: ""}
	output, err := HandleComp098(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp098(t *testing.T) {
	t.Run("success", TestComp098Success)
	t.Run("failure", TestComp098Failure)
}
