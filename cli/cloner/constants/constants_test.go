package constants

import "testing"

func TestConstants(t *testing.T) {
	if len(GitPullFastForwardArgs) == 0 {
		t.Fatal("empty fast forward args")
	}
	if len(GitPullAutoMergeArgs) == 0 {
		t.Fatal("empty auto merge args")
	}
	if EnvGitLFSNoSmudge == "" {
		t.Fatal("empty LFS smudge env")
	}
}
