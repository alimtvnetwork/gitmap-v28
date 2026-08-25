package cmd

import (
	"testing"
)

func TestComp028Success(testRunner *testing.T) {
	input := Input028{ID: comp028Uniqueness}
	output, err := HandleComp028(input)
	if err != nil {
		testRunner.Fatalf("unexpected error: %v", err)
	}

	if !output.Result {
		testRunner.Fatalf("expected Result to be true, got false")
	}
}

func TestComp028Failure(testRunner *testing.T) {
	input := Input028{ID: ""}
	output, err := HandleComp028(input)
	if err == nil {
		testRunner.Fatalf("expected error, got nil")
	}

	if output.Result {
		testRunner.Fatalf("expected Result to be false, got true")
	}
}

func TestComp028(testRunner *testing.T) {
	testRunner.Run("success", TestComp028Success)
	testRunner.Run("failure", TestComp028Failure)
}
