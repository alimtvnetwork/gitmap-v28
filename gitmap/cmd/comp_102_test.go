package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func TestComp102(t *testing.T) {
	out, err := HandleComp102(Input102{ID: "test-id"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("expected Result to be true")
	}

	_, err = HandleComp102(Input102{ID: ""})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	appErr, ok := err.(*apperror.AppError)
	if !ok {
		t.Fatalf("expected error to be of type *apperror.AppError")
	}

	if appErr.Code != "E_COMP_102_FAIL" {
		t.Errorf("expected code E_COMP_102_FAIL, got %s", appErr.Code)
	}
}
