package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"testing"
)

func TestComp053(t *testing.T) {
	in := Input053{ID: "482d9673cfee"}
	out, err := HandleComp053(in)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !out.Result {
		t.Errorf("expected true result")
	}

	inFail := Input053{ID: ""}
	_, errFail := HandleComp053(inFail)
	
	if errFail == nil {
		t.Errorf("expected error, got nil")
	} else {
		appErr, ok := errFail.(*apperror.AppError)
		if !ok {
			t.Errorf("expected *apperror.AppError, got %T", errFail)
		} else if appErr.Code != "E_COMP_053_FAIL" {
			t.Errorf("expected code E_COMP_053_FAIL, got %s", appErr.Code)
		}
	}
}
