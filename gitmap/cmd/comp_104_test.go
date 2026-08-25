package cmd

import (
	"testing"
)

func TestComp104(t *testing.T) {
	in := Input104{ID: "valid"}
	out, err := HandleComp104(in)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Fatalf("Expected Result to be true, got %v", out.Result)
	}

	failIn := Input104{ID: "fail"}
	_, err = HandleComp104(failIn)
	if err == nil {
		t.Fatalf("Expected error for fail ID, got nil")
	}
}
