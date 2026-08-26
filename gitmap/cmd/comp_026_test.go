package cmd

import (
	"testing"
)

func TestComp026Success(testRunner *testing.T) {
	input := Input026{ID: comp026Uniqueness}
	output, err := HandleComp026(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp026Failure(testRunner *testing.T) {
	input := Input026{ID: ""}
	output, err := HandleComp026(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp026(testRunner *testing.T) {
	testRunner.Run("success", TestComp026Success)
	testRunner.Run("failure", TestComp026Failure)
}
