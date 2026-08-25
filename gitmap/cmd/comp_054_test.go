package cmd

import (
	"testing"
)

func TestComp054(t *testing.T) {
	in := Input054{ID: "9537f32ec759"}
	out, err := HandleComp054(in)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected Result to be true, got false")
	}

	// Test error case
	inErr := Input054{ID: ""}
	_, err2 := HandleComp054(inErr)
	if err2 == nil {
		t.Errorf("Expected error E_COMP_054_FAIL, got nil")
	}
}
