package cmd

import "testing"

func TestHandleCloneNextDryRun_WhenDisabledReturnsFalse(t *testing.T) {
	isHandled := handleCloneNextDryRun(false, "https://github.com/example/repo", "/tmp/dest")
	if isHandled {
		t.Errorf("handleCloneNextDryRun(false) = true, want false")
	}
}
