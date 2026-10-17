// Package cmdrelease — release_gh.go: shared `gh` CLI helpers for the
// program-275 release audit commands (verify, winscan, train). Read-only.
package cmdrelease

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// requireGhCLI returns an AppError naming the tool when `gh` is missing.
func requireGhCLI() *appfault.AppError {
	if _, err := exec.LookPath("gh"); err != nil {
		return appfault.NewNotFoundError("gh CLI not found — install from https://cli.github.com " +
			"or set it on PATH; release audit commands shell out to it")
	}

	return nil
}

// runGhCLI runs `gh` with args and returns trimmed stdout.
func runGhCLI(args ...string) (string, *appfault.AppError) {
	cmd := exec.Command("gh", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", appfault.NewValidationError(
			"gh " + strings.Join(args, " ") + " failed: " + ghErrDetail(err))
	}

	return strings.TrimSpace(string(out)), nil
}

// ghErrDetail extracts the stderr text from a failed gh invocation.
func ghErrDetail(err error) string {
	if exitErr, isExit := err.(*exec.ExitError); isExit {
		if stderr := strings.TrimSpace(string(exitErr.Stderr)); stderr != "" {
			return stderr
		}
	}

	return err.Error()
}

// normalizeReleaseRepo accepts `owner/repo` or a bare repo name (resolved to
// the alimtvnetwork org, the owner's GitHub identity).
func normalizeReleaseRepo(repo string) string {
	repo = strings.TrimSpace(repo)
	if repo == "" || strings.Contains(repo, "/") {
		return repo
	}

	return "alimtvnetwork/" + repo
}
