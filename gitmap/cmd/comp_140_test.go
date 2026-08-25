package cmd

import (
	"testing"
)

func TestComp140(t *testing.T) {
	out, err := HandleComp140(Input140{ID: Comp140Uniqueness})
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected true result")
	}

	out2, err2 := HandleComp140(Input140{ID: "invalid"})
	if err2 == nil {
		t.Fatalf("Expected error, got nil")
	}
	if out2.Result {
		t.Errorf("Expected false result")
	}
}
