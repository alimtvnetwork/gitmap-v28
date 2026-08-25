package cmd

import (
	"testing"
)

func TestComp034Success(testRunner *testing.T) {
	input := Input034{ID: comp034Uniqueness}
	output, err := HandleComp034(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp034Failure(testRunner *testing.T) {
	input := Input034{ID: ""}
	output, err := HandleComp034(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp034(testRunner *testing.T) {
	testRunner.Run("success", TestComp034Success)
	testRunner.Run("failure", TestComp034Failure)
}
