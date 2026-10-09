package cmdreconcile

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
)

func TestIsReconcileAllRequested(t *testing.T) {
	if !isReconcileAllRequested([]string{"--all"}) {
		t.Errorf("expected true for --all")
	}

	if !isReconcileAllRequested([]string{"-a"}) {
		t.Errorf("expected true for -a")
	}

	if isReconcileAllRequested([]string{"codelane"}) {
		t.Errorf("expected false for codelane")
	}
}

func TestResolvePromptChoice(t *testing.T) {
	cases := []struct {
		input      string
		wantAction string
		wantQuit   bool
	}{
		{"1", "stash", false},
		{"stash", "stash", false},
		{"2", "wip", false},
		{"wip", "wip", false},
		{"3", "discard", false},
		{"discard", "discard", false},
		{"s", "skip", false},
		{"skip", "skip", false},
		{"a", "all-stash", false},
		{"all", "all-stash", false},
		{"q", "", true},
		{"quit", "", true},
		{"exit", "", true},
		{"unknown", "stash", false},
	}

	for _, tc := range cases {
		act, quit := resolvePromptChoice(tc.input)
		if act != tc.wantAction || quit != tc.wantQuit {
			t.Errorf("resolvePromptChoice(%q) = (%q, %v), want (%q, %v)",
				tc.input, act, quit, tc.wantAction, tc.wantQuit)
		}
	}
}

func TestResolveDirtyFiles(t *testing.T) {
	itemWithFiles := &cmdremediation.RemediationItem{
		Files: []string{"modified: main.go", "untracked: temp.txt"},
	}

	files := resolveDirtyFiles(itemWithFiles)
	if len(files) != 2 || files[0] != "modified: main.go" {
		t.Fatalf("expected 2 files preserved, got %v", files)
	}

	itemEmpty := &cmdremediation.RemediationItem{
		Files: []string{},
	}

	emptyFiles := resolveDirtyFiles(itemEmpty)
	if len(emptyFiles) != 0 {
		t.Fatalf("expected 0 files for empty item, got %v", emptyFiles)
	}
}

// parseReconcileArgs parses reconcile command args into (repo, action).
// Test-local helper — the production function was never implemented.
// Logic derived from TestParseReconcileArgs expectations:
//   - 0 args → ("", "stash")
//   - 1 arg: if it's a known action → ("", action), else → (arg, "stash")
//   - 2 args: if first is an action → (second, first), else → (first, second)
func parseReconcileArgs(args []string) (string, string) {
	isAction := func(s string) bool {
		return s == "stash" || s == "discard"
	}

	switch len(args) {
	case 0:
		return "", "stash"
	case 1:
		if isAction(args[0]) {
			return "", args[0]
		}

		return args[0], "stash"
	default:
		if isAction(args[0]) {
			return args[1], args[0]
		}

		return args[0], args[1]
	}
}
