package committransfer

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// gitOut runs `git -C dir <args...>` and returns trimmed stdout. On
// non-zero exit it returns the combined output as part of the error so
// callers can surface git's diagnostic verbatim.
func gitOut(dir string, args ...string) (string, error) {
	full := append([]string{constants.GitDirFlag, dir}, args...)
	out, err := exec.Command(constants.GitBin, full...).CombinedOutput()
	trimmed := strings.TrimSpace(string(out))
	if err != nil {
		return trimmed, fmt.Errorf("git %s: %w (%s)", strings.Join(args, " "), err, trimmed)
	}

	return trimmed, nil
}

// currentRefName returns the symbolic-ref short name (e.g. "main") if
// HEAD is on a branch, or the SHA otherwise. Used to restore the source
// working dir after the replay loop.
func currentRefName(dir string) (string, error) {
	if name, err := gitOut(dir, "symbolic-ref", "--short", "HEAD"); err == nil {
		return name, nil
	}

	return gitOut(dir, "rev-parse", "HEAD")
}

// mergeBase returns `git merge-base A B`. An error from git (e.g. no
// shared history) becomes ("", nil) — the planner treats unrelated
// histories as "use the source's full reachable history."
func mergeBase(dir, a, b string) (string, error) {
	out, err := gitOut(dir, "merge-base", a, b)
	if err != nil {
		return "", nil
	}

	return out, nil
}

// revListReverse lists commits in `<base>..<head>` oldest-first. Pass
// includeMerges=false to add `--no-merges` (the spec §3 default).
func revListReverse(dir, base, head string, includeMerges bool) ([]string, error) {
	args := []string{"rev-list", "--reverse"}
	if !includeMerges {
		args = append(args, "--no-merges")
	}
	rangeSpec := head
	if base != "" {
		rangeSpec = base + ".." + head
	}
	args = append(args, rangeSpec)
	out, err := gitOut(dir, args...)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}

	return strings.Split(out, "\n"), nil
}

// readCommit returns subject, body, author, author-date, short-sha for sha.
// Uses an ASCII unit-separator delimiter so subjects can contain anything.
const commitFormat = "%s%x1f%b%x1f%an <%ae>%x1f%aI%x1f%h"

func readCommit(dir, sha string) (subject, body, author, shortSHA string, when time.Time, err error) {
	out, gErr := gitOut(dir, "show", "-s", "--format="+commitFormat, sha)
	if gErr != nil {
		err = gErr

		return
	}
	parts := strings.SplitN(out, "\x1f", 5)
	if len(parts) != 5 {
		err = fmt.Errorf("unexpected git show output for %s: %q", sha, out)

		return
	}
	subject = parts[0]
	body = strings.TrimSpace(parts[1])
	author = parts[2]
	when, _ = time.Parse(time.RFC3339, parts[3])
	shortSHA = parts[4]

	return
}

// checkoutDetached moves HEAD to sha in detached state. Used for the
// per-commit snapshot during replay.
func checkoutDetached(dir, sha string) error {
	_, err := gitOut(dir, "checkout", "--quiet", "--detach", sha)

	return err
}

// checkoutRef restores the working dir to a named ref or sha (used in the
// deferred cleanup at the end of the replay loop).
func checkoutRef(dir, ref string) error {
	_, err := gitOut(dir, "checkout", "--quiet", ref)

	return err
}

// commitWithEnv runs `git commit -m <msg>` with GIT_AUTHOR_*/COMMITTER_*
// honoring the source author. Returns the new SHA.
func commitWithEnv(dir, msg, author string, when time.Time) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir, "commit",
		"--allow-empty-message", "-m", msg)
	cmd.Env = appendCommitEnv(author, when)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git commit: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	return gitOut(dir, "rev-parse", "HEAD")
}

// appendCommitEnv builds the GIT_AUTHOR_*/COMMITTER_* env vars so the
// replayed commit preserves the source author. Inherits the parent
// process env so the user's name/email become the *committer*.
func appendCommitEnv(author string, when time.Time) []string {
	name, email := splitAuthor(author)
	stamp := when.Format(time.RFC3339)
	base := []string{
		"GIT_AUTHOR_NAME=" + name,
		"GIT_AUTHOR_EMAIL=" + email,
		"GIT_AUTHOR_DATE=" + stamp,
	}

	return append(currentEnv(), base...)
}

// splitAuthor parses "Name <email>" into the two halves. Falls back to
// using the whole string as the name when the email block is missing.
func splitAuthor(author string) (name, email string) {
	open := strings.LastIndex(author, "<")
	close := strings.LastIndex(author, ">")
	if open < 0 || close < 0 || close < open {
		return strings.TrimSpace(author), ""
	}

	return strings.TrimSpace(author[:open]), strings.TrimSpace(author[open+1 : close])
}

// addAll stages every change in the target working dir (`git add -A`).
func addAll(dir string) error {
	_, err := gitOut(dir, "add", "-A")

	return err
}

// hasStagedChanges reports whether `git diff --cached --quiet` exits
// non-zero (which means there ARE staged changes).
func hasStagedChanges(dir string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir,
		"diff", "--cached", "--quiet")
	err := cmd.Run()

	return err != nil
}

// pushHEAD pushes the current branch. Returns the trimmed git output.
func pushHEAD(dir string) (string, error) {
	return gitOut(dir, "push")
}

// recentLogSubjectsAndBodies returns recent commits' full subject + body
// concatenated, used by the idempotence check (§10). When n <= 0 the
// entire target history is returned — required to avoid false-negative
// "fresh commit" classification on long-history targets (spec 114 Gap A).
func recentLogSubjectsAndBodies(dir string, n int) (string, error) {
	if n <= 0 {
		return gitOut(dir, "log", "--format=%B%n---commit-sep---")
	}

	return gitOut(dir, "log", fmt.Sprintf("-n%d", n), "--format=%B%n---commit-sep---")
}
