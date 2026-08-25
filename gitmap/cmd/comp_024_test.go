package cmd

import (
	"testing"
)

func TestComp024Success(testRunner *testing.T) {
	input := Input024{ID: comp024Uniqueness}
	output, err := HandleComp024(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp024Failure(testRunner *testing.T) {
	input := Input024{ID: ""}
	output, err := HandleComp024(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp024(testRunner *testing.T) {
	testRunner.Run("success", TestComp024Success)
	testRunner.Run("failure", TestComp024Failure)
}
