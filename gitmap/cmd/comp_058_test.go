package cmd

import (
	"testing"
)

func TestComp058(t *testing.T) {
	in := Input058{ID: "e5b861a6d8a9"}
	out, err := HandleComp058(in)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected Result to be true, got false")
	}

	// Test error case
	inErr := Input058{ID: ""}
	_, err2 := HandleComp058(inErr)
	if err2 == nil {
		t.Errorf("Expected error E_COMP_058_FAIL, got nil")
	}
}
