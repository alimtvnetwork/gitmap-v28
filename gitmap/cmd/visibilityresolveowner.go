// Package cmd — visibilityresolveowner.go: owner-only resolver used by
// the bulk wildcard visibility commands (make-all-public,
// make-all-private, MAPUB, MAPRI). Accepts a full provider URL, a bare
// "host/owner" token, a folder path, or "." (origin of cwd). Returns the
// classified provider and the bare owner — NO repo slug, because the
// caller will enumerate repos under that owner.
//
// Kept in its own file to honor the ≤200-line per-file rule and to keep
// the existing single-repo resolver (visibilityresolve.go) untouched.
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §2.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
)

// ownerContext is the owner-scoped sibling of visibilityContext.
type ownerContext struct {
	Provider  string // "github" | "gitlab"
	Owner     string // bare owner/org slug
	TargetRaw string // exact arg as typed, for audit persistence
}

// ResolveOwnerOnly classifies the supplied target and extracts the
// owner. Order of attempts: (1) explicit URL → host + first path
// segment; (2) bare "host/owner"; (3) folder path or "." → read
// origin of that folder's .git/config. Returns an error with full
// path/operation/reason context per Code Red rule on failure.
func ResolveOwnerOnly(arg string) (ownerContext, error) {
	trimmed := strings.TrimSpace(arg)
	for strings.HasSuffix(trimmed, "/") {
		trimmed = strings.TrimSuffix(trimmed, "/")
	}
	if len(trimmed) == 0 {
		return ownerContext{}, fmt.Errorf("Error: empty target (operation: resolve-owner, reason: arg is blank)")
	}

	if isURLLike(trimmed) {
		return ownerFromURL(trimmed)
	}

	if isBareHostOwner(trimmed) {
		return ownerFromBareHostOwner(trimmed)
	}

	return ownerFromFolder(trimmed)
}

// isURLLike returns true for https://, http://, git@host:owner forms.
func isURLLike(arg string) bool {
	lower := strings.ToLower(arg)
	if strings.HasPrefix(lower, "https://") {
		return true
	}
	if strings.HasPrefix(lower, "http://") {
		return true
	}

	return strings.HasPrefix(lower, "git@")
}

// isBareHostOwner matches "<host>/<owner>" with exactly one slash and a
// dot-bearing host segment (rules out plain folder names like "./acme").
func isBareHostOwner(arg string) bool {
	if strings.Contains(arg, "://") {
		return false
	}
	if strings.HasPrefix(arg, ".") {
		return false
	}
	parts := strings.Split(arg, "/")
	if len(parts) != 2 {
		return false
	}

	return strings.Contains(parts[0], ".") && len(parts[1]) > 0
}

// ownerFromURL extracts host + first path segment from a full URL.
func ownerFromURL(url string) (ownerContext, error) {
	provider := classifyProvider(url)
	if len(provider) == 0 {
		return ownerContext{}, fmt.Errorf("Error: unrecognized provider at %s (operation: classify-host, reason: not github.com or gitlab.com)", url)
	}

	owner := firstPathSegment(url)
	if len(owner) == 0 {
		return ownerContext{}, fmt.Errorf("Error: missing owner in URL %s (operation: parse-url, reason: empty path segment)", url)
	}

	return ownerContext{Provider: provider, Owner: owner, TargetRaw: url}, nil
}

// firstPathSegment pulls the first non-empty path component (the
// owner/org segment) from any supported URL form, stripping ".git"
// and any leading/trailing slashes. Skips the scheme and host so
// https://github.com/alice and https://github.com/alice/ both yield
// "alice" (NOT "github.com").
func firstPathSegment(url string) string {
	trimmed := strings.TrimSpace(url)
	for strings.HasSuffix(trimmed, "/") {
		trimmed = strings.TrimSuffix(trimmed, "/")
	}
	trimmed = strings.TrimSuffix(trimmed, ".git")

	// git@host:owner/repo form — owner is the segment after ':'.
	if idx := strings.Index(trimmed, "@"); idx >= 0 && !strings.Contains(trimmed[:idx], "/") {
		if colon := strings.Index(trimmed[idx:], ":"); colon >= 0 {
			rest := trimmed[idx+colon+1:]

			return strings.Split(strings.TrimLeft(rest, "/"), "/")[0]
		}
	}

	// Strip scheme (https://, http://, ssh://, git://).
	if schemeIdx := strings.Index(trimmed, "://"); schemeIdx >= 0 {
		trimmed = trimmed[schemeIdx+3:]
	}

	// Now trimmed is "<host>[/<owner>[/<repo>...]]" — drop the host.
	parts := strings.Split(trimmed, "/")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			return parts[i]
		}
	}

	return ""
}

// ownerFromBareHostOwner handles "github.com/acme" style.
func ownerFromBareHostOwner(arg string) (ownerContext, error) {
	parts := strings.SplitN(arg, "/", 2)
	host, owner := parts[0], parts[1]
	provider := classifyProvider(host)
	if len(provider) == 0 {
		return ownerContext{}, fmt.Errorf("Error: unrecognized host %s (operation: classify-host, reason: not github.com or gitlab.com)", host)
	}

	return ownerContext{Provider: provider, Owner: owner, TargetRaw: arg}, nil
}

// ownerFromFolder reads the origin URL from a folder's .git/config and
// recurses through ownerFromURL.
func ownerFromFolder(path string) (ownerContext, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return ownerContext{}, fmt.Errorf("Error: failed to resolve folder at %s: %v (operation: filepath.Abs, reason: %s)", path, err, err.Error())
	}

	if _, statErr := os.Stat(abs); statErr != nil {
		return ownerContext{}, fmt.Errorf("Error: folder not found at %s: %v (operation: stat, reason: %s)", abs, statErr, statErr.Error())
	}

	url, err := gitutil.RemoteURL(abs)
	if err != nil || len(url) == 0 {
		return ownerContext{}, fmt.Errorf("Error: no origin remote at %s: %v (operation: gitutil.RemoteURL, reason: %s)", abs, err, errString(err))
	}

	ctx, err := ownerFromURL(url)
	if err != nil {
		return ownerContext{}, err
	}
	ctx.TargetRaw = path

	return ctx, nil
}

// errString returns "<nil>" for a nil error so the Code Red reason field
// always renders a non-empty token. Avoids a nil-deref in callers.
func errString(err error) string {
	if err == nil {
		return constants.ProviderUnknownReason
	}

	return err.Error()
}
