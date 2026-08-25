package cmd

import (
	"testing"
)

func TestComp009Success(testRunner *testing.T) {
	input := Input009{ID: comp009Uniqueness}
	output, err := HandleComp009(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp009Failure(testRunner *testing.T) {
	input := Input009{ID: ""}
	output, err := HandleComp009(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp009(testRunner *testing.T) {
	testRunner.Run("success", TestComp009Success)
	testRunner.Run("failure", TestComp009Failure)
}
