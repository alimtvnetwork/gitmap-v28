package clonenow

// Plain-text parser. The text artifact `gitmap scan` writes is one
// line per repo, each line being a runnable `git clone <url> <dir>`
// command (with `<dir>` optional and `<url>` either https or ssh).
// Parsing is deliberately tolerant: blank lines and `#` comments are
// skipped so users can hand-edit the file before re-running.

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// parseTextRows reads the scan-text artifact (one `git clone ...`
// line per repo) into Row slices. Branch is always empty because
// the text format carries no branch information; users who need
// per-row branch pinning should use the JSON or CSV input.
func parseTextRows(r io.Reader) ([]Row, error) {
	scanner := bufio.NewScanner(r)
	// 1 MiB line cap -- scan output lines stay well under this even
	// for monorepo-style nested URLs, but the default 64 KiB would
	// truncate pathological cases without warning.
	scanner.Buffer(make([]byte, 0, 1<<16), 1<<20)
	out := make([]Row, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		row, ok := textRowFromLine(line)
		if !ok {
			continue
		}
		out = append(out, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(constants.ErrCloneNowTextRead, err)
	}

	return out, nil
}

// textRowFromLine extracts a Row from a single `git clone ...` line.
// Returns (_, false) for lines that don't look like a clone command
// so we silently skip section headers / shebangs / etc. -- a strict
// "every non-comment line must be a clone" rule would break the
// PowerShell-script artifact format, which interleaves comments and
// metadata that aren't worth special-casing here.
func textRowFromLine(line string) (Row, bool) {
	fields := strings.Fields(line)
	url, dest, ok := extractCloneArgs(fields)
	if !ok {
		return Row{}, false
	}
	row := Row{
		RepoName:     deriveRepoName(url),
		RelativePath: dest,
	}
	if isSSHURL(url) {
		row.SSHUrl = url
	} else {
		row.HTTPSUrl = url
	}
	if len(row.RelativePath) == 0 {
		row.RelativePath = DeriveDest(url)
	}

	return row, true
}

// extractCloneArgs walks the tokenized line looking for `git clone
// <url> [dest]`. Skips leading flags (e.g., `-b main`) so a user-
// edited file with shallow / branch flags still parses cleanly --
// we drop those flags on purpose so clone-now's mode + path
// contract isn't quietly subverted by leftover scan-time options.
func extractCloneArgs(fields []string) (string, string, bool) {
	if len(fields) < 3 || fields[0] != constants.GitBin || fields[1] != constants.GitClone {
		return "", "", false
	}
	rest := fields[2:]
	rest = skipCloneFlags(rest)
	if len(rest) == 0 {
		return "", "", false
	}
	url := rest[0]
	dest := ""
	if len(rest) > 1 {
		dest = rest[1]
	}

	return url, dest, true
}

// skipCloneFlags drops leading `-x` / `--foo` / `--foo=bar` tokens
// (and the value of `-b <branch>` style two-token flags) so the URL
// always lands at index 0 of the returned slice.
func skipCloneFlags(toks []string) []string {
	for len(toks) > 0 && strings.HasPrefix(toks[0], "-") {
		// Two-token flags: -b <branch>, --branch <branch>.
		if (toks[0] == "-b" || toks[0] == "--branch") && len(toks) > 1 {
			toks = toks[2:]

			continue
		}
		toks = toks[1:]
	}

	return toks
}

// isSSHURL classifies a URL as ssh-shaped. Accepts both `ssh://...`
// and the scp-style `user@host:path` form that git also speaks.
func isSSHURL(url string) bool {
	if strings.HasPrefix(url, "ssh://") {
		return true
	}
	// scp-style: must contain '@' before ':' and not be http(s)://
	if strings.Contains(url, "://") {
		return false
	}
	at := strings.Index(url, "@")
	colon := strings.Index(url, ":")

	return at >= 0 && colon > at
}

// deriveRepoName returns a human-friendly repo label from a URL --
// the last path segment, with `.git` stripped. Used for progress
// lines when the input format doesn't carry an explicit RepoName.
func deriveRepoName(url string) string {
	dest := DeriveDest(url)

	return dest
}

// DeriveDest returns the directory git would create by default for
// the given URL: the trailing path segment, with a single trailing
// `.git` removed. Public so executor + parser agree on the fallback.
func DeriveDest(url string) string {
	url = strings.TrimRight(url, "/")
	// scp-style splits on ':' -- the path is whatever follows.
	if i := strings.LastIndex(url, ":"); i >= 0 && !strings.Contains(url[:i], "/") {
		url = url[i+1:]
	}
	if i := strings.LastIndex(url, "/"); i >= 0 {
		url = url[i+1:]
	}

	return strings.TrimSuffix(url, ".git")
}
