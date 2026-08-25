package cmd

import (
	"testing"
)

func TestComp012Success(testRunner *testing.T) {
	input := Input012{ID: comp012Uniqueness}
	output, err := HandleComp012(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp012Failure(testRunner *testing.T) {
	input := Input012{ID: ""}
	output, err := HandleComp012(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp012(testRunner *testing.T) {
	testRunner.Run("success", TestComp012Success)
	testRunner.Run("failure", TestComp012Failure)
}
