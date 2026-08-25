package cmd

import (
	"testing"
)

func TestComp020Success(testRunner *testing.T) {
	input := Input020{ID: comp020Uniqueness}
	output, err := HandleComp020(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp020Failure(testRunner *testing.T) {
	input := Input020{ID: ""}
	output, err := HandleComp020(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp020(testRunner *testing.T) {
	testRunner.Run("success", TestComp020Success)
	testRunner.Run("failure", TestComp020Failure)
}
