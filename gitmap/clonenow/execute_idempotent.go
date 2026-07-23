package clonenow

// Idempotency layer: the inspect-then-dispatch logic that lets
// clone-now make a smart decision about an already-populated
// destination instead of the old "non-empty dir = skip" heuristic.
//
// Why a separate file? execute.go is on the per-file budget already
// (Execute + executeRow + the runGitClone helpers + the trim-error
// utilities). Pulling the inspect / dispatch / update / force code
// here keeps each file focused on a single concern -- the executor
// orchestrates, this file decides what to do with what's on disk.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// existingRepoState captures the result of probing an on-disk
// destination. It is intentionally a small value type with no
// methods -- the consumer (dispatchOnExists) reads the fields
// directly so each branch's decision is auditable in one place
// instead of being hidden behind state.IsX() helpers.
type existingRepoState struct {
	// Exists is true when the destination directory exists on disk.
	// A false value is the trivial "no conflict, just clone" case.
	Exists bool
	// IsRepo is true when the destination is a git work tree (we
	// detect this by probing for a .git entry -- file or dir, since
	// worktrees use a file). Reused by every on-exists branch.
	IsRepo bool
	// Empty is true when Exists && !IsRepo and the directory has
	// no children. An empty dir is safe to remove + clone into;
	// a populated non-repo dir is treated as a hard failure under
	// every policy (we never destroy unrelated user data).
	Empty bool
	// RemoteURL is the origin remote URL as reported by `git config
	// remote.origin.url`. Empty when IsRepo is false or when the
	// repo has no origin remote configured.
	RemoteURL string
	// Branch is the currently checked-out branch name. Empty when
	// IsRepo is false or when HEAD is detached.
	Branch string
}

// inspectExistingRepo probes the destination and returns a populated
// existingRepoState. Always returns a value (never errors): a probe
// failure simply means fewer fields are populated, and downstream
// branches treat "missing field" as "can't confirm match", which
// safely falls through to the more conservative skip path.
func inspectExistingRepo(absDest string) existingRepoState {
	state := existingRepoState{}
	info, err := os.Stat(absDest)
	if err != nil || !info.IsDir() {
		return state
	}
	state.Exists = true
	entries, _ := os.ReadDir(absDest)
	state.Empty = len(entries) == 0
	if !isGitWorkTree(absDest) {
		return state
	}
	state.IsRepo = true
	if remote, err := gitutil.RemoteURL(absDest); err == nil {
		state.RemoteURL = strings.TrimSpace(remote)
	}
	if branch, err := gitutil.CurrentBranch(absDest); err == nil {
		state.Branch = strings.TrimSpace(branch)
	}

	return state
}

// isGitWorkTree returns true when absDest contains a `.git` entry
// (directory for normal clones, file for worktree-style links).
// We probe directly rather than shelling out to `git rev-parse`
// because the syscall is ~100x cheaper and runs on every row.
func isGitWorkTree(absDest string) bool {
	_, err := os.Stat(filepath.Join(absDest, ".git"))

	return err == nil
}

// dispatchOnExists is the policy switch. Returns a Result with only
// Status + Detail populated; the caller (executeRow) overlays the
// row/url/dest/duration fields so this function doesn't need to
// know about timing or row identity.
//
// The four possible outcomes per branch are:
//
//   - Nothing on disk -> clone fresh (status = ok or failed).
//   - Match (skip/update) -> "already matches" (status = skipped).
//   - Mismatch under skip -> explanatory skip detail.
//   - Mismatch under update -> fetch + checkout (status = ok/failed).
//   - Force -> remove + reclone (status = ok/failed).
//   - Non-repo populated dir -> hard failure under every policy.
func dispatchOnExists(r Row, url, absDest, cwd, policy string, state existingRepoState) Result {
	if !state.Exists || state.Empty {
		return cloneFresh(r, url, absDest, cwd)
	}
	if !state.IsRepo {
		return Result{
			Status: constants.CloneNowStatusFailed,
			Detail: constants.MsgCloneNowNotARepo,
		}
	}
	if policy == constants.CloneNowOnExistsForce {
		return forceReclone(r, url, absDest, cwd)
	}
	if repoMatches(r, url, state) {
		return Result{
			Status: constants.CloneNowStatusSkipped,
			Detail: constants.MsgCloneNowAlreadyMatches,
		}
	}
	if policy == constants.CloneNowOnExistsUpdate {
		return updateExisting(r, url, absDest, state)
	}

	return mismatchSkipResult(r, url, state)
}

// cloneFresh removes an empty leftover directory (so `git clone`
// doesn't refuse with "destination path already exists"), then
// runs the clone. Empty-dir removal is bounded to the exact dest
// path; we do NOT recursively remove anything we didn't just confirm
// is empty.
func cloneFresh(r Row, url, absDest, cwd string) Result {
	if info, err := os.Stat(absDest); err == nil && info.IsDir() {
		_ = os.Remove(absDest) // empty -> ok; non-empty -> git will error.
	}
	dest := relOrAbs(absDest, cwd)
	detail, ok := runGitClone(r, url, dest, cwd)
	if !ok {
		return Result{Status: constants.CloneNowStatusFailed, Detail: detail}
	}

	return Result{Status: constants.CloneNowStatusOK, Detail: detail}
}

// repoMatches returns true when the on-disk repo's remote URL +
// current branch match the planned row. URL comparison is shape-
// agnostic (we collapse https/ssh forms to a canonical owner/repo
// string) so a row whose plan-recorded URL is HTTPS but whose local
// clone uses SSH is still considered a match.
//
// Branch comparison treats an empty planned branch as "any branch
// is fine" -- the user didn't pin one, so whatever is checked out
// locally satisfies the plan.
func repoMatches(r Row, url string, state existingRepoState) bool {
	if !urlsMatch(state.RemoteURL, url) && !urlsMatch(state.RemoteURL, r.HTTPSUrl) &&
		!urlsMatch(state.RemoteURL, r.SSHUrl) {
		return false
	}
	if len(r.Branch) == 0 {
		return true
	}

	return state.Branch == r.Branch
}

// urlsMatch returns true when two git URLs point at the same repo
// regardless of HTTPS vs SSH form. The canonicalization is "lower-
// case host + lower-case path with .git stripped" -- coarse enough
// to absorb the two URL shapes git accepts, strict enough to NOT
// false-match repos on the same host with similar names.
func urlsMatch(a, b string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	return canonicalRepoID(a) == canonicalRepoID(b)
}

// canonicalRepoID collapses an https / ssh / scp-style git URL down
// to a "host/path" string suitable for equality comparison. Returns
// the original input on parse failure so unfamiliar URL shapes
// degrade to byte-equality rather than silently false-matching.
func canonicalRepoID(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSuffix(s, ".git")
	if strings.HasPrefix(s, "https://") {
		s = strings.TrimPrefix(s, "https://")
	} else if strings.HasPrefix(s, "http://") {
		s = strings.TrimPrefix(s, "http://")
	} else if strings.HasPrefix(s, "ssh://") {
		s = strings.TrimPrefix(s, "ssh://")
		s = strings.TrimPrefix(s, "git@")
	} else if at := strings.Index(s, "@"); at >= 0 {
		// scp-style: git@host:owner/repo -> host:owner/repo
		s = s[at+1:]
		s = strings.Replace(s, ":", "/", 1)
	}

	return strings.ToLower(s)
}

// mismatchSkipResult builds a "skip with explanation" Result for the
// default skip policy. The detail names which dimension differs
// (URL vs branch) so the user sees actionable drift information,
// not a generic "skipped".
func mismatchSkipResult(r Row, url string, state existingRepoState) Result {
	if !urlsMatch(state.RemoteURL, url) {
		return Result{
			Status: constants.CloneNowStatusSkipped,
			Detail: fmt.Sprintf(constants.MsgCloneNowURLMismatch,
				state.RemoteURL, url),
		}
	}

	return Result{
		Status: constants.CloneNowStatusSkipped,
		Detail: fmt.Sprintf(constants.MsgCloneNowBranchMismatch,
			state.Branch, r.Branch),
	}
}

// updateExisting brings a mismatched on-disk repo into alignment
// with the plan: fetch from origin, then checkout the planned
// branch. Skips the checkout when the plan didn't pin a branch
// (a fetch alone is enough to surface upstream changes without
// switching what the user has open).
//
// We deliberately do NOT run `git pull` -- pull would merge or
// rebase, both of which can fail or produce a dirty work tree on
// repos with local changes. Fetch + checkout is the strongest
// non-destructive option: it advances refs and switches branches
// without touching tracked files outside the checkout transition.
func updateExisting(r Row, url, absDest string, state existingRepoState) Result {
	if detail, ok := runGitFetch(absDest); !ok {
		return Result{
			Status: constants.CloneNowStatusFailed,
			Detail: fmt.Sprintf(constants.MsgCloneNowFetchFail, detail),
		}
	}
	if len(r.Branch) > 0 && state.Branch != r.Branch {
		if detail, ok := runGitCheckout(absDest, r.Branch); !ok {
			return Result{
				Status: constants.CloneNowStatusFailed,
				Detail: fmt.Sprintf(constants.MsgCloneNowCheckoutFail,
					r.Branch, detail),
			}
		}
	}
	_ = url // url unused once update succeeds -- the existing remote stays.

	return Result{Status: constants.CloneNowStatusOK, Detail: constants.MsgCloneNowUpdated}
}

// forceReclone removes the existing destination and re-clones from
// scratch. Guarded by the caller (dispatchOnExists) so we only
// reach here when state.IsRepo is true -- we never blow away an
// unrelated user directory.
func forceReclone(r Row, url, absDest, cwd string) Result {
	if err := os.RemoveAll(absDest); err != nil {
		return Result{
			Status: constants.CloneNowStatusFailed,
			Detail: fmt.Sprintf(constants.MsgCloneNowForceRemoveFail, absDest, err),
		}
	}
	dest := relOrAbs(absDest, cwd)
	detail, ok := runGitClone(r, url, dest, cwd)
	if !ok {
		return Result{Status: constants.CloneNowStatusFailed, Detail: detail}
	}
	_ = detail

	return Result{Status: constants.CloneNowStatusOK, Detail: constants.MsgCloneNowForceRecloned}
}

// runGitFetch wraps `git fetch --all --prune` in the standard
// (detail, ok) shape the rest of the executor uses. Stderr is
// trimmed via the same helper as runGitClone so error formatting
// stays consistent across the per-row Result.Detail field.
func runGitFetch(absDest string) (string, bool) {
	cmd := exec.Command(constants.GitBin, "fetch", "--all", "--prune")
	cmd.Dir = absDest
	out, err := cmd.CombinedOutput()
	if err != nil {
		return trimGitError(string(out), err), false
	}

	return "", true
}

// runGitCheckout wraps `git checkout <branch>` in the same
// (detail, ok) shape. Plain checkout (not `-B`) so a non-existent
// branch surfaces as a clear failure rather than silently creating
// a new branch that wouldn't track the remote.
func runGitCheckout(absDest, branch string) (string, bool) {
	cmd := exec.Command(constants.GitBin, "checkout", branch)
	cmd.Dir = absDest
	out, err := cmd.CombinedOutput()
	if err != nil {
		return trimGitError(string(out), err), false
	}

	return "", true
}

// relOrAbs returns the destination path the way the cloner should
// pass it to `git clone`: relative when it's inside cwd (so progress
// output is short), absolute otherwise. Centralized so cloneFresh
// and forceReclone build identical command lines for the same dest.
func relOrAbs(absDest, cwd string) string {
	rel, err := filepath.Rel(cwd, absDest)
	if err != nil || strings.HasPrefix(rel, "..") {
		return absDest
	}

	return rel
}
