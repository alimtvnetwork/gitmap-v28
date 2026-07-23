package clonenext

// Remote version discovery for `gitmap cn`.
//
// Given a parsed local repo (BaseName + CurrentVersion), this probes
// GitHub for the highest existing -v<M> sibling repo and returns the
// best target. The probe is fail-fast: as soon as -v<N+k> returns 404,
// we stop walking upward (matches install.sh's resolve_effective_repo).
//
// Spec reference: spec/01-app/95-installer-script-find-latest-repo.md
// Issue context: cn batch update detection (v3.43.3).

import (
	"fmt"
)

// DefaultRemoteProbeCeiling caps how far above the current version we
// probe before giving up. Mirrors install.sh's PROBE_CEILING default.
const DefaultRemoteProbeCeiling = 30

// RemoteUpdateCheck is the result of comparing a local repo to its
// versioned remote siblings.
type RemoteUpdateCheck struct {
	LocalVersion  int    // current -v<N> derived from the local folder name
	RemoteVersion int    // highest -v<M> that exists on GitHub (== local when none higher)
	UpdateNeeded  bool   // true iff RemoteVersion > LocalVersion
	TargetRepo    string // "<owner>/<base>-v<RemoteVersion>" — empty when no update
}

// repoExistsFn is the indirection point used by tests to stub the
// network probe without hitting the GitHub API.
type repoExistsFn func(owner, repo string) (bool, error)

// CheckRemoteForUpdate probes GitHub for higher -v<M> siblings of the
// given parsed repo and returns whether a re-clone is warranted.
//
// `owner` is the GitHub owner (org or user) parsed from the local repo's
// origin URL. `parsed` is the result of ParseRepoName on the local
// folder name. `ceiling` is the inclusive upper bound on M; pass
// DefaultRemoteProbeCeiling for the standard cap.
//
// Returns UpdateNeeded=false when the repo has no -v<N> suffix at all
// (cannot reason about "next version" without a baseline).
func CheckRemoteForUpdate(owner string, parsed ParsedRepo, ceiling int) (RemoteUpdateCheck, error) {
	return checkRemoteForUpdateWith(owner, parsed, ceiling, RepoExists)
}

// checkRemoteForUpdateWith is the test seam: identical to
// CheckRemoteForUpdate but with an injectable probe function.
func checkRemoteForUpdateWith(owner string, parsed ParsedRepo, ceiling int, probe repoExistsFn) (RemoteUpdateCheck, error) {
	out := RemoteUpdateCheck{
		LocalVersion:  parsed.CurrentVersion,
		RemoteVersion: parsed.CurrentVersion,
	}
	if !parsed.HasVersion {
		// No -v<N> baseline → no notion of "next version" to check.
		return out, nil
	}
	if ceiling <= 0 {
		ceiling = DefaultRemoteProbeCeiling
	}

	highest, err := probeHighestSibling(owner, parsed, ceiling, probe)
	if err != nil {
		return out, err
	}

	out.RemoteVersion = highest
	if highest > parsed.CurrentVersion {
		out.UpdateNeeded = true
		out.TargetRepo = fmt.Sprintf("%s/%s", owner, TargetRepoName(parsed.BaseName, highest))
	}

	return out, nil
}

// probeHighestSibling walks -v<N+1>, -v<N+2>, ... up to ceiling and
// returns the highest M that exists. Stops at the first 404 (fail-fast).
func probeHighestSibling(owner string, parsed ParsedRepo, ceiling int, probe repoExistsFn) (int, error) {
	highest := parsed.CurrentVersion
	for m := parsed.CurrentVersion + 1; m <= ceiling; m++ {
		exists, err := probe(owner, TargetRepoName(parsed.BaseName, m))
		if err != nil {
			return highest, fmt.Errorf("probe %s/%s-v%d: %w", owner, parsed.BaseName, m, err)
		}
		if !exists {
			return highest, nil
		}
		highest = m
	}

	return highest, nil
}
