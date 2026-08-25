package cmd

import (
	"testing"
)

func TestComp035Success(testRunner *testing.T) {
	input := Input035{ID: comp035Uniqueness}
	output, err := HandleComp035(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp035Failure(testRunner *testing.T) {
	input := Input035{ID: ""}
	output, err := HandleComp035(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp035(testRunner *testing.T) {
	testRunner.Run("success", TestComp035Success)
	testRunner.Run("failure", TestComp035Failure)
}
