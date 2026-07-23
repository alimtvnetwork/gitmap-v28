package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// verifyHistoryRewrite runs the spec §3 / §4 verification loops. Exit
// 6 on any disagreement.
func verifyHistoryRewrite(mode historyMode, sandbox string, paths []string) {
	fmt.Fprint(os.Stderr, constants.HistoryMsgPhaseVerify)
	for _, p := range paths {
		if mode == historyModePurge {
			verifyPurgePath(sandbox, p)
		} else {
			verifyPinPath(sandbox, p)
		}
	}
	fmt.Fprint(os.Stderr, constants.HistoryMsgVerifyOk)
}

// verifyPurgePath asserts no commit on any branch still references P.
func verifyPurgePath(sandbox, path string) {
	cmd := exec.Command(constants.HistoryGitBin, "-C", sandbox, "log", "--all",
		"--oneline", "--", path)
	out, _ := cmd.Output()
	count := nonEmptyLineCount(string(out))
	if count > 0 {
		fmt.Fprintf(os.Stderr, constants.HistoryErrVerifyPurge, path, count)
		os.Exit(constants.HistoryExitVerifyFailed)
	}
}

// verifyPinPath asserts every commit's content for P hashes to the
// same SHA-256 (i.e. the file is byte-identical across all history).
func verifyPinPath(sandbox, path string) {
	shas := commitsTouchingPath(sandbox, path)
	uniq := make(map[string]bool)
	for _, commit := range shas {
		hash, ok := hashPathAtCommit(sandbox, commit, path)
		if !ok {
			continue
		}
		uniq[hash] = true
	}
	if len(uniq) > 1 {
		fmt.Fprintf(os.Stderr, constants.HistoryErrVerifyPin, path, len(uniq))
		os.Exit(constants.HistoryExitVerifyFailed)
	}
}

// commitsTouchingPath lists every commit SHA on any branch whose tree
// contains the given path.
func commitsTouchingPath(sandbox, path string) []string {
	cmd := exec.Command(constants.HistoryGitBin, "-C", sandbox, "log", "--all",
		"--pretty=format:%H", "--", path)
	out, _ := cmd.Output()
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	clean := make([]string, 0, len(lines))
	for _, l := range lines {
		if l != "" {
			clean = append(clean, l)
		}
	}
	return clean
}

// hashPathAtCommit returns the hex sha256 of `git show <commit>:<path>`.
// The bool is false when the path didn't exist at that commit.
func hashPathAtCommit(sandbox, commit, path string) (string, bool) {
	cmd := exec.Command(constants.HistoryGitBin, "-C", sandbox, "show",
		commit+":"+path)
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	sum := sha256.Sum256(out)
	return hex.EncodeToString(sum[:]), true
}

// nonEmptyLineCount counts lines containing any non-whitespace.
func nonEmptyLineCount(s string) int {
	n := 0
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	return n
}
