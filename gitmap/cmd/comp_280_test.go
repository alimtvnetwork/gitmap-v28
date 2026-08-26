package cmd

import (
	"testing"
)

func TestComp280Success(t *testing.T) {
	input := Input280{ID: Comp280Uniqueness}
	output, err := HandleComp280(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		t.Fatalf("expected Result to be true, got false")
	}
}

func TestComp280Failure(t *testing.T) {
	input := Input280{ID: ""}
	output, err := HandleComp280(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output.Result {
		t.Fatalf("expected Result to be false, got true")
	}
}

func TestComp280(t *testing.T) {
	t.Run("success", TestComp280Success)
	t.Run("failure", TestComp280Failure)
}
