package cmdinstall

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/enums/ctxmodetype"
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
	// Direct top-level actions
	"01_show_log":       {"", []string{"show-log"}},
	"02_diff":           {"", []string{"diff"}},
	"03_branch_compare": {"", []string{"show-log", "--branch-compare"}},
	"04_commit":         {"", []string{"commit", "-b", "main"}},
	"05_pull":           {"", []string{"pull"}},
	"06_push":           {"", []string{"push"}},
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
	"40_repos/10_web_ui":   {"", []string{"ui"}},
	"40_repos/20_list":     {"", []string{"list"}},
	"40_repos/30_status":   {"", []string{"status-all"}},
	"40_repos/40_pull_all": {"", []string{"pull-all"}},
	// Visibility
	"50_visibility/10_public":  {"", []string{"make-public"}},
	"50_visibility/20_private": {"", []string{"make-private"}},
	// Install
	"60_install/01_installer_ui": {"", []string{"ui", "installer"}},
	"60_install/05_install_ctx":  {"", []string{"install", "context-menu"}},
	"60_install/10_wizard":       {"", []string{"install", "--wizard"}},
	"60_install/20_dev":          {"", []string{"install", "profile", "dev"}},
	"60_install/30_tools":        {"", []string{"install", "list"}},
	"60_install/40_agy":          {"", []string{"install", "antigravity"}},
	"60_install/50_clean_pwsh":   {"", []string{"context-menu", "clean-powershell"}},
	// Macro
	"65_macro/10_macro_ui":   {"", []string{"ui", "macro"}},
	"65_macro/20_macro_list": {"", []string{"macro", "list"}},
	// Top-level standalone entries
	"70_settings":    {"", []string{"ui", "settings"}},
	"80_antigravity": {"", []string{"ag"}},
	"90_terminal":    {"", nil}, // Prefill mode; no argv
	"91_docs":        {"", []string{"docs"}},
	"92_help":        {"", []string{constants.CmdHelp}},
}

// flattenCtxMenu walks ctxMenu() and returns a path→entry map keyed by the
// "/"-joined KeyName chain. Categories (Children non-nil) are descended into
// but not themselves recorded. Duplicate KeyNames at the same level (e.g. the
// double-listed 90_terminal / 91_docs) are tolerated — the second occurrence
// must match the first.
func flattenCtxMenuByPath(t *testing.T) map[string]ctxEntry {
	t.Helper()
	// WHY: same Windows-runner registry leak as collectCtxLeaves (see
	// installctx_harness_test.go): pin the direct entries visible so
	// the argv table covers the full menu regardless of test order.
	_ = SetGitActionsHidden(false)
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

			prev, isDuplicate := out[path]
			isEqual := reflect.DeepEqual(prev, e)
			if isDuplicate && !isEqual {
				t.Fatalf("duplicate KeyName %q with divergent definition: %+v vs %+v", path, prev, e)
			}

			if isDuplicate {
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
	// 02-spec/04-generic-cli/30-install-ctx.md §3.
	if got := strings.Join(entry.Args, " "); got != "release --bump minor" {
		t.Fatalf("release-next composed argv = %q, want %q", got, "release --bump minor")
	}
}

var expectedCtxFileArgv = map[string][]string{
	"10_diff":       {"diff", `"%1"`},
	"20_diff_later": {"diff", "--stash-path", `"%1"`},
	"30_diff_prev":  {"diff", "HEAD~1", `"%1"`},
	"40_show_log":   {"show-log", `"%1"`},
	"50_blame":      {"blame", `"%1"`},
	"60_commit":     {"commit", "-b", "main"},
}

func TestCtxFileMenuArgvMapping(t *testing.T) {
	entries := ctxFileMenu()
	if len(entries) != len(expectedCtxFileArgv) {
		t.Fatalf("ctxFileMenu count = %d, want %d", len(entries), len(expectedCtxFileArgv))
	}

	for _, e := range entries {
		want, hasEntry := expectedCtxFileArgv[e.KeyName]
		if !hasEntry {
			t.Errorf("unexpected file menu KeyName %q", e.KeyName)

			continue
		}

		if !reflect.DeepEqual(e.Args, want) {
			t.Errorf("%s: Args = %v, want %v", e.KeyName, e.Args, want)
		}

		if e.Mode != ctxmodetype.Terminal {
			t.Errorf("%s: Mode = %v, want %v", e.KeyName, e.Mode, ctxmodetype.Terminal)
		}
	}
}
