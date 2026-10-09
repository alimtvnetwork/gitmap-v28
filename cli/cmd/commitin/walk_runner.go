package commitin

import ()

// gitRunner is the swappable git executor for the walk package. Tests
// use SetGitRunnerForTest to inject a fake; production runs the real
// `git` binary via defaultGitRunner.
//
// Signature: (repoDir, subcommand, args...) -> (stdout, error). On
// non-zero exit, the error message includes git's combined output so
// callers can surface diagnostics verbatim.
// defaultGitRunner shells out to `git -C <repoDir> <sub> <args...>`.
// SetGitRunnerForTest replaces the package-level git runner; returns
// a restore func meant for `defer`. Exported only for the test suite.
func SetGitRunnerForTest(fake func(repoDir, sub string, args ...string) (string, error)) func() {
	prev := gitRunner
	gitRunner = fake

	return func() { gitRunner = prev }
}
