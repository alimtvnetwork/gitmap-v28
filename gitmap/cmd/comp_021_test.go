package cmd

import (
	"testing"
)

func TestComp021Success(testRunner *testing.T) {
	input := Input021{ID: comp021Uniqueness}
	output, err := HandleComp021(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp021Failure(testRunner *testing.T) {
	input := Input021{ID: ""}
	output, err := HandleComp021(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp021(testRunner *testing.T) {
	testRunner.Run("success", TestComp021Success)
	testRunner.Run("failure", TestComp021Failure)
}
