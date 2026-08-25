package cmd

import (
	"testing"
)

func TestComp011Success(testRunner *testing.T) {
	input := Input011{ID: comp011Uniqueness}
	output, err := HandleComp011(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp011Failure(testRunner *testing.T) {
	input := Input011{ID: ""}
	output, err := HandleComp011(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp011(testRunner *testing.T) {
	testRunner.Run("success", TestComp011Success)
	testRunner.Run("failure", TestComp011Failure)
}
