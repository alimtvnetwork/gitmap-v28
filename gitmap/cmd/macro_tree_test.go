package cmd

import (
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func TestPrintMacroStepsTree_EmptyAndNil(t *testing.T) {
	t.Parallel()

	// Should not panic on nil macro
	printMacroStepsTree(nil)

	// Should not panic on empty steps macro
	emptyMacro := &macro.Macro{
		Name:  "empty-macro",
		Steps: []macro.MacroStep{},
	}
	printMacroStepsTree(emptyMacro)
}

func TestPrintMacroStepsTree_RenderSteps(t *testing.T) {
	t.Parallel()

	testMacro := &macro.Macro{
		ID:          1,
		Name:        "test-pipeline",
		Description: "Runs build and tests",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Steps: []macro.MacroStep{
			{
				StepNum:     1,
				CommandLine: "git status",
			},
			{
				StepNum:     2,
				CommandLine: "go build ./...",
				WorkingDir:  "/workspace",
			},
			{
				StepNum:     3,
				CommandLine: "go test ./...",
			},
		},
	}

	// Verify it executes without panic
	printMacroStepsTree(testMacro)
}

func TestPrintMacroStepsTree_NoDescription(t *testing.T) {
	t.Parallel()

	testMacro := &macro.Macro{
		Name: "simple-macro",
		Steps: []macro.MacroStep{
			{
				StepNum:     1,
				CommandLine: "echo hello",
			},
		},
	}

	printMacroStepsTree(testMacro)
}

func TestPrintMacroTree_NonExistent(t *testing.T) {
	t.Parallel()

	// Loading non-existent macro should gracefully return without panic
	printMacroTree("non-existent-macro-xyz-123")
}

func TestFormatStepDescription(t *testing.T) {
	t.Parallel()

	stepWithDir := macro.MacroStep{WorkingDir: "/app"}
	descWithDir := formatStepDescription(stepWithDir)
	if descWithDir != "(dir: /app)" {
		t.Errorf("expected '(dir: /app)', got %q", descWithDir)
	}

	stepEmpty := macro.MacroStep{}
	descEmpty := formatStepDescription(stepEmpty)
	if descEmpty != "" {
		t.Errorf("expected empty description, got %q", descEmpty)
	}
}

func TestSelectTreeConnector(t *testing.T) {
	t.Parallel()

	if connector := selectTreeConnector(true); connector != constants.TreeCorner {
		t.Errorf("expected corner %q, got %q", constants.TreeCorner, connector)
	}

	if connector := selectTreeConnector(false); connector != constants.TreeBranch {
		t.Errorf("expected branch %q, got %q", constants.TreeBranch, connector)
	}
}

func TestPrintStepBranch(t *testing.T) {
	t.Parallel()

	printStepBranch(constants.TreeBranch, "git status", "(dir: /tmp)")
	printStepBranch(constants.TreeCorner, "go test ./...", "")
}
