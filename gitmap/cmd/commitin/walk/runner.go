package walk

import (
	"fmt"
	"os/exec"
	"strings"
)

// gitRunner is the swappable git executor for the walk package. Tests
// use SetGitRunnerForTest to inject a fake; production runs the real
// `git` binary via defaultGitRunner.
//
// Signature: (repoDir, subcommand, args...) -> (stdout, error). On
// non-zero exit, the error message includes git's combined output so
// callers can surface diagnostics verbatim.
var gitRunner = defaultGitRunner

// defaultGitRunner shells out to `git -C <repoDir> <sub> <args...>`.
func defaultGitRunner(repoDir, sub string, args ...string) (string, error) {
	full := append([]string{"-C", repoDir, sub}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	trimmed := strings.TrimRight(string(out), "\n")
	if err != nil {
		return trimmed, fmt.Errorf("git %s: %w (%s)", sub, err, strings.TrimSpace(trimmed))
	}
	return trimmed, nil
}

// SetGitRunnerForTest replaces the package-level git runner; returns
// a restore func meant for `defer`. Exported only for the test suite.
func SetGitRunnerForTest(fake func(repoDir, sub string, args ...string) (string, error)) func() {
	prev := gitRunner
	gitRunner = fake
	return func() { gitRunner = prev }
}
