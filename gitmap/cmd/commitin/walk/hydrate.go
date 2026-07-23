package walk

import (
	"fmt"
	"strings"
	"time"
)

// commitFormat uses ASCII unit-separators (\x1f) so subjects/bodies
// can carry any byte. Order: AuthorName · AuthorEmail · AuthorDate
// (RFC3339) · CommitterDate (RFC3339) · subject · body.
const commitFormat = "%an\x1f%ae\x1f%aI\x1f%cI\x1f%s\x1f%b"

// hydrate populates one SourceCommit by issuing two git calls:
//  1. `git show -s --format=...` for metadata
//  2. `git show --name-only --format=` for the file list
func hydrate(repoDir, sha string, orderIndex int) (SourceCommit, error) {
	meta, err := readCommitMeta(repoDir, sha)
	if err != nil {
		return SourceCommit{}, err
	}
	files, err := readCommitFiles(repoDir, sha)
	if err != nil {
		return SourceCommit{}, err
	}
	meta.OrderIndex = orderIndex
	meta.Sha = sha
	meta.Files = files
	return meta, nil
}

// readCommitMeta runs the unit-separator-formatted `git show` and
// parses the six fields into a partially-populated SourceCommit.
func readCommitMeta(repoDir, sha string) (SourceCommit, error) {
	out, err := gitRunner(repoDir, "show", "-s", "--format="+commitFormat, sha)
	if err != nil {
		return SourceCommit{}, fmt.Errorf("show %s: %w", sha, err)
	}
	parts := strings.SplitN(out, "\x1f", 6)
	if len(parts) != 6 {
		return SourceCommit{}, fmt.Errorf("malformed git show output for %s: %q", sha, out)
	}
	return assembleMeta(parts)
}

// assembleMeta builds the SourceCommit value from the six split parts.
// Pure function — no I/O.
func assembleMeta(parts []string) (SourceCommit, error) {
	authorDate, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[2]))
	if err != nil {
		return SourceCommit{}, fmt.Errorf("parse author date: %w", err)
	}
	committerDate, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[3]))
	if err != nil {
		return SourceCommit{}, fmt.Errorf("parse committer date: %w", err)
	}
	return SourceCommit{
		AuthorName:      parts[0],
		AuthorEmail:     parts[1],
		AuthorDate:      authorDate,
		CommitterDate:   committerDate,
		OriginalMessage: assembleMessage(parts[4], parts[5]),
	}, nil
}

// assembleMessage rejoins subject + body, trimming trailing whitespace
// so empty bodies don't leave a dangling blank line.
func assembleMessage(subject, body string) string {
	subject = strings.TrimRight(subject, "\n")
	body = strings.TrimRight(body, "\n")
	if body == "" {
		return subject
	}
	return subject + "\n\n" + body
}

// readCommitFiles returns the list of paths touched by sha (POSIX
// separators, no leading `./`). Honors the spec §3 first-parent walk:
// for merges, only the diff against parent #1 is reported.
func readCommitFiles(repoDir, sha string) ([]string, error) {
	out, err := gitRunner(repoDir, "show", "--first-parent", "--name-only", "--format=", sha)
	if err != nil {
		return nil, fmt.Errorf("show files %s: %w", sha, err)
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}
	return splitAndCleanPaths(out), nil
}

// splitAndCleanPaths breaks the newline-delimited list and drops any
// empty / dot-prefixed entries.
func splitAndCleanPaths(out string) []string {
	raw := strings.Split(out, "\n")
	cleaned := make([]string, 0, len(raw))
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		cleaned = append(cleaned, strings.TrimPrefix(p, "./"))
	}
	return cleaned
}
