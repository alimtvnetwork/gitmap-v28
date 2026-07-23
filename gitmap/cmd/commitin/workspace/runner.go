package workspace

import (
	"fmt"
	"os/exec"
	"strings"
)

// gitRunner is the default git executor; tests swap it for a fake.
// Signature: (subcommand, args...) -> error. Implementation runs in
// the caller's CWD; subcommands needing a target dir pass it as an
// arg (e.g. clone) — that's why we don't expose Dir here.
var gitRunner = defaultGitRunner

// defaultGitRunner shells out to the system `git` binary, capturing
// combined output for inclusion in any returned error.
func defaultGitRunner(sub string, args ...string) error {
	full := append([]string{sub}, args...)
	cmd := exec.Command("git", full...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s: %w (output: %s)", sub, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// SetGitRunnerForTest swaps the package-level git runner with a
// caller-provided fake. Returns a restore func; tests `defer` it.
// Exposed for the workspace_test.go suite — not for production use.
func SetGitRunnerForTest(fake func(sub string, args ...string) error) func() {
	prev := gitRunner
	gitRunner = fake
	return func() { gitRunner = prev }
}
