package cmd

import (
	"testing"
)

func TestComp017Success(testRunner *testing.T) {
	input := Input017{ID: comp017Uniqueness}
	output, err := HandleComp017(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp017Failure(testRunner *testing.T) {
	input := Input017{ID: ""}
	output, err := HandleComp017(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp017(testRunner *testing.T) {
	testRunner.Run("success", TestComp017Success)
	testRunner.Run("failure", TestComp017Failure)
}
