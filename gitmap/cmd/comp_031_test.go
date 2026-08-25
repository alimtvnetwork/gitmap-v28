package cmd

import (
	"testing"
)

func TestComp031Success(testRunner *testing.T) {
	input := Input031{ID: comp031Uniqueness}
	output, err := HandleComp031(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp031Failure(testRunner *testing.T) {
	input := Input031{ID: ""}
	output, err := HandleComp031(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp031(testRunner *testing.T) {
	testRunner.Run("success", TestComp031Success)
	testRunner.Run("failure", TestComp031Failure)
}
