package cmd

import (
	"errors"
	"testing"
)

func TestIsProcessLockError(t *testing.T) {
	if isProcessLockError(nil) {
		t.Errorf("expected false for nil error")
	}

	lockedErr := errors.New("unlinkat D:\\work\\Antigravity: The process cannot access the file because it is being used by another process")
	if !isProcessLockError(lockedErr) {
		t.Errorf("expected true for windows process lock error")
	}

	deniedErr := errors.New("remove D:\\test: Access is denied")
	if !isProcessLockError(deniedErr) {
		t.Errorf("expected true for access is denied error")
	}

	otherErr := errors.New("file does not exist")
	if isProcessLockError(otherErr) {
		t.Errorf("expected false for non-lock error")
	}
}
