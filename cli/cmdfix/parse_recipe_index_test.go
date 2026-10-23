package cmdfix

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestParseRecipeIndex(t *testing.T) {
	recipes := []gitutil.RemediationRecipe{
		{Title: "Opt 1"},
		{Title: "Opt 2"},
		{Title: "Opt 3"},
	}

	testRecipeMap(t, recipes, map[string]int{
		"1": 0, "stash": 0, "s": 0,
		"2": 1, "wip": 1, "w": 1,
		"3": 2, "discard": 2, "clean": 2, "d": 2,
		"unknown": -1,
	})
}

func testRecipeMap(t *testing.T, recipes []gitutil.RemediationRecipe, cases map[string]int) {
	t.Helper()
	for input, expected := range cases {
		idx := parseRecipeIndex(input, recipes)
		if idx != expected {
			t.Errorf("parseRecipeIndex(%q) = %d, want %d", input, idx, expected)
		}
	}
}

func TestResolveRecipeIdx(t *testing.T) {
	recipes := []gitutil.RemediationRecipe{
		{Title: "Option 1: Stash"},
		{Title: "Option 2: Commit WIP"},
	}

	if idx := resolveRecipeIdx("1", recipes); idx != 0 {
		t.Fatalf("expected 0, got %d", idx)
	}
	if idx := resolveRecipeIdx("2", recipes); idx != 1 {
		t.Fatalf("expected 1, got %d", idx)
	}
	if idx := resolveRecipeIdx("unknown", recipes); idx != 0 {
		t.Fatalf("expected 0 for fallback, got %d", idx)
	}
	if idx := resolveRecipeIdx("1", nil); idx != 0 {
		t.Fatalf("expected 0 for empty recipes, got %d", idx)
	}
}
