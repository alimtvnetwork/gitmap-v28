package walk

import (
	"fmt"
	"strings"
	"time"
)

// SourceCommit is one fully-hydrated commit ready for downstream
// stages (DedupeCheck, BuildFileSet, BuildMessage, Commit).
// Pure value type — zero git, zero filesystem after construction.
type SourceCommit struct {
	OrderIndex      int       // 1-based, oldest=1
	Sha             string    // 40-char full SHA
	AuthorName      string    // "Name" portion of the source author
	AuthorEmail     string    // "email@host" portion of the source author
	AuthorDate      time.Time // source GIT_AUTHOR_DATE
	CommitterDate   time.Time // source GIT_COMMITTER_DATE
	OriginalMessage string    // full message: subject + blank + body
	Files           []string  // POSIX relative paths touched by this commit
}

// WalkFirstParent lists every commit reachable from HEAD via the first
// parent only, oldest→newest, then hydrates each one. Maps to spec §3.1
// stage 09 (`WalkCommits`) with the determinism rule from §3.4.
//
// `repoDir` is the staged input path (a real git working tree). The
// caller is responsible for staging clones before this runs.
func WalkFirstParent(repoDir string) ([]SourceCommit, error) {
	shas, err := listFirstParentShas(repoDir)
	if err != nil {
		return nil, fmt.Errorf("walk: list shas: %w", err)
	}
	if len(shas) == 0 {
		return nil, nil
	}
	out := make([]SourceCommit, 0, len(shas))
	for i, sha := range shas {
		c, hyErr := hydrate(repoDir, sha, i+1)
		if hyErr != nil {
			return nil, fmt.Errorf("walk: hydrate %s: %w", sha, hyErr)
		}
		out = append(out, c)
	}
	return out, nil
}

// listFirstParentShas runs `git rev-list --first-parent --reverse HEAD`
// in repoDir. Empty repos return an empty slice (no error).
func listFirstParentShas(repoDir string) ([]string, error) {
	out, err := gitRunner(repoDir, "rev-list", "--first-parent", "--reverse", "HEAD")
	if err != nil {
		if isEmptyRepoError(out, err) {
			return nil, nil
		}
		return nil, err
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// isEmptyRepoError matches git's "unknown revision HEAD" diagnostic so
// brand-new repos walk as zero commits instead of erroring out.
func isEmptyRepoError(out string, err error) bool {
	if err == nil {
		return false
	}
	lc := strings.ToLower(out + " " + err.Error())
	return strings.Contains(lc, "unknown revision") || strings.Contains(lc, "ambiguous argument 'head'")
}
