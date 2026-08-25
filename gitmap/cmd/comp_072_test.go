package cmd

import (
	"testing"
)

func TestComp072(t *testing.T) {
	in := Input072{ID: "5ec1a0c99d42"}
	out, err := HandleComp072(in)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("expected true, got %v", out.Result)
	}

	inFail := Input072{ID: ""}
	_, errFail := HandleComp072(inFail)
	if errFail == nil {
		t.Errorf("expected error E_COMP_072_FAIL")
	}
}
