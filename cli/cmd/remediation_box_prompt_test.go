package cmd

import (
	"testing"
)

func TestRemediationChoiceMatchers(t *testing.T) {
	allChoices := []string{"a", "1", "all", "y", "yes", "A", "YES"}
	for _, c := range allChoices {
		if !isFixAllChoice(c) {
			t.Errorf("expected %q to be recognized as fix-all choice", c)
		}
		if isFixSingleChoice(c) {
			t.Errorf("did not expect %q to be recognized as fix-single choice", c)
		}
	}

	singleChoices := []string{"s", "2", "single", "step", "S", "STEP"}
	for _, c := range singleChoices {
		if !isFixSingleChoice(c) {
			t.Errorf("expected %q to be recognized as fix-single choice", c)
		}
		if isFixAllChoice(c) {
			t.Errorf("did not expect %q to be recognized as fix-all choice", c)
		}
	}

	skipChoices := []string{"k", "q", "skip", "n", "no", "exit", "", "other"}
	for _, c := range skipChoices {
		if isFixAllChoice(c) {
			t.Errorf("did not expect %q to be fix-all choice", c)
		}
		if isFixSingleChoice(c) {
			t.Errorf("did not expect %q to be fix-single choice", c)
		}
	}
}
