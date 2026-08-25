package cmd

import (
	"testing"
)

func TestComp036Success(testRunner *testing.T) {
	input := Input036{ID: comp036Uniqueness}
	output, err := HandleComp036(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp036Failure(testRunner *testing.T) {
	input := Input036{ID: ""}
	output, err := HandleComp036(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp036(testRunner *testing.T) {
	testRunner.Run("success", TestComp036Success)
	testRunner.Run("failure", TestComp036Failure)
}
