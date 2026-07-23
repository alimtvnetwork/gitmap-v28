package cmd

import (
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ctxFlatLeaf is the per-leaf record the harness uses to drive every
// platform-specific assertion. It captures the post-flatten view of a
// menu entry plus a stable identity (Slug) so tests on Windows can
// correlate with tests on macOS/Linux that consume the flattened tree.
type ctxFlatLeaf struct {
	Path     string // dotted KeyName path: "30_release.20_release_next"
	Slug     string // flatten slug, e.g. "gitmap-release-release-next"
	Label    string // user-visible flatten label, e.g. "gitmap: Release — Release next (bump minor)"
	Args     []string
	Mode     constants.CtxMode
	Exe      string
	Extended bool
}

// collectCtxLeaves walks ctxMenu() once and returns one ctxFlatLeaf per
// terminal entry, in declaration order. Categories (Children non-nil)
// are descended into but never recorded. Duplicate KeyNames at the
// same level (e.g. the documented duplicate 90_terminal / 91_docs in
// installctxentries.go) are tolerated — both occurrences are emitted
// so callers can assert on / dedupe as they need.
func collectCtxLeaves(t *testing.T) []ctxFlatLeaf {
	t.Helper()
	flat := flattenCtxMenu()
	out := make([]ctxFlatLeaf, 0, len(flat))
	pathByLeaf := buildLeafPaths()
	for _, e := range flat {
		out = append(out, ctxFlatLeaf{
			Path:     pathByLeaf[e.Slug],
			Slug:     e.Slug,
			Label:    e.Label,
			Args:     append([]string(nil), e.Args...),
			Mode:     e.Mode,
			Exe:      e.Exe,
			Extended: e.Extended,
		})
	}

	return out
}

// buildLeafPaths returns slug → "parent.child" KeyName path for every
// leaf currently in ctxMenu(). Used to give each ctxFlatLeaf a stable
// identity that is independent of the visible label (which can change
// for UX reasons without breaking the contract).
func buildLeafPaths() map[string]string {
	out := map[string]string{}
	for _, e := range ctxMenu() {
		if len(e.Children) == 0 {
			out[topSlug(e)] = e.KeyName

			continue
		}
		for _, c := range e.Children {
			out[childSlug(e, c)] = e.KeyName + "." + c.KeyName
		}
	}

	return out
}

func topSlug(e ctxEntry) string {
	return slugifyCtx(constants.CtxFlatPrefix + constants.CtxFlatSeparator + e.MUIVerb)
}

func childSlug(parent, child ctxEntry) string {
	return slugifyCtx(constants.CtxFlatPrefix + constants.CtxFlatSeparator + parent.MUIVerb + constants.CtxFlatChildJoiner + child.MUIVerb)
}

// withExplain flips ctxExplainEnabled for the duration of f and
// guarantees the previous value is restored even on test failure.
// Tests must serialize on this — never run subtests that mutate the
// flag in t.Parallel mode without their own guard.
func withExplain(t *testing.T, on bool, f func()) {
	t.Helper()
	prev := ctxExplainEnabled
	ctxExplainEnabled = on
	defer func() { ctxExplainEnabled = prev }()
	f()
}

// resolvedTarget returns the executable path the platform templates
// will bake into the menu entry: the per-entry Exe override (e.g.
// "git") if set, otherwise the gitmap binary path the harness pinned.
func (l ctxFlatLeaf) resolvedTarget(gitmapExe string) string {
	if l.Exe != "" {
		return l.Exe
	}

	return gitmapExe
}

// fakeGitmapExe returns a deterministic, OS-appropriate path string
// the harness uses in place of os.Executable() so generated registry
// values / .desktop Exec= lines are byte-stable across hosts.
func fakeGitmapExe(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		return `C:\fake\gitmap.exe`
	}

	return filepath.Join("/fake/bin", "gitmap")
}

// sortedSlugs returns every leaf slug, sorted, for set-equality
// assertions across platforms (parity test in step 5).
func sortedSlugs(leaves []ctxFlatLeaf) []string {
	out := make([]string, 0, len(leaves))
	for _, l := range leaves {
		out = append(out, l.Slug)
	}
	sort.Strings(out)

	return out
}

// containsAll asserts every needle is a substring of haystack and
// returns the first miss for a focused failure message.
func containsAll(haystack string, needles []string) (string, bool) {
	for _, n := range needles {
		if !strings.Contains(haystack, n) {
			return n, false
		}
	}

	return "", true
}
