package cmd

import (
	"testing"
)

func TestComp006Success(testRunner *testing.T) {
	input := Input006{ID: comp006Uniqueness}
	output, err := HandleComp006(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp006Failure(testRunner *testing.T) {
	input := Input006{ID: ""}
	output, err := HandleComp006(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp006(testRunner *testing.T) {
	testRunner.Run("success", TestComp006Success)
	testRunner.Run("failure", TestComp006Failure)
}
