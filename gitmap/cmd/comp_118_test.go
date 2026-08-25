package cmd

import (
	"testing"
)

func TestComp118(t *testing.T) {
	out, err := HandleComp118(Input118{ID: "9a049b03f6fc"})
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected true result")
	}

	out2, err2 := HandleComp118(Input118{ID: "invalid"})
	if err2 == nil {
		t.Fatalf("Expected error, got nil")
	}
	if out2.Result {
		t.Errorf("Expected false result")
	}
}
