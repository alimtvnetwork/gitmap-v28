package cmd

import (
	"testing"
)

func TestComp100(t *testing.T) {
	// Test success case
	inSuccess := Input100{ID: "27badc983df1"}
	out, err := HandleComp100(inSuccess)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("expected Result to be true, got false")
	}

	// Test failure case
	inFail := Input100{ID: "wrong-id"}
	_, err = HandleComp100(inFail)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
