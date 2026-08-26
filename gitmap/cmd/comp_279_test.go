package cmd

import (
	"testing"
)

func TestComp279Success(t *testing.T) {
	input := Input279{ID: Comp279Uniqueness}
	output, err := HandleComp279(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp279Failure(t *testing.T) {
	input := Input279{ID: ""}
	output, err := HandleComp279(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp279(t *testing.T) {
	t.Run("success", TestComp279Success)
	t.Run("failure", TestComp279Failure)
}
