package cmd

import (
	"testing"
)

func TestComp016Success(testRunner *testing.T) {
	input := Input016{ID: comp016Uniqueness}
	output, err := HandleComp016(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp016Failure(testRunner *testing.T) {
	input := Input016{ID: ""}
	output, err := HandleComp016(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp016(testRunner *testing.T) {
	testRunner.Run("success", TestComp016Success)
	testRunner.Run("failure", TestComp016Failure)
}
