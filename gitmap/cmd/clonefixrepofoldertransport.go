package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// preferExistingFolderTransport closes the cfr/cfrp gap documented in
// .lovable/audits/2026-06-07-reclone-pickers.md: when the destination
// folder already contains a `.git/` directory, the existing origin's
// transport is the source of truth — re-cloning with a different
// scheme silently downgrades SSH-origin repos to HTTPS and re-triggers
// the browser-auth prompt on private remotes. Returns the (possibly
// rewritten) URL the actual `git clone` should use.
//
// No-ops when:
//   - absPath has no `.git/` (fresh clone — nothing to honor).
//   - `gitutil.RemoteURL` fails (warned to stderr, original URL kept).
//   - existing origin and positional URL already share a transport.
//   - URL rewrite fails (warned to stderr, original URL kept).
func preferExistingFolderTransport(positional, absPath string) string {
	if !hasDotGitDir(absPath) {
		return positional
	}
	existing, err := gitutil.RemoteURL(absPath)
	if err != nil || existing == "" {
		warnPreferTransport(absPath, "could not read existing origin", err)
		return positional
	}
	posIsSSH := isSSHURL(positional)
	exIsSSH := isSSHURL(existing)
	if posIsSSH == exIsSSH {
		return positional
	}
	return rewriteToMatchExisting(positional, exIsSSH)
}

// rewriteToMatchExisting flips the positional URL to the existing
// folder's transport. Logs the swap to stderr so the user can audit it.
func rewriteToMatchExisting(positional string, targetSSH bool) string {
	if targetSSH {
		out, ok := ConvertURLToSSH(positional)
		if !ok {
			warnPreferTransport("", "ssh rewrite failed", nil)
			return positional
		}
		fmt.Fprintf(os.Stderr, constants.MsgCFRFolderTransport, "ssh", positional, out)
		return out
	}
	out, ok := ConvertURLToHTTPS(positional)
	if !ok {
		warnPreferTransport("", "https rewrite failed", nil)
		return positional
	}
	fmt.Fprintf(os.Stderr, constants.MsgCFRFolderTransport, "https", positional, out)
	return out
}

// hasDotGitDir reports whether absPath contains a `.git` entry (dir
// or worktree file). Mirrors isGitRepoDir in reporeclone.go but is
// duplicated here to avoid coupling cfr to the reclone package.
func hasDotGitDir(absPath string) bool {
	_, err := os.Stat(filepath.Join(absPath, ".git"))
	return err == nil
}

// isSSHURL classifies a URL as SSH (shorthand or scheme). Anything
// else — including https://, http://, file://, or unrecognised — is
// treated as not-SSH so the caller defaults to the HTTPS-friendly
// branch.
func isSSHURL(url string) bool {
	lower := strings.ToLower(strings.TrimSpace(url))
	if strings.HasPrefix(lower, "git@") {
		return true
	}
	return strings.HasPrefix(lower, "ssh://")
}

// warnPreferTransport surfaces non-fatal transport-detection failures
// to stderr per the zero-swallow error policy. The clone still
// proceeds with the original URL — fail-open is the right default
// because the user's positional URL might still be correct.
func warnPreferTransport(absPath, reason string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnCFRFolderTransport, absPath, reason, err)
		return
	}
	fmt.Fprintf(os.Stderr, constants.WarnCFRFolderTransportNoErr, absPath, reason)
}
