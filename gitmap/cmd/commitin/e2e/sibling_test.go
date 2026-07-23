package e2e

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestKeywordAllDiscoversEveryVersionedSibling builds three sibling
// repos (project, project-v2, project-v3) in one parent and runs
// commit-in with the `all` keyword. All three siblings of `project`
// (= project itself as v0 + v2 + v3) — minus the source — must be
// staged and replayed in ascending-version order.
//
// Note: sibling discovery EXCLUDES the source itself, so the source
// repo `project` is NOT replayed into itself; only project-v2 and
// project-v3 are. We seed each sibling with one easily-distinguished
// commit so we can assert order via subjects.
func TestKeywordAllDiscoversEveryVersionedSibling(t *testing.T) {
	parent, src := buildSiblingFixture(t,
		siblingSpec{name: "project-v2", subject: "v2-only"},
		siblingSpec{name: "project-v3", subject: "v3-only"},
	)
	_ = parent
	raw := NewRawArgs(src.Path)
	raw.Keyword = constants.CommitInInputKeywordAll

	res := Run(t, raw)
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d, want 0\nstderr=%s", res.ExitCode, res.Stderr)
	}
	dst := src.LogFirstParent(t)
	if len(dst) != 2 {
		t.Fatalf("dst commits=%d, want 2\nstderr=%s", len(dst), res.Stderr)
	}
	// Ascending-version order means v2 replays before v3.
	if dst[0].Subject != "v2-only" || dst[1].Subject != "v3-only" {
		t.Fatalf("subjects=%q,%q, want v2-only, v3-only",
			dst[0].Subject, dst[1].Subject)
	}
}

// TestKeywordTailNReturnsLastNSiblings builds four siblings and runs
// commit-in with `-2` — only the two most-recent (highest-version)
// siblings should be staged. Per spec §2.4, when fewer than N
// siblings are present `-N` returns all of them; we use N=2 with
// 3 available siblings to exercise the slicing branch.
func TestKeywordTailNReturnsLastNSiblings(t *testing.T) {
	parent, src := buildSiblingFixture(t,
		siblingSpec{name: "project-v2", subject: "v2"},
		siblingSpec{name: "project-v3", subject: "v3"},
		siblingSpec{name: "project-v4", subject: "v4"},
	)
	_ = parent
	raw := NewRawArgs(src.Path)
	raw.Keyword = "-2"
	raw.KeywordTail = 2

	res := Run(t, raw)
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d, want 0\nstderr=%s", res.ExitCode, res.Stderr)
	}
	dst := src.LogFirstParent(t)
	if len(dst) != 2 {
		t.Fatalf("dst commits=%d, want 2 (last-2 of 3 siblings)\nstderr=%s",
			len(dst), res.Stderr)
	}
	if dst[0].Subject != "v3" || dst[1].Subject != "v4" {
		t.Fatalf("subjects=%q,%q, want v3,v4", dst[0].Subject, dst[1].Subject)
	}
}

// TestAutoInitOnMissingSourceCreatesRepoAndReplays exercises spec §2.3
// case 4 (missing path → mkdir + git init). We point --source at a
// path that does NOT exist yet; orchestrator must create + init it,
// then replay the input's commits into it.
func TestAutoInitOnMissingSourceCreatesRepoAndReplays(t *testing.T) {
	tmp := t.TempDir()
	missing := filepath.Join(tmp, "fresh-source")
	input := NewRepoIn(t, tmp, "input")
	input.Commit("a.txt", "1\n", "seed", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	raw := NewRawArgs(missing, input.Path)
	res := Run(t, raw)
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d, want 0\nstderr=%s", res.ExitCode, res.Stderr)
	}
	// Wrap the now-init'd source in a Repo so we can re-use the log
	// helpers. NewRepoIn would re-init, so build the value manually.
	created := &Repo{t: t, Path: missing}
	if log := created.LogFirstParent(t); len(log) != 1 || log[0].Subject != "seed" {
		t.Fatalf("auto-init source did not receive the seed commit: %+v", log)
	}
	if !strings.Contains(res.Stderr, "fresh-source") {
		t.Errorf("stderr should mention freshly-init'd source path\nstderr=%s", res.Stderr)
	}
}

// siblingSpec describes one versioned-sibling input for the helper.
type siblingSpec struct {
	name    string // folder name, e.g. "project-v2"
	subject string // commit subject so tests can assert order
}

// buildSiblingFixture creates `source` ("project") + every spec in
// `sibs` inside a single shared parent directory so the orchestrator's
// sibling discovery (parent-of-source scope) finds them.
func buildSiblingFixture(t *testing.T, sibs ...siblingSpec) (parent string, src *Repo) {
	t.Helper()
	parent = t.TempDir()
	src = NewRepoIn(t, parent, "project")
	for i, s := range sibs {
		repo := NewRepoIn(t, parent, s.name)
		when := time.Date(2024, 1, 1+i, 0, 0, 0, 0, time.UTC)
		repo.Commit("f.txt", s.subject+"\n", s.subject, when)
	}
	return parent, src
}
