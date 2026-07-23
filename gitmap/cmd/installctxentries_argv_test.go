package cmd

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// expectedCtxArgv is the authoritative KeyName → argv mapping for every
// right-click context-menu entry. The literal strings here intentionally
// duplicate the resolved values of constants such as constants.CmdRelease,
// constants.FlagBumpDash, and constants.BumpMinor: any drift between the
// constant value and what the menu actually emits will fail this test.
//
// Path-prefixed key (e.g. "30_release/20_release_next") locates an entry
// inside its parent submenu so we can assert nested values without losing
// the order-bearing numeric prefix.
var expectedCtxArgv = map[string]struct {
	exe  string // "" => gitmap binary
	argv []string
}{
	// Scan
	"10_scan/10_scan_here": {"", []string{"scan"}},
	"10_scan/20_rescan":    {"", []string{"rescan"}},
	"10_scan/30_find_next": {"", []string{"find-next"}},
	// Clone
	"20_clone/10_clone_next": {"", []string{"clone-next"}},
	"20_clone/20_pull":       {"", []string{"pull"}},
	"20_clone/30_pull_all":   {"", []string{"pull-all"}},
	// Release — verifies FlagBumpDash + BumpMinor compose to "--bump minor"
	"30_release/10_release":         {"", []string{"release"}},
	"30_release/20_release_next":    {"", []string{"release", "--bump", "minor"}},
	"30_release/30_release_pull":    {"", []string{"pull-release"}},
	"30_release/40_release_pending": {"", []string{"release-pending"}},
	"30_release/50_list_releases":   {"", []string{"list-releases"}},
	"30_release/60_list_versions":   {"", []string{"list-versions"}},
	// Repos
	"40_repos/10_go":     {"", []string{"go-repos"}},
	"40_repos/20_node":   {"", []string{"node-repos"}},
	"40_repos/30_react":  {"", []string{"react-repos"}},
	"40_repos/40_cpp":    {"", []string{"cpp-repos"}},
	"40_repos/50_csharp": {"", []string{"csharp-repos"}},
	// Visibility
	"50_visibility/10_public":  {"", []string{"make-public"}},
	"50_visibility/20_private": {"", []string{"make-private"}},
	// Tools
	"60_tools/10_fix_repo": {"", []string{"fix-repo"}},
	"60_tools/20_diff":     {"", []string{"diff"}},
	"60_tools/30_history":  {"", []string{"history"}},
	"60_tools/40_update":   {"", []string{"update"}},
	// Raw git — Exe override, argv pulled from constants for stability.
	"70_git/10_history": {constants.CtxExeGit, constants.CtxGitHistoryArgs},
	"70_git/20_diff":    {constants.CtxExeGit, constants.CtxGitDiffArgs},
	"70_git/30_log":     {constants.CtxExeGit, constants.CtxGitLogArgs},
	"70_git/40_status":  {constants.CtxExeGit, constants.CtxGitStatusArgs},
	// Top-level standalone entries
	"90_terminal": {"", nil}, // Prefill mode; no argv
	"91_docs":     {"", []string{"docs"}},
	"92_help":     {"", []string{constants.CmdHelp}},
}

// flattenCtxMenu walks ctxMenu() and returns a path→entry map keyed by the
// "/"-joined KeyName chain. Categories (Children non-nil) are descended into
// but not themselves recorded. Duplicate KeyNames at the same level (e.g. the
// double-listed 90_terminal / 91_docs) are tolerated — the second occurrence
// must match the first.
func flattenCtxMenuByPath(t *testing.T) map[string]ctxEntry {
	t.Helper()
	out := map[string]ctxEntry{}
	var walk func(prefix string, entries []ctxEntry)
	walk = func(prefix string, entries []ctxEntry) {
		for _, e := range entries {
			path := e.KeyName
			if prefix != "" {
				path = prefix + "/" + e.KeyName
			}
			if e.Children != nil {
				walk(path, e.Children)
				continue
			}
			if prev, dup := out[path]; dup {
				if !reflect.DeepEqual(prev, e) {
					t.Fatalf("duplicate KeyName %q with divergent definition: %+v vs %+v", path, prev, e)
				}
				continue
			}
			out[path] = e
		}
	}
	walk("", ctxMenu())
	return out
}

func TestCtxMenuKeyNameToArgvMapping(t *testing.T) {
	got := flattenCtxMenuByPath(t)

	for path, want := range expectedCtxArgv {
		entry, ok := got[path]
		if !ok {
			t.Errorf("missing context-menu entry for path %q", path)
			continue
		}
		if entry.Exe != want.exe {
			t.Errorf("%s: Exe = %q, want %q", path, entry.Exe, want.exe)
		}
		if !reflect.DeepEqual(entry.Args, want.argv) {
			t.Errorf("%s: Args = %v, want %v", path, entry.Args, want.argv)
		}
	}

	// Reverse direction: every emitted leaf must be covered by the table so
	// new menu entries don't sneak in without an argv assertion.
	for path := range got {
		if _, ok := expectedCtxArgv[path]; !ok {
			t.Errorf("unexpected context-menu entry %q (Args=%v) — add it to expectedCtxArgv", path, got[path].Args)
		}
	}
}

// TestCtxReleaseNextUsesBumpConstants pins the release-next entry to the
// exact constant identifiers (not just their resolved string values) so a
// rename of FlagBumpDash or BumpMinor that silently swapped to a different
// string would still trip a focused failure with a clear blame target.
func TestCtxReleaseNextUsesBumpConstants(t *testing.T) {
	got := flattenCtxMenuByPath(t)
	entry, ok := got["30_release/20_release_next"]
	if !ok {
		t.Fatal("release-next entry missing")
	}
	want := []string{constants.CmdRelease, constants.FlagBumpDash, constants.BumpMinor}
	if !reflect.DeepEqual(entry.Args, want) {
		t.Fatalf("release-next argv = %v, want %v", entry.Args, want)
	}
	// Sanity: the resolved literals must match the contract documented in
	// spec/04-generic-cli/30-install-ctx.md §3.
	if got := strings.Join(entry.Args, " "); got != "release --bump minor" {
		t.Fatalf("release-next composed argv = %q, want %q", got, "release --bump minor")
	}
}
