package cmd

import (
	"testing"
)

// TestRenderDynamicCommandHelp verifies that any arbitrary command produces
// a valid dynamic help menu without errors.
func TestRenderDynamicCommandHelp(t *testing.T) {
	rendered := RenderDynamicCommandHelp("unknowncommand")
	if !rendered {
		t.Fatalf("expected RenderDynamicCommandHelp to return true")
	}

	renderedClean := RenderDynamicCommandHelp("clean")
	if !renderedClean {
		t.Fatalf("expected RenderDynamicCommandHelp to return true for clean")
	}
}
