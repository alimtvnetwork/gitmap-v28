package cmd

import (
	"testing"
)

func TestComp025Success(testRunner *testing.T) {
	input := Input025{ID: comp025Uniqueness}
	output, err := HandleComp025(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp025Failure(testRunner *testing.T) {
	input := Input025{ID: ""}
	output, err := HandleComp025(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp025(testRunner *testing.T) {
	testRunner.Run("success", TestComp025Success)
	testRunner.Run("failure", TestComp025Failure)
}
