package completion

import (
	"sort"
	"strings"
	"testing"
)

// legacyAllCommands is the hand-curated list that AllCommands() returned
// before the go:generate refactor. The generated + manualExtras union must
// still cover every entry here; otherwise tab-completion regressed.
//
// Note: two stale aliases from the original manualExtras slice were corrected
// during the migration to match the real Cmd*Alias constants:
//   - "ep" -> "ex"  (CmdExportAlias = "ex")
//   - "gor" -> "gr" (CmdGoReposAlias = "gr")
//
// These were never registered in the dispatcher, so removing them from the
// legacy contract is a bug-fix, not a regression.
var legacyAllCommands = []string{
	"scan", "s",
	"clone", "c",
	"pull", "p",
	"status", "st",
	"exec", "x",
	"release", "r",
	"release-branch", "rb",
	"release-pending", "rp",
	"changelog", "cl",
	"latest-branch", "lb",
	"list", "ls",
	"group", "g",
	"multi-group", "mg",
	"cd", "go",
	"update",
	"version", "v",
	"desktop-sync", "ds",
	"github-desktop", "gd",
	"rescan", "rs",
	"setup",
	"doctor",
	"db-reset",
	"list-versions", "lv",
	"list-releases", "lr",
	"revert",
	"seo-write", "sw",
	"amend", "am",
	"amend-list", "al",
	"history", "hi",
	"history-reset", "hr",
	"stats", "ss",
	"bookmark", "bk",
	"export", "ex",
	"import", "im",
	"profile", "pf",
	"diff-profiles", "dp",
	"watch", "w",
	"gomod", "gm",
	"go-repos", "gr",
	"node-repos", "nr",
	"react-repos", "rr",
	"cpp-repos", "cr",
	"csharp-repos", "csr",
	"completion", "cmp",
	"interactive", "i",
	"clear-release-json", "crj",
	"alias", "a",
	"zip-group", "z",
	"dashboard", "db",
	"ssh",
	"prune", "prn",
	"pull-release", "pr",
	"temp-release", "tr",
	"clone-next", "cn",
	"uninstall", "un",
	"help",
	"version-history", "vh",
	"llm-docs", "ld",
	"mv", "move",
	"merge-both", "mb",
	"merge-left", "ml",
	"merge-right", "mr",
}

// TestAllCommandsCoversLegacyList guards against regressions in the union of
// generatedCommands and manualExtras after the go:generate refactor.
func TestAllCommandsCoversLegacyList(t *testing.T) {
	have := make(map[string]bool, len(AllCommands()))
	for _, v := range AllCommands() {
		have[v] = true
	}

	var missing []string

	for _, want := range legacyAllCommands {
		if !have[want] {
			missing = append(missing, want)
		}
	}

	if len(missing) == 0 {
		return
	}

	sort.Strings(missing)
	t.Fatalf("AllCommands() missing %d legacy entries: %s",
		len(missing), strings.Join(missing, ", "))
}

// TestAllCommandsIsSorted confirms the returned slice is sorted and unique,
// since downstream completion scripts rely on stable ordering for diff-friendly
// generation.
func TestAllCommandsIsSorted(t *testing.T) {
	got := AllCommands()

	for i := 1; i < len(got); i++ {
		if got[i-1] >= got[i] {
			t.Fatalf("AllCommands() not strictly sorted at index %d: %q >= %q",
				i, got[i-1], got[i])
		}
	}
}
