package cloner

import "testing"

func TestIsLFSSmudgeFailure(t *testing.T) {
	positives := []string{
		"fatal: public/favicon.ico: smudge filter lfs failed\n",
		"error: external filter 'git-lfs filter-process' failed",
		"Errors logged to '.git/lfs/logs/x.log'. external filter 'git-lfs filter-process' failed",
	}
	for _, s := range positives {
		if !isLFSSmudgeFailure(s) {
			t.Errorf("expected LFS smudge detection for: %q", s)
		}
	}
	negatives := []string{
		"", "fatal: repository not found", "Permission denied (publickey)",
	}
	for _, s := range negatives {
		if isLFSSmudgeFailure(s) {
			t.Errorf("unexpected LFS smudge detection for: %q", s)
		}
	}
}
