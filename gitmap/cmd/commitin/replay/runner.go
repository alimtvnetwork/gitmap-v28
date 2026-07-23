package replay

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// gitRunner is the swappable text-stdout git executor. Tests inject
// fakes; production runs the real binary.
var gitRunner = defaultGitRunner

// gitRunnerBytes is the swappable binary-stdout variant (used to read
// raw blob contents — `git cat-file blob` is binary-safe).
var gitRunnerBytes = defaultGitRunnerBytes

// gitRunnerEnv runs git with extra env vars on TOP of the parent env;
// used by ApplyCommit to set GIT_AUTHOR_*/GIT_COMMITTER_*.
var gitRunnerEnv = defaultGitRunnerEnv

// hashObjectStdin pipes data through `git hash-object -w --stdin` and
// returns the resulting blob SHA. Swappable for tests.
var hashObjectStdin = defaultHashObjectStdin

// defaultGitRunner runs `git -C <repoDir> <sub> <args...>` and returns
// trimmed text stdout.
func defaultGitRunner(repoDir, sub string, args ...string) (string, error) {
	full := append([]string{"-C", repoDir, sub}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	trimmed := strings.TrimRight(string(out), "\n")
	if err != nil {
		return trimmed, fmt.Errorf("git %s: %w (%s)", sub, err, strings.TrimSpace(trimmed))
	}
	return trimmed, nil
}

// defaultGitRunnerBytes is the binary-safe twin: returns raw stdout,
// stderr is captured into the error.
func defaultGitRunnerBytes(repoDir, sub string, args ...string) ([]byte, error) {
	full := append([]string{"-C", repoDir, sub}, args...)
	cmd := exec.Command("git", full...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return out, fmt.Errorf("git %s: %w (%s)", sub, err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

// defaultGitRunnerEnv runs git with extra env vars appended to os.Environ().
func defaultGitRunnerEnv(repoDir string, extraEnv []string, args ...string) (string, error) {
	full := append([]string{"-C", repoDir}, args...)
	cmd := exec.Command("git", full...)
	cmd.Env = append(os.Environ(), extraEnv...)
	out, err := cmd.CombinedOutput()
	trimmed := strings.TrimRight(string(out), "\n")
	if err != nil {
		return trimmed, fmt.Errorf("git %s: %w (%s)", args[0], err, strings.TrimSpace(trimmed))
	}
	return trimmed, nil
}

// defaultHashObjectStdin pipes raw bytes into `git hash-object -w --stdin`.
func defaultHashObjectStdin(repoDir string, data []byte) (string, error) {
	cmd := exec.Command("git", "-C", repoDir, "hash-object", "-w", "--stdin")
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(out))
	if err != nil {
		return trimmed, fmt.Errorf("git hash-object: %w (%s)", err, trimmed)
	}
	return trimmed, nil
}

// SetTestHooks swaps every package-level git hook with caller-provided
// fakes. Returns a single restore func; tests `defer` it.
func SetTestHooks(
	textFn func(repoDir, sub string, args ...string) (string, error),
	bytesFn func(repoDir, sub string, args ...string) ([]byte, error),
	envFn func(repoDir string, extraEnv []string, args ...string) (string, error),
	hashFn func(repoDir string, data []byte) (string, error),
) func() {
	prevText, prevBytes, prevEnv, prevHash := gitRunner, gitRunnerBytes, gitRunnerEnv, hashObjectStdin
	if textFn != nil {
		gitRunner = textFn
	}
	if bytesFn != nil {
		gitRunnerBytes = bytesFn
	}
	if envFn != nil {
		gitRunnerEnv = envFn
	}
	if hashFn != nil {
		hashObjectStdin = hashFn
	}
	return func() {
		gitRunner, gitRunnerBytes, gitRunnerEnv, hashObjectStdin = prevText, prevBytes, prevEnv, prevHash
	}
}
