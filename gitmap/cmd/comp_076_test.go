package cmd

import (
	"testing"
)

func TestComp076(t *testing.T) {
	in := Input076{ID: "f74efabef12e"}
	out, err := HandleComp076(in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Fatalf("expected Result to be true")
	}

	inErr := Input076{ID: "invalid"}
	outErr, err2 := HandleComp076(inErr)
	if err2 == nil {
		t.Fatalf("expected error, got nil")
	}
	if outErr.Result {
		t.Fatalf("expected Result to be false on error")
	}
}
