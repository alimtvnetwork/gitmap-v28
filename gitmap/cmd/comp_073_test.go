package cmd

import (
	"testing"
)

func TestComp073(t *testing.T) {
	in := Input073{ID: "0a5b046d07f6"}
	out, err := HandleComp073(in)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected result true, got false")
	}

	inEmpty := Input073{ID: ""}
	_, errEmpty := HandleComp073(inEmpty)
	if errEmpty == nil {
		t.Errorf("Expected error for empty ID, got nil")
	}
}
