package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	appErr, ok := err.(*apperror.AppError)
	if !ok {
		t.Fatalf("expected *apperror.AppError, got %T", err)
	}

	if appErr.Code != "E_COMP_109_FAIL" {
		t.Errorf("expected code E_COMP_109_FAIL, got %s", appErr.Code)
	}
}
