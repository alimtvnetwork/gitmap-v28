package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/helptext"
)

// pullReleaseAliasNames lists every spelling that must route to the
// pull-release handler and render the canonical pull-release help page.
// Derived directly from the constants so a future alias addition forces
// a test update at the same site as the constant.
var pullReleaseAliasNames = []string{
	constants.CmdReleasePull,       // pull-release (canonical)
	constants.CmdReleasePullAlias,  // pr
	constants.CmdReleasePullAlias2, // release-pull (legacy)
	constants.CmdReleasePullAlias3, // relp (legacy)
	constants.CmdReleasePullAlias4, // rlp (legacy)
}

// TestPullReleaseAliasesShareOneDispatchEntry guards that every alias
// is registered in a single dispatchEntry — split entries would let
// handlers drift apart silently.
func TestPullReleaseAliasesShareOneDispatchEntry(t *testing.T) {
	entry := findPullReleaseDispatchEntry(t)

	for _, alias := range pullReleaseAliasNames {
		if !matchAny(alias, entry.names) {
			t.Fatalf("alias %q missing from pull-release dispatch entry %v", alias, entry.names)
		}
	}
	if len(entry.names) != len(pullReleaseAliasNames) {
		t.Fatalf("pull-release dispatch entry has %d names, expected %d (%v)",
			len(entry.names), len(pullReleaseAliasNames), entry.names)
	}
}

// TestPullReleaseAliasesAreUnique guards against accidental duplication
// (e.g. a legacy alias colliding with the canonical name).
func TestPullReleaseAliasesAreUnique(t *testing.T) {
	seen := make(map[string]struct{}, len(pullReleaseAliasNames))
	for _, alias := range pullReleaseAliasNames {
		if alias == "" {
			t.Fatalf("empty alias in pullReleaseAliasNames")
		}
		if _, dup := seen[alias]; dup {
			t.Fatalf("duplicate pull-release alias: %q", alias)
		}
		seen[alias] = struct{}{}
	}
}

// TestPullReleaseAliasesResolveSameHelpPage guards that --help for any
// alias renders the canonical pull-release.md page. checkHelp calls
// helptext.PrintWithMode with the dispatcher-supplied command name; the
// dispatcher passes constants.CmdReleasePull for every alias because
// they share one handler closure. We assert the help file resolves and
// looks like pull-release help, which transitively proves every alias
// surfaces identical help output.
func TestPullReleaseAliasesResolveSameHelpPage(t *testing.T) {
	raw := captureHelp(t, constants.CmdReleasePull)
	if !strings.Contains(raw, "pull-release") {
		t.Fatalf("canonical help page missing 'pull-release' marker:\n%s", raw)
	}

	for _, alias := range pullReleaseAliasNames {
		if alias == constants.CmdReleasePull {
			continue
		}
		// Every alias is also documented in the canonical help page's
		// alias table, so the same help text serves every spelling.
		if !strings.Contains(raw, alias) {
			t.Fatalf("canonical pull-release help page omits alias %q", alias)
		}
	}
}

// findPullReleaseDispatchEntry locates the dispatchEntry that owns the
// canonical pull-release command name.
func findPullReleaseDispatchEntry(t *testing.T) dispatchEntry {
	t.Helper()
	for _, entry := range releaseDispatchEntries() {
		if matchAny(constants.CmdReleasePull, entry.names) {
			return entry
		}
	}
	t.Fatalf("pull-release dispatch entry not found")
	return dispatchEntry{}
}

// captureHelp reads the embedded help markdown via helptext.PrintRaw's
// underlying loader without invoking os.Exit. We assert against the raw
// markdown so the test is independent of pretty-render state.
func captureHelp(t *testing.T, command string) string {
	t.Helper()
	data, err := helptext.ReadRaw(command)
	if err != nil {
		t.Fatalf("helptext.ReadRaw(%q) failed: %v", command, err)
	}
	return string(data)
}
