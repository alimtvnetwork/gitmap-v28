// Package cmd — visibilitybulkhelpers.go: small helpers shared by
// visibilitybulk.go and cfrppriorversion.go. Kept in their own file
// to stay under the 200-line cap.
package cmdvisibility

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ownerFromSlug splits "owner/repo[-vN]" → "owner". Returns the slug
// unchanged when it has no slash (defensive — callers pre-validate).
func OwnerFromSlug(slug string) string {
	idx := strings.Index(slug, "/")
	if idx < 0 {
		return slug
	}

	return slug[:idx]
}

// providerHost maps the provider token back to its hostname.
func providerHost(provider string) string {
	if provider == constants.ProviderGitLab {
		return constants.HostGitLab
	}

	return constants.HostGitHub
}

// readVisibilitySoft is a non-exit version of mustReadCurrentVisibility.
// Returns the parsed state and any underlying CLI error so bulk
// callers can keep going on individual failures.
func ReadVisibilitySoft(ctx VisibilityContext, verbose bool) (string, error) {
	args := readVisibilityArgs(ctx.Provider, ctx.Slug)

	out, err := runProviderCLI(ctx.Provider, args, verbose)
	if err != nil {
		return "", fmt.Errorf("read visibility for %s: %w", ctx.Slug, err)
	}

	return parseVisibilityOutput(ctx.Provider, out), nil
}
