package cmd

// clonetermdetect.go — shared `git ls-remote --symref <url> HEAD`
// branch detector used by every clone-related command's
// `--output terminal` adapter (clone, clone-now, clone-pick,
// clone-from). clone-next has its own currentBranch helper because
// it reads the LOCAL HEAD of the source repo, not the remote.
//
// Keeping the detector here (not in the render package) lets us
// shell out to `git` without forcing render/ to depend on os/exec
// — render/ stays a pure formatter.

import (
	"context"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// detectRemoteHEAD runs `git ls-remote --symref <url> HEAD` with a
// short timeout and returns the resolved branch name (e.g. "main").
// Returns "" on any error so RenderRepoTermBlock falls back to its
// "(unknown)" placeholder rather than failing the whole command.
//
// The timeout is intentionally tight: we're previewing one repo
// before cloning it, and the user already accepted that the clone
// itself will take longer. A hung ls-remote should not block the
// terminal preview.
func detectRemoteHEAD(url string) string {
	if len(strings.TrimSpace(url)) == 0 {
		return ""
	}
	ctx, cancel := context.WithTimeout(
		context.Background(), constants.CloneTermDetectTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, constants.GitBin,
		"ls-remote", "--symref", url, "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return parseSymrefHEAD(string(out))
}

// parseSymrefHEAD pulls the branch out of `git ls-remote --symref`
// output. The first line looks like:
//
//	ref: refs/heads/main	HEAD
//
// We split on whitespace, take the second token, and strip the
// "refs/heads/" prefix. Anything we can't parse returns "" so the
// caller renders "(unknown)".
func parseSymrefHEAD(raw string) string {
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "ref:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return ""
		}

		return strings.TrimPrefix(fields[1], "refs/heads/")
	}

	return ""
}

// remoteBranchSource returns the canonical BranchSource label used
// by RepoTermBlock when the branch came from `ls-remote --symref`.
// Empty branch maps to empty source so the renderer drops the
// parenthesized segment entirely.
func remoteBranchSource(branch string) string {
	if len(strings.TrimSpace(branch)) == 0 {
		return ""
	}

	return constants.BranchSourceRemoteHEAD
}
