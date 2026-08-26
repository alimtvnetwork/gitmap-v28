// Package installer — ui_examples_test.go tests example output.
package installer

import (
	"strings"
	"testing"
)

func TestComposerExample(t *testing.T) {
	ex := PrintComposerExample()
	if !strings.Contains(ex, "Composer Installer Configuration Example") || !strings.Contains(ex, "PHP Composer") {
		t.Errorf("unexpected example output: %s", ex)
	}
}
