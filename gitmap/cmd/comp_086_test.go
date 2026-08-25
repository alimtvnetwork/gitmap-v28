package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func TestComp086(t *testing.T) {
	// Success case
	in := Input086{ID: "68519a9eca55"}
	out, err := HandleComp086(in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("expected Result to be true")
	}

	// Failure case
	inFail := Input086{ID: "wrong"}
	_, err = HandleComp086(inFail)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	
	appErr, ok := err.(*apperror.AppError)
	if !ok {
		t.Fatalf("expected apperror.AppError, got %T", err)
	}
	if appErr.Code != "E_COMP_086_FAIL" {
		t.Errorf("expected error code E_COMP_086_FAIL, got %s", appErr.Code)
	}
}
