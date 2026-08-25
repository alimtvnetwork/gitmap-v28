package cmd

import (
	"testing"
)

func TestComp014Success(testRunner *testing.T) {
	input := Input014{ID: comp014Uniqueness}
	output, err := HandleComp014(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp014Failure(testRunner *testing.T) {
	input := Input014{ID: ""}
	output, err := HandleComp014(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp014(testRunner *testing.T) {
	testRunner.Run("success", TestComp014Success)
	testRunner.Run("failure", TestComp014Failure)
}
