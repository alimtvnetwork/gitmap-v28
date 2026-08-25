package cmd

import (
	"testing"
)

func TestComp013Success(testRunner *testing.T) {
	input := Input013{ID: comp013Uniqueness}
	output, err := HandleComp013(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp013Failure(testRunner *testing.T) {
	input := Input013{ID: ""}
	output, err := HandleComp013(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp013(testRunner *testing.T) {
	testRunner.Run("success", TestComp013Success)
	testRunner.Run("failure", TestComp013Failure)
}
