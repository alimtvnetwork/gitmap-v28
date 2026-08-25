package cmd

import (
	"testing"
)

func TestComp033Success(testRunner *testing.T) {
	input := Input033{ID: comp033Uniqueness}
	output, err := HandleComp033(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp033Failure(testRunner *testing.T) {
	input := Input033{ID: ""}
	output, err := HandleComp033(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp033(testRunner *testing.T) {
	testRunner.Run("success", TestComp033Success)
	testRunner.Run("failure", TestComp033Failure)
}
