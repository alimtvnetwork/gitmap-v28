package cmd

import (
	"testing"
)

func TestComp247Success(t *testing.T) {
	input := Input247{ID: Comp247Uniqueness}
	output, err := HandleComp247(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp247Failure(t *testing.T) {
	input := Input247{ID: ""}
	output, err := HandleComp247(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp247(t *testing.T) {
	t.Run("success", TestComp247Success)
	t.Run("failure", TestComp247Failure)
}
