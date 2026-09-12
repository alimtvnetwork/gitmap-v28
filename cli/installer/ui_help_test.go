// Package installer — ui_help_test.go tests detailed help output.
package installer

import (
	"strings"
	"testing"
)

func TestDetailedHelp(t *testing.T) {
	help := PrintDetailedHelp()
	if !strings.Contains(help, "Gitmap Installer Management System") || !strings.Contains(help, "export-all") {
		t.Errorf("unexpected help content: %s", help)
	}
}
