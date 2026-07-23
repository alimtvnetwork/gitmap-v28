package runner

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/drift"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/gitlog"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/group"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/render"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/writer"
)

// Execute runs the chosen mode and returns the process exit code.
func Execute(a Args, warnOut io.Writer) (int, error) {
	commits, lower, err := collectCommits(a)
	if err != nil {
		return 0, fmt.Errorf("collecting commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Fprintf(warnOut, "no commits in %s..%s — nothing to do.\n",
			displayRev(lower), displayRev(a.ReleaseTag))

		return 0, nil
	}

	groups := classify(commits, warnOut)
	if len(groups) == 0 {
		return 0, nil
	}

	entry := render.Entry{
		Version: resolveVersion(a, lower),
		Date:    time.Now().UTC().Format("2006-01-02"),
		Groups:  groups,
	}

	return dispatchMode(a.Mode, entry, a.RepoRoot)
}

// collectCommits picks the explicit `--since`..`--release-tag` path
// when either bound is set, otherwise falls back to the auto-detected
// "since latest semver tag" path. Returns the resolved lower bound so
// the version-label fallback (`<lower>+next`) stays accurate.
func collectCommits(a Args) ([]gitlog.Commit, string, error) {
	if a.Since != "" || a.ReleaseTag != "" {
		commits, err := gitlog.CommitsInRange(a.RepoRoot, a.Since, a.ReleaseTag)

		return commits, a.Since, err
	}

	return gitlog.CommitsSinceLastTag(a.RepoRoot)
}

func classify(commits []gitlog.Commit, warnOut io.Writer) []group.Section {
	groups, skipped := group.ByPrefix(commits)
	for _, s := range skipped {
		fmt.Fprintf(warnOut, "skipped (no Conventional Commit prefix): %s %s\n", s.Hash, s.Subject)
	}

	if len(groups) == 0 {
		fmt.Fprintf(warnOut, "no Conventional Commit entries found — nothing to write.\n")
	}

	return groups
}

func dispatchMode(mode Mode, entry render.Entry, repoRoot string) (int, error) {
	mdPath := filepath.Join(repoRoot, "CHANGELOG.md")
	tsPath := filepath.Join(repoRoot, "src/data/changelog.ts")

	switch mode {
	case ModeWrite:
		err := writer.PrependBoth(mdPath, tsPath, entry)
		if err != nil {
			return 0, err
		}

		return 0, nil
	case ModeCheck:
		return drift.Check(mdPath, tsPath, entry)
	default:
		return 0, fmt.Errorf("unhandled mode %q", mode)
	}
}

// resolveVersion picks the entry label in priority order:
// explicit --version > explicit --release-tag > `<lowerBound>+next` >
// `vNEXT`. Lets `make changelog RELEASE_TAG=v3.92.0` Just Work without
// also passing --version.
func resolveVersion(a Args, lowerBound string) string {
	if a.Version != "" {
		return a.Version
	}

	if a.ReleaseTag != "" {
		return a.ReleaseTag
	}

	if lowerBound != "" {
		return lowerBound + "+next"
	}

	return "vNEXT"
}

func displayRev(rev string) string {
	if rev == "" {
		return "HEAD"
	}

	return rev
}
