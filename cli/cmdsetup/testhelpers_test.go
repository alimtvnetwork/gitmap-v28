package cmdsetup

import (
	"strings"
	"testing"
)

func assertContains(t *testing.T, content, substr string) {
	t.Helper()
	hasMatch := strings.Contains(content, substr)
	if hasMatch {
		return
	}

	t.Errorf("expected content to contain %q", substr)
}
