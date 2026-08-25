package cmd

import (
	"testing"
)

func TestComp097(t *testing.T) {
	in := Input097{ID: "7559ca4a957c"}
	out, err := HandleComp097(in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("expected Result to be true")
	}

	inErr := Input097{ID: ""}
	outErr, errErr := HandleComp097(inErr)
	if errErr == nil {
		t.Fatalf("expected error for empty ID")
	}
	if outErr.Result {
		t.Errorf("expected Result to be false")
	}
}
