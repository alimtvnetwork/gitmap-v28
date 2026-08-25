package cmd

import (
	"testing"
)

func TestComp099(t *testing.T) {
	in := Input099{ID: "a4e00d7e6aa8"}
	out, err := HandleComp099(in)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if !out.Result {
		t.Errorf("Expected true result, got %v", out.Result)
	}

	inFail := Input099{ID: "wrong"}
	outFail, errFail := HandleComp099(inFail)
	if errFail == nil {
		t.Errorf("Expected error, got nil")
	}
	if outFail.Result {
		t.Errorf("Expected false result on failure")
	}
}
