package cmd

import (
	"testing"
)

func TestComp121(t *testing.T) {
	out, err := HandleComp121(Input121{ID: "14063697603e"})
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected true result")
	}

	out2, err2 := HandleComp121(Input121{ID: "invalid"})
	if err2 == nil {
		t.Fatalf("Expected error, got nil")
	}
	if out2.Result {
		t.Errorf("Expected false result")
	}
}
