package gitutil

import (
	"strings"
	"testing"
)

func TestRemediationStepsGeneration(t *testing.T) {
	testPath := "D:/My Work/test-repo"
	commitRecipe := GenerateCommitRecipe(testPath)

	if len(commitRecipe.Steps) != 3 {
		t.Fatalf("expected 3 steps in commit recipe, got %d", len(commitRecipe.Steps))
	}
	if commitRecipe.Steps[0].Name != "git" {
		t.Errorf("expected step 0 name 'git', got %s", commitRecipe.Steps[0].Name)
	}

	// Verify that rawPath in Steps does NOT contain quotes, even when path has spaces
	for _, step := range commitRecipe.Steps {
		isCDir := len(step.Args) >= 2 && step.Args[0] == "-C"
		hasQuotes := isCDir && (strings.HasPrefix(step.Args[1], `"`) || strings.HasSuffix(step.Args[1], `"`))
		if hasQuotes {
			t.Errorf("rawPath in step args must not have quotes: got %s", step.Args[1])
		}
	}

	// Verify commit step has exact message as single argument
	commitStep := commitRecipe.Steps[1]
	hasMsgFlag := false
	for i, arg := range commitStep.Args {
		isMatchingMsg := arg == "-m" && i+1 < len(commitStep.Args) && commitStep.Args[i+1] == "wip: local changes"
		if isMatchingMsg {
			hasMsgFlag = true
		}
	}
	if !hasMsgFlag {
		t.Errorf("expected -m flag with 'wip: local changes' in commit step: %v", commitStep.Args)
	}

	stashRecipe := GenerateStashRecipe(testPath)
	if len(stashRecipe.Steps) != 3 {
		t.Errorf("expected 3 steps in stash recipe, got %d", len(stashRecipe.Steps))
	}

	discardRecipe := GenerateDiscardRecipe(testPath)
	if len(discardRecipe.Steps) != 3 {
		t.Errorf("expected 3 steps in discard recipe, got %d", len(discardRecipe.Steps))
	}
}
