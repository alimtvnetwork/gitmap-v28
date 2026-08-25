package cmd

import (
	"testing"
)

func TestComp109(t *testing.T) {
	in := Input109{ID: "5966abd0cbfc"}
	out, err := HandleComp109(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Result {
		t.Errorf("expected Result true, got false")
	}

	inErr := Input109{ID: "fail"}
	_, err = HandleComp109(inErr)
	if err != ErrComp109Fail {
		t.Errorf("expected ErrComp109Fail, got %v", err)
	}
}
