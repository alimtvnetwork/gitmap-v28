package cmd

import (
	"testing"
)

func TestComp174(t *testing.T) {
	out, err := HandleComp174(Input174{ID: Comp174Uniqueness})
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected true result")
	}

	out2, err2 := HandleComp174(Input174{ID: "invalid"})
	if err2 == nil {
		t.Fatalf("Expected error, got nil")
	}
	if out2.Result {
		t.Errorf("Expected false result")
	}
}
