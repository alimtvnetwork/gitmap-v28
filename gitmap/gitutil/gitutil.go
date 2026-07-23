// Package gitutil extracts Git metadata by running git commands.
package gitutil

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// RepoStatus holds the live state of a Git repository.
type RepoStatus struct {
	Branch      string
	Dirty       bool
	Untracked   int
	Modified    int
	Staged      int
	Ahead       int
	Behind      int
	StashCount  int
	Unreachable bool
}

// RemoteURL returns the origin remote URL for a repo at the given path.
func RemoteURL(repoPath string) (string, error) {
	out, err := runGit(repoPath,
		constants.GitConfigCmd, constants.GitGetFlag, constants.GitRemoteOrigin)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

// RepoRoot returns the absolute path of the working-tree root for the
// repo containing `startPath`. Returns the underlying git error when
// startPath is not inside a git working tree.
func RepoRoot(startPath string) (string, error) {
	out, err := runGit(startPath, constants.GitRevParse, "--show-toplevel")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

// Branch source labels describe how a repo's branch was determined.
const (
	BranchSourceHEAD           = "HEAD"
	BranchSourceDetached       = "detached"
	BranchSourceRemoteTracking = "remote-tracking"
	BranchSourceDefault        = "default"
	BranchSourceUnknown        = "unknown"
)

// CurrentBranch returns the current branch name for a repo.
func CurrentBranch(repoPath string) (string, error) {
	out, err := runGit(repoPath,
		constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if err != nil {
		return constants.DefaultBranch, err
	}

	return strings.TrimSpace(out), nil
}

// DetectBranch returns the branch name and a label describing how it was
// detected, using constants.DefaultBranch as the last-resort fallback.
// Equivalent to DetectBranchWithDefault(repoPath, constants.DefaultBranch).
// Kept as the legacy entry point so existing callers compile unchanged.
func DetectBranch(repoPath string) (branch, source string) {
	return DetectBranchWithDefault(repoPath, constants.DefaultBranch)
}

// DetectBranchWithDefault is DetectBranch with a caller-supplied
// fallback branch name. Resolution order (each step falls through to the
// next on miss):
//
//  1. HEAD via `git rev-parse --abbrev-ref HEAD` → "HEAD" when on a named
//     branch. NOTE: a detached HEAD or empty result no longer terminates
//     resolution — we fall through to the remote-default lookups so the
//     returned name is always usable for `git clone -b` when possible.
//  2. Local remote-tracking ref via `git symbolic-ref refs/remotes/origin/
//     HEAD` → "remote-tracking". Works when the repo was cloned with
//     `--single-branch` skipped (the common case).
//  3. Live remote query via `git ls-remote --symref origin HEAD` →
//     "remote-tracking". Covers single-branch clones and shallow mirrors
//     where step 2's local ref is absent.
//  4. Caller-supplied `fallback` → "default". Last-resort guess so
//     callers always have SOMETHING to attempt. When fallback is the
//     empty string, this step is skipped — the function continues to
//     the detached-HEAD / unknown sentinels below so callers that want
//     to *see* "we have nothing" can pass "" instead of a placeholder.
//
// Only when every step fails does the function return ("HEAD", "detached")
// for an actual detached HEAD or ("", "unknown") for a totally opaque repo —
// preserving the original sentinels for downstream consumers (cloner/strategy.go)
// that branch on BranchSource.
func DetectBranchWithDefault(repoPath, fallback string) (branch, source string) {
	if name, ok := detectFromLocalHEAD(repoPath); ok {
		return name, BranchSourceHEAD
	}
	if name, ok := detectFromLocalRemoteRef(repoPath); ok {
		return name, BranchSourceRemoteTracking
	}
	if name, ok := detectFromLiveRemote(repoPath); ok {
		return name, BranchSourceRemoteTracking
	}
	if len(fallback) > 0 {
		return fallback, BranchSourceDefault
	}
	if isDetachedHEAD(repoPath) {

		return constants.GitHEAD, BranchSourceDetached
	}

	return "", BranchSourceUnknown
}

// detectFromLocalHEAD reads the named branch from `git rev-parse
// --abbrev-ref HEAD`. Returns ok=false for detached HEAD (literal "HEAD")
// or empty output so the caller can fall through to remote-default lookups.
func detectFromLocalHEAD(repoPath string) (string, bool) {
	out, err := runGit(repoPath,
		constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if err != nil {

		return "", false
	}
	name := strings.TrimSpace(out)
	if len(name) == 0 || name == constants.GitHEAD {

		return "", false
	}

	return name, true
}

// detectFromLocalRemoteRef reads the locally-tracked default via
// `git symbolic-ref refs/remotes/origin/HEAD`. Returns ok=false when the
// ref is missing (single-branch / mirror clones) so the caller can fall
// through to a live `ls-remote` query.
func detectFromLocalRemoteRef(repoPath string) (string, bool) {
	out, err := runGit(repoPath,
		"symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {

		return "", false
	}
	const prefix = "refs/remotes/origin/"
	ref := strings.TrimSpace(out)
	if !strings.HasPrefix(ref, prefix) {

		return "", false
	}

	return strings.TrimPrefix(ref, prefix), true
}

// detectFromLiveRemote queries the origin remote directly with
// `git ls-remote --symref origin HEAD` and parses the first symref line:
//
//	ref: refs/heads/main\tHEAD
//
// Network-dependent — failure is silent so offline runs simply fall through
// to the built-in default. Bounded by git's own timeout / credential helper.
func detectFromLiveRemote(repoPath string) (string, bool) {
	out, err := runGit(repoPath,
		"ls-remote", "--symref", "origin", constants.GitHEAD)
	if err != nil {

		return "", false
	}

	return parseLsRemoteSymref(out)
}

// parseLsRemoteSymref extracts the branch name from the first `ref: ...`
// line of `git ls-remote --symref` output. Split out for unit-testability —
// no git invocation, pure string parsing.
func parseLsRemoteSymref(output string) (string, bool) {
	const refPrefix = "refs/heads/"
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "ref: ") {
			continue
		}
		// Format: "ref: refs/heads/<name>\tHEAD"
		body := strings.TrimPrefix(trimmed, "ref: ")
		fields := strings.Fields(body)
		if len(fields) == 0 {
			continue
		}
		if !strings.HasPrefix(fields[0], refPrefix) {
			continue
		}

		return strings.TrimPrefix(fields[0], refPrefix), true
	}

	return "", false
}

// isDetachedHEAD reports whether `git rev-parse --abbrev-ref HEAD` returns
// the literal "HEAD" sentinel. Used as a final disambiguator after every
// fallback has been exhausted, so the returned source label accurately
// distinguishes a true detached state from a fully opaque repo.
func isDetachedHEAD(repoPath string) bool {
	out, err := runGit(repoPath,
		constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if err != nil {

		return false
	}

	return strings.TrimSpace(out) == constants.GitHEAD
}

// Status returns the full live status of a repository.
// If the path does not exist or is not a git repo, Unreachable is set.
func Status(repoPath string) RepoStatus {
	rs := RepoStatus{}

	if _, err := os.Stat(repoPath); err != nil {
		rs.Unreachable = true
		return rs
	}

	branch, err := CurrentBranch(repoPath)
	if err != nil {
		rs.Unreachable = true
		return rs
	}

	rs.Branch = branch
	rs.Dirty, rs.Untracked, rs.Modified, rs.Staged = parsePortcelainStatus(repoPath)
	rs.Ahead, rs.Behind = parseAheadBehind(repoPath)
	rs.StashCount = countStashes(repoPath)

	return rs
}

// parsePortcelainStatus runs git status --porcelain and counts file states.
func parsePortcelainStatus(repoPath string) (dirty bool, untracked, modified, staged int) {
	out, err := runGit(repoPath, "status", "--porcelain")
	if err != nil {
		return false, 0, 0, 0
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		x := line[0]
		y := line[1]
		if x == '?' && y == '?' {
			untracked++
		} else if x != ' ' && x != '?' {
			staged++
		}
		if y != ' ' && y != '?' {
			modified++
		}
	}
	dirty = (untracked + modified + staged) > 0

	return dirty, untracked, modified, staged
}

// parseAheadBehind extracts ahead/behind counts from rev-list.
func parseAheadBehind(repoPath string) (ahead, behind int) {
	out, err := runGit(repoPath, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if err != nil {
		return 0, 0
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) == 2 {
		ahead, _ = strconv.Atoi(parts[0])
		behind, _ = strconv.Atoi(parts[1])
	}

	return ahead, behind
}

// countStashes returns the number of stash entries.
func countStashes(repoPath string) int {
	out, err := runGit(repoPath, "stash", "list")
	if err != nil || len(strings.TrimSpace(out)) == 0 {
		return 0
	}

	return len(strings.Split(strings.TrimSpace(out), "\n"))
}

// runGit executes a git command in the given directory and returns stdout.
func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
