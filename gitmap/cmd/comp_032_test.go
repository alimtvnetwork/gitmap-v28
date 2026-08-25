package cmd

import (
	"testing"
)

func TestComp032Success(testRunner *testing.T) {
	input := Input032{ID: comp032Uniqueness}
	output, err := HandleComp032(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp032Failure(testRunner *testing.T) {
	input := Input032{ID: ""}
	output, err := HandleComp032(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp032(testRunner *testing.T) {
	testRunner.Run("success", TestComp032Success)
	testRunner.Run("failure", TestComp032Failure)
}
