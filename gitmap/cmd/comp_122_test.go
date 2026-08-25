package cmd

import (
	"testing"
)

func TestComp122Success(t *testing.T) {
	input := Input122{ID: Comp122Uniqueness}
	output, err := HandleComp122(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp122Failure(t *testing.T) {
	input := Input122{ID: ""}
	output, err := HandleComp122(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp122(t *testing.T) {
	t.Run("success", TestComp122Success)
	t.Run("failure", TestComp122Failure)
}
