package cmd

import (
	"testing"
)

func TestComp041(t *testing.T) {
	in := Input041{ID: "a46e37632fa6"}
	out, err := HandleComp041(in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Fatalf("expected result to be true")
	}

	inFail := Input041{ID: ""}
	_, err = HandleComp041(inFail)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
