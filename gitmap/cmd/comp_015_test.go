package cmd

import (
	"testing"
)

func TestComp015Success(testRunner *testing.T) {
	input := Input015{ID: comp015Uniqueness}
	output, err := HandleComp015(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp015Failure(testRunner *testing.T) {
	input := Input015{ID: ""}
	output, err := HandleComp015(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp015(testRunner *testing.T) {
	testRunner.Run("success", TestComp015Success)
	testRunner.Run("failure", TestComp015Failure)
}
