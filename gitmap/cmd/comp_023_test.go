package cmd

import (
	"testing"
)

func TestComp023Success(testRunner *testing.T) {
	input := Input023{ID: comp023Uniqueness}
	output, err := HandleComp023(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp023Failure(testRunner *testing.T) {
	input := Input023{ID: ""}
	output, err := HandleComp023(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp023(testRunner *testing.T) {
	testRunner.Run("success", TestComp023Success)
	testRunner.Run("failure", TestComp023Failure)
}
