package cmd

import (
	"testing"
)

func TestComp027Success(testRunner *testing.T) {
	input := Input027{ID: comp027Uniqueness}
	output, err := HandleComp027(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp027Failure(testRunner *testing.T) {
	input := Input027{ID: ""}
	output, err := HandleComp027(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp027(testRunner *testing.T) {
	testRunner.Run("success", TestComp027Success)
	testRunner.Run("failure", TestComp027Failure)
}
