package clonefrom

// depthflag_format_test.go — locks the EXACT spelling of clone-from's
// shallow-clone flag in BOTH the executed argv (BuildGitArgs, which
// shells out to git) and the human-facing preview (cloneCommandForRow,
// which feeds the dry-run + the --output terminal block). The two MUST
// stay byte-identical so the printed cmd: line is faithful and the
// --verify-cmd-faithful checker has zero false positives.
//
// If you intentionally switch to the split form (`--depth N`), you must
// update constants.CloneFromDepthFlagFmt + the golden fixture
// cmd/testdata/clonetermblock_clonefrom.golden in the same commit and
// then re-run the goldens with -update. This test will fail loudly until
// then, by design.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestDepthFlagFormat_Locked pins the constant itself. A typo here
// (e.g. dropping the `=`) would silently change every render site;
// this guard catches that at the constant level before any caller
// runs.
func TestDepthFlagFormat_Locked(t *testing.T) {
	got := fmt.Sprintf(constants.CloneFromDepthFlagFmt, 1)
	const want = "--depth=1"
	if got != want {
		t.Fatalf("CloneFromDepthFlagFmt rendered %q, want %q "+
			"(joined form is mandatory — see constant doc)", got, want)
	}
	if strings.Contains(constants.CloneFromDepthFlagFmt, " ") {
		t.Fatalf("CloneFromDepthFlagFmt contains a space (%q); "+
			"split form `--depth N` is forbidden for clone-from",
			constants.CloneFromDepthFlagFmt)
	}
}

// TestBuildGitArgs_DepthJoined asserts the executor passes a SINGLE
// joined token (`--depth=N`) — not two tokens (`--depth`, `N`). git
// accepts both, but the printed cmd: line and exec.Command argv must
// agree token-for-token or --verify-cmd-faithful flags a false drift.
func TestBuildGitArgs_DepthJoined(t *testing.T) {
	row := Row{URL: "https://example.com/x.git", Branch: "main", Depth: 5}
	args := BuildGitArgs(row, "x")

	const wantTok = "--depth=5"
	if !containsTok(args, wantTok) {
		t.Fatalf("BuildGitArgs argv missing %q\n got: %v",
			wantTok, args)
	}
	// Reject the split form explicitly: `--depth` followed by a bare
	// number must NEVER appear in this executor's argv.
	for i, a := range args {
		if a == "--depth" {
			t.Fatalf("BuildGitArgs argv contains split form "+
				"`--depth` at index %d (followed by %q) — must "+
				"be joined `--depth=N`\n got: %v",
				i, safeIdx(args, i+1), args)
		}
	}
}

// TestCloneCommandForRow_DepthJoined mirrors the executor check on the
// preview side: the string the user SEES in the terminal block must
// contain `--depth=N`, not `--depth N`.
func TestCloneCommandForRow_DepthJoined(t *testing.T) {
	row := Row{URL: "https://example.com/x.git", Branch: "main", Depth: 7}
	got := cloneCommandForRow(row, "x")

	if !strings.Contains(got, "--depth=7") {
		t.Fatalf("cloneCommandForRow missing `--depth=7`\n got: %s",
			got)
	}
	if strings.Contains(got, "--depth 7") {
		t.Fatalf("cloneCommandForRow rendered split form "+
			"`--depth 7` — must be joined `--depth=7`\n got: %s",
			got)
	}
}

// containsTok is a tiny helper so the assertions above read top-down.
// Kept package-private to this _test.go (no production callers).
func containsTok(args []string, tok string) bool {
	for _, a := range args {
		if a == tok {
			return true
		}
	}

	return false
}

// safeIdx returns args[i] or "<end-of-argv>" so the failure message
// in TestBuildGitArgs_DepthJoined is informative even when `--depth`
// is the last token (which would itself be a separate bug).
func safeIdx(args []string, i int) string {
	if i < 0 || i >= len(args) {
		return "<end-of-argv>"
	}

	return args[i]
}
