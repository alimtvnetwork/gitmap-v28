package cmd

import (
	"testing"
)

func TestComp010Success(testRunner *testing.T) {
	input := Input010{ID: comp010Uniqueness}
	output, err := HandleComp010(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp010Failure(testRunner *testing.T) {
	input := Input010{ID: ""}
	output, err := HandleComp010(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp010(testRunner *testing.T) {
	testRunner.Run("success", TestComp010Success)
	testRunner.Run("failure", TestComp010Failure)
}
