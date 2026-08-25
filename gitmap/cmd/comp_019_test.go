package cmd

import (
	"testing"
)

func TestComp019Success(testRunner *testing.T) {
	input := Input019{ID: comp019Uniqueness}
	output, err := HandleComp019(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp019Failure(testRunner *testing.T) {
	input := Input019{ID: ""}
	output, err := HandleComp019(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp019(testRunner *testing.T) {
	testRunner.Run("success", TestComp019Success)
	testRunner.Run("failure", TestComp019Failure)
}
