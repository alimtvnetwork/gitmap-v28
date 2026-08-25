package cmd

import (
	"testing"
)

func TestComp120(t *testing.T) {
	out, err := HandleComp120(Input120{ID: "69cce754be0e"})
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected true result")
	}

	out2, err2 := HandleComp120(Input120{ID: "invalid"})
	if err2 == nil {
		t.Fatalf("Expected error, got nil")
	}
	if out2.Result {
		t.Errorf("Expected false result")
	}
}
