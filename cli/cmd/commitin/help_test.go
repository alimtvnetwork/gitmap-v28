package commitin

import (
	"strings"
	"testing"
)

func TestCommitInHelp(t *testing.T) {
	help := PrintCommitInHelp()
	if !strings.Contains(help, "FuncIntel") || !strings.Contains(help, "SEO Commit Scheduling") {
		t.Fatalf("unexpected help text: %s", help)
	}
}
