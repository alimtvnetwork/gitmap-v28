package cmd

import (
	"testing"
)

func TestComp236Success(testRunner *testing.T) {
	input := Input236{ID: comp236Uniqueness}
	output, err := HandleComp236(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp236Failure(testRunner *testing.T) {
	input := Input236{ID: ""}
	output, err := HandleComp236(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp236(testRunner *testing.T) {
	testRunner.Run("success", TestComp236Success)
	testRunner.Run("failure", TestComp236Failure)
}
