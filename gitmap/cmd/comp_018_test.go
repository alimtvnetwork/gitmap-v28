package cmd

import (
	"testing"
)

func TestComp018Success(testRunner *testing.T) {
	input := Input018{ID: comp018Uniqueness}
	output, err := HandleComp018(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp018Failure(testRunner *testing.T) {
	input := Input018{ID: ""}
	output, err := HandleComp018(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp018(testRunner *testing.T) {
	testRunner.Run("success", TestComp018Success)
	testRunner.Run("failure", TestComp018Failure)
}
