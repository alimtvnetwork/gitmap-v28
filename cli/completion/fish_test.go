// Package completion — fish_test.go verifies fish completion script generation.
package completion

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestGenerateFish(t *testing.T) {
	script, err := Generate(constants.ShellFish)
	if err != nil {
		t.Fatalf("Generate(fish) failed: %v", err)
	}

	if len(script) == 0 {
		t.Fatal("Generate(fish) returned empty script")
	}

	if !strings.Contains(script, "complete -c gitmap") {
		t.Errorf("expected complete -c gitmap in script, got: %s", script)
	}
}
