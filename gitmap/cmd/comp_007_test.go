package cmd

import (
	"testing"
)

func TestComp007Success(testRunner *testing.T) {
	input := Input007{ID: comp007Uniqueness}
	output, err := HandleComp007(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp007Failure(testRunner *testing.T) {
	input := Input007{ID: ""}
	output, err := HandleComp007(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp007(testRunner *testing.T) {
	testRunner.Run("success", TestComp007Success)
	testRunner.Run("failure", TestComp007Failure)
}
