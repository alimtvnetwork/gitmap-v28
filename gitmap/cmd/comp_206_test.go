package cmd

import (
	"testing"
)

func TestComp206(t *testing.T) {
	out, err := HandleComp206(Input206{ID: Comp206Uniqueness})
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected true result")
	}

	out2, err2 := HandleComp206(Input206{ID: "invalid"})
	if err2 == nil {
		t.Fatalf("Expected error, got nil")
	}
	if out2.Result {
		t.Errorf("Expected false result")
	}
}
