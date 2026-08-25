package cmd

import (
	"testing"
)

func TestComp030Success(testRunner *testing.T) {
	input := Input030{ID: comp030Uniqueness}
	output, err := HandleComp030(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp030Failure(testRunner *testing.T) {
	input := Input030{ID: ""}
	output, err := HandleComp030(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp030(testRunner *testing.T) {
	testRunner.Run("success", TestComp030Success)
	testRunner.Run("failure", TestComp030Failure)
}
