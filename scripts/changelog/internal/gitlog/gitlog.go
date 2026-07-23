// Package gitlog wraps the small subset of `git` we need: the highest
// semver tag in the repo and the commits between that tag and HEAD.
//
// Ordering contract (deterministic across machines and reruns):
//
//   - Tags are sorted by `--sort=-v:refname` (semver descending) so the
//     "latest" tag is the highest version number, never whichever tag
//     happens to be most-recently reachable from HEAD. Two tags on
//     diverged branches can no longer flip the chosen boundary.
//   - Commits are read with `--date-order` and then re-sorted in Go by
//     (author-unix-timestamp ASC, hash ASC). Author date is stable
//     across rebases/cherry-picks where committer date is not, and the
//     hash tiebreak keeps the order total when two commits share a
//     timestamp.
package gitlog

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// Commit is the minimal commit shape consumed by group.ByPrefix.
type Commit struct {
	Hash       string
	Subject    string
	AuthorUnix int64
}

// CommitsSinceLastTag returns commits between the highest-semver tag and
// HEAD in deterministic order. When no tag exists, every commit
// reachable from HEAD is returned and lastTag is "".
func CommitsSinceLastTag(repoRoot string) ([]Commit, string, error) {
	tag, err := highestSemverTag(repoRoot)
	if err != nil {
		return nil, "", err
	}

	commits, err := CommitsInRange(repoRoot, tag, "")
	if err != nil {
		return nil, "", err
	}

	return commits, tag, nil
}

// CommitsInRange returns commits in `since..until` order. Empty `since`
// means "from the repository start"; empty `until` means HEAD. Used by
// the runner when --since / --release-tag explicitly bound the range.
func CommitsInRange(repoRoot, since, until string) ([]Commit, error) {
	rev := buildRev(since, until)

	commits, err := readCommits(repoRoot, rev)
	if err != nil {
		return nil, err
	}

	sortCommits(commits)

	return commits, nil
}

// buildRev assembles the `<since>..<until>` token git log expects,
// collapsing each empty bound to its sane default (start of history,
// HEAD).
func buildRev(since, until string) string {
	if until == "" {
		until = "HEAD"
	}

	if since == "" {
		return until
	}

	return since + ".." + until
}

// highestSemverTag returns the tag with the highest semver, not the
// most-recently-reachable tag. `git tag --sort=-v:refname` yields the
// list in descending version order; we take the first entry.
func highestSemverTag(repoRoot string) (string, error) {
	cmd := exec.Command("git", "-C", repoRoot, "tag", "--sort=-v:refname")

	out, err := cmd.Output()
	if err != nil {
		_, isExitErr := err.(*exec.ExitError)
		if isExitErr {
			return "", nil
		}

		return "", fmt.Errorf("git tag failed: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed, nil
		}
	}

	return "", nil
}

func readCommits(repoRoot, rev string) ([]Commit, error) {
	cmd := exec.Command(
		"git", "-C", repoRoot, "log",
		"--date-order",
		"--pretty=format:%H\x1f%at\x1f%s",
		rev,
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}

	return parseLog(string(out)), nil
}

func parseLog(raw string) []Commit {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}

	lines := strings.Split(trimmed, "\n")
	commits := make([]Commit, 0, len(lines))

	for _, line := range lines {
		c, ok := parseCommitLine(line)
		if ok {
			commits = append(commits, c)
		}
	}

	return commits
}

func parseCommitLine(line string) (Commit, bool) {
	parts := strings.SplitN(line, "\x1f", 3)
	if len(parts) != 3 {
		return Commit{}, false
	}

	ts := parseUnix(parts[1])

	return Commit{Hash: parts[0], AuthorUnix: ts, Subject: parts[2]}, true
}

func parseUnix(s string) int64 {
	var n int64

	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	if err != nil {
		return 0
	}

	return n
}

// sortCommits enforces (author-date ASC, hash ASC) so the same input
// produces the same output regardless of how git happened to emit it.
func sortCommits(commits []Commit) {
	sort.SliceStable(commits, func(i, j int) bool {
		if commits[i].AuthorUnix != commits[j].AuthorUnix {
			return commits[i].AuthorUnix < commits[j].AuthorUnix
		}

		return commits[i].Hash < commits[j].Hash
	})
}
