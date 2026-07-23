package gitlog

import "testing"

// TestSortCommitsByDateThenHash pins the (author-date ASC, hash ASC)
// ordering contract that lets two reruns produce byte-identical output
// even when git emits commits in different orders (rebase, replay, …).
func TestSortCommitsByDateThenHash(t *testing.T) {
	in := []Commit{
		{Hash: "ccc", AuthorUnix: 200, Subject: "feat: c"},
		{Hash: "aaa", AuthorUnix: 100, Subject: "feat: a"},
		{Hash: "bbb", AuthorUnix: 100, Subject: "feat: b"},
		{Hash: "ddd", AuthorUnix: 300, Subject: "feat: d"},
	}

	sortCommits(in)

	want := []string{"aaa", "bbb", "ccc", "ddd"}
	for i, w := range want {
		if in[i].Hash != w {
			t.Fatalf("position %d: want %s got %s", i, w, in[i].Hash)
		}
	}
}

// TestParseCommitLineExtractsAuthorTimestamp guarantees the gitlog
// reader populates AuthorUnix so the sort key is never zero in
// production. A regression here would silently collapse the sort to
// hash-only and break chronological ordering.
func TestParseCommitLineExtractsAuthorTimestamp(t *testing.T) {
	c, ok := parseCommitLine("abc123\x1f1700000000\x1ffeat: hello")
	if !ok {
		t.Fatalf("parseCommitLine rejected a well-formed line")
	}

	if c.Hash != "abc123" || c.AuthorUnix != 1700000000 || c.Subject != "feat: hello" {
		t.Fatalf("parsed wrong fields: %+v", c)
	}
}
