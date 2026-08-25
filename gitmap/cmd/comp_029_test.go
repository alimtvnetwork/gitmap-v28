package cmd

import (
	"testing"
)

func TestComp029Success(testRunner *testing.T) {
	input := Input029{ID: comp029Uniqueness}
	output, err := HandleComp029(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp029Failure(testRunner *testing.T) {
	input := Input029{ID: ""}
	output, err := HandleComp029(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp029(testRunner *testing.T) {
	testRunner.Run("success", TestComp029Success)
	testRunner.Run("failure", TestComp029Failure)
}
