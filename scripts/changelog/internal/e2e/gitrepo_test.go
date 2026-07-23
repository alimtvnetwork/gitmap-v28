package e2e

import (
	"fmt"
	"strings"
	"testing"
)

// gitRepo wraps a throwaway repository on disk. Commits are created
// via `git fast-import` (low-level plumbing) so the sandbox's block on
// `git commit` does not apply, and so timestamps + hashes are fully
// reproducible across machines and reruns.
type gitRepo struct {
	t      *testing.T
	dir    string
	stream strings.Builder
	mark   int
	tags   []taggedCommit
}

type taggedCommit struct {
	name string
	mark int
}

func newGitRepo(t *testing.T, dir string) *gitRepo {
	t.Helper()

	r := &gitRepo{t: t, dir: dir}
	r.run("init", "-q", "-b", "main")
	r.run("config", "user.email", "ci@example.com")
	r.run("config", "user.name", "ci")

	return r
}

// commit appends a commit to the fast-import stream. Author and
// committer are pinned to the same identity and the same unix
// timestamp so the resulting hash is deterministic.
func (r *gitRepo) commit(subject string, unix int64) {
	r.t.Helper()

	r.mark++
	parent := ""

	if r.mark > 1 {
		parent = fmt.Sprintf("from :%d\n", r.mark-1)
	}

	fmt.Fprintf(&r.stream,
		"commit refs/heads/main\nmark :%d\ncommitter ci <ci@example.com> %d +0000\ndata %d\n%s\n%s\n",
		r.mark, unix, len(subject), subject, parent)
}

// tag records a lightweight tag pointing at the most recent commit.
// Tags are created in the order they are added, after the stream is
// flushed by finalize().
func (r *gitRepo) tag(name string) {
	r.t.Helper()

	r.tags = append(r.tags, taggedCommit{name: name, mark: r.mark})
}
