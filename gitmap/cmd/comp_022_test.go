package cmd

import (
	"testing"
)

func TestComp022Success(testRunner *testing.T) {
	input := Input022{ID: comp022Uniqueness}
	output, err := HandleComp022(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp022Failure(testRunner *testing.T) {
	input := Input022{ID: ""}
	output, err := HandleComp022(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp022(testRunner *testing.T) {
	testRunner.Run("success", TestComp022Success)
	testRunner.Run("failure", TestComp022Failure)
}
