package clonenow

// Tests for the idempotent re-clone behavior. Cover the inspect +
// dispatch logic without spawning real `git clone` (network-free):
//
//   - inspectExistingRepo on missing / empty / non-repo dirs
//   - urlsMatch + canonicalRepoID for HTTPS<->SSH equivalence
//   - dispatchOnExists "already matches" skip path
//   - dispatchOnExists URL/branch mismatch under "skip" policy
//   - dispatchOnExists "not a repo" hard-failure under every policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestInspectExistingRepo_Missing(t *testing.T) {
	state := inspectExistingRepo(filepath.Join(t.TempDir(), "nope"))
	if state.Exists || state.IsRepo {
		t.Errorf("missing dir reported as present: %+v", state)
	}
}

func TestInspectExistingRepo_EmptyDir(t *testing.T) {
	state := inspectExistingRepo(t.TempDir())
	if !state.Exists || !state.Empty || state.IsRepo {
		t.Errorf("empty dir misclassified: %+v", state)
	}
}

func TestInspectExistingRepo_NonRepoPopulated(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	state := inspectExistingRepo(dir)
	if !state.Exists || state.Empty || state.IsRepo {
		t.Errorf("populated non-repo misclassified: %+v", state)
	}
}

func TestUrlsMatch_HTTPSandSSHEquivalence(t *testing.T) {
	cases := [][2]string{
		{"https://github.com/owner/repo.git", "git@github.com:owner/repo.git"},
		{"https://github.com/owner/repo", "https://github.com/owner/repo.git"},
		{"ssh://git@github.com/owner/repo.git", "https://github.com/Owner/Repo"},
	}
	for _, c := range cases {
		if !urlsMatch(c[0], c[1]) {
			t.Errorf("expected match: %q vs %q (canon=%q,%q)",
				c[0], c[1], canonicalRepoID(c[0]), canonicalRepoID(c[1]))
		}
	}
	if urlsMatch("https://github.com/a/b", "https://github.com/a/c") {
		t.Errorf("different repos should not match")
	}
}

func TestDispatchOnExists_AlreadyMatchesIsSkip(t *testing.T) {
	state := existingRepoState{
		Exists: true, IsRepo: true,
		RemoteURL: "https://github.com/owner/repo.git",
		Branch:    "main",
	}
	row := Row{HTTPSUrl: "https://github.com/owner/repo.git", Branch: "main"}
	res := dispatchOnExists(row, row.HTTPSUrl, "/abs", "/cwd",
		constants.CloneNowOnExistsSkip, state)
	if res.Status != constants.CloneNowStatusSkipped {
		t.Errorf("status = %q, want skipped", res.Status)
	}
	if res.Detail != constants.MsgCloneNowAlreadyMatches {
		t.Errorf("detail = %q", res.Detail)
	}
}

func TestDispatchOnExists_URLMismatchSkipsWithReason(t *testing.T) {
	state := existingRepoState{
		Exists: true, IsRepo: true,
		RemoteURL: "https://github.com/old/repo.git",
		Branch:    "main",
	}
	row := Row{HTTPSUrl: "https://github.com/new/repo.git", Branch: "main"}
	res := dispatchOnExists(row, row.HTTPSUrl, "/abs", "/cwd",
		constants.CloneNowOnExistsSkip, state)
	if res.Status != constants.CloneNowStatusSkipped {
		t.Errorf("status = %q, want skipped", res.Status)
	}
	if !strings.Contains(res.Detail, "remote url differs") {
		t.Errorf("detail must explain url drift: %q", res.Detail)
	}
}

func TestDispatchOnExists_BranchMismatchSkipsWithReason(t *testing.T) {
	state := existingRepoState{
		Exists: true, IsRepo: true,
		RemoteURL: "https://github.com/owner/repo.git",
		Branch:    "develop",
	}
	row := Row{HTTPSUrl: "https://github.com/owner/repo.git", Branch: "main"}
	res := dispatchOnExists(row, row.HTTPSUrl, "/abs", "/cwd",
		constants.CloneNowOnExistsSkip, state)
	if res.Status != constants.CloneNowStatusSkipped {
		t.Errorf("status = %q, want skipped", res.Status)
	}
	if !strings.Contains(res.Detail, "branch differs") {
		t.Errorf("detail must explain branch drift: %q", res.Detail)
	}
}

func TestDispatchOnExists_NonRepoFailsUnderEveryPolicy(t *testing.T) {
	state := existingRepoState{Exists: true, IsRepo: false, Empty: false}
	row := Row{HTTPSUrl: "https://x/a.git"}
	for _, policy := range []string{
		constants.CloneNowOnExistsSkip,
		constants.CloneNowOnExistsUpdate,
		constants.CloneNowOnExistsForce,
	} {
		res := dispatchOnExists(row, row.HTTPSUrl, "/abs", "/cwd", policy, state)
		if res.Status != constants.CloneNowStatusFailed {
			t.Errorf("policy=%s: status = %q, want failed", policy, res.Status)
		}
		if res.Detail != constants.MsgCloneNowNotARepo {
			t.Errorf("policy=%s: detail = %q", policy, res.Detail)
		}
	}
}

func TestRepoMatches_EmptyBranchMeansAny(t *testing.T) {
	state := existingRepoState{
		IsRepo:    true,
		RemoteURL: "https://github.com/owner/repo.git",
		Branch:    "anything",
	}
	row := Row{HTTPSUrl: "https://github.com/owner/repo.git", Branch: ""}
	if !repoMatches(row, row.HTTPSUrl, state) {
		t.Errorf("empty planned branch should match any local branch")
	}
}
