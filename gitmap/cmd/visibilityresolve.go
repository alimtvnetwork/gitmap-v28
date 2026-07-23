// Package cmd — visibilityresolve.go: provider/slug detection,
// CLI-availability checks, and the interactive confirm prompt.
//
// Kept separate from visibility.go so each file stays well under
// the 200-line limit and helpers can be unit-tested in isolation
// without dragging in the flag-parsing surface.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// visibilityContext bundles the data resolved from the current repo.
// Provider is "github" or "gitlab"; Slug is "owner/repo".
type visibilityContext struct {
	URL      string
	Provider string
	Slug     string
}

// mustResolveVisibilityContext reads origin + classifies the host.
// Exits with the matching exit code on any failure (parity with the
// PowerShell reference).
func mustResolveVisibilityContext() visibilityContext {
	if !release.IsInsideGitRepo() {
		fmt.Fprint(os.Stderr, constants.ErrVisNotInRepo)
		os.Exit(constants.ExitVisNotARepo)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, constants.ErrVisNotInRepo)
		os.Exit(constants.ExitVisNotARepo)
	}

	url, err := gitutil.RemoteURL(cwd)
	if err != nil || len(url) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrVisNoOrigin)
		os.Exit(constants.ExitVisNoOrigin)
	}

	return resolveProviderAndSlugOrExit(url)
}

// resolveProviderAndSlugOrExit classifies the URL and parses the slug.
// Exits with ExitVisBadProvider if anything is unrecognized. Local
// remotes (file://, bare filesystem paths) warn-and-skip with exit 0
// so CI fixtures using local bare repos never fail the visibility step.
func resolveProviderAndSlugOrExit(url string) visibilityContext {
	if isLocalRemote(url) {
		fmt.Fprintf(os.Stderr, constants.MsgVisLocalSkipFmt, url)
		os.Exit(constants.ExitVisOK)
	}

	provider := classifyProvider(url)
	if len(provider) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrVisBadProviderFmt, url)
		os.Exit(constants.ExitVisBadProvider)
	}

	slug := parseOwnerRepo(url)
	if len(slug) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrVisBadSlugFmt, url)
		os.Exit(constants.ExitVisBadProvider)
	}

	return visibilityContext{URL: url, Provider: provider, Slug: slug}
}

// isLocalRemote reports whether url refers to a filesystem-local
// remote — file:// scheme, an absolute Unix path, or a Windows drive
// path. Local remotes have no provider API and should warn-and-skip.
func isLocalRemote(url string) bool {
	lower := strings.ToLower(strings.TrimSpace(url))
	if strings.HasPrefix(lower, "file://") {
		return true
	}
	if strings.HasPrefix(lower, "/") {
		return true
	}
	if len(lower) >= 3 && lower[1] == ':' && (lower[2] == '/' || lower[2] == '\\') {
		return true
	}

	return false
}


// classifyProvider returns "github", "gitlab", or "" for unknown.
// Matches both HTTPS and SSH URL forms.
func classifyProvider(url string) string {
	lower := strings.ToLower(url)
	if strings.Contains(lower, constants.HostGitHub) {
		return constants.ProviderGitHub
	}
	if strings.Contains(lower, constants.HostGitLab) {
		return constants.ProviderGitLab
	}

	return ""
}

// parseOwnerRepo returns "owner/repo" stripped of any .git suffix
// or trailing slash. Handles SSH (git@host:owner/repo[.git]) and
// HTTPS (https://host/owner/repo[.git]) forms.
func parseOwnerRepo(url string) string {
	trimmed := strings.TrimSuffix(strings.TrimSpace(url), "/")
	trimmed = strings.TrimSuffix(trimmed, ".git")

	if idx := strings.Index(trimmed, "@"); idx >= 0 {
		if colon := strings.Index(trimmed[idx:], ":"); colon >= 0 {
			return trimmed[idx+colon+1:]
		}
	}

	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return ""
	}

	return parts[len(parts)-2] + "/" + parts[len(parts)-1]
}

// mustEnsureProviderCLI exits ExitVisAuthFailed if `gh` / `glab`
// is not on PATH. Verbose mode prints the lookup before checking.
func mustEnsureProviderCLI(provider string, verbose bool) {
	cli := providerCLI(provider)
	if verbose {
		fmt.Fprintf(os.Stderr, constants.MsgVisVerboseExec, "which", cli)
	}
	if _, err := exec.LookPath(cli); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrVisCLIMissingFmt, cli)
		os.Exit(constants.ExitVisAuthFailed)
	}
}

// providerCLI maps a provider token to its CLI binary name.
func providerCLI(provider string) string {
	if provider == constants.ProviderGitLab {
		return constants.CLIGitLab
	}

	return constants.CLIGitHub
}

// confirmPublicOrExit prompts for explicit "yes" before exposing a
// private repo. Any other input (or EOF) exits ExitVisConfirmReq.
func confirmPublicOrExit(ctx visibilityContext) {
	fmt.Printf(constants.MsgVisConfirmFmt, ctx.Slug, ctx.Provider)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprint(os.Stderr, constants.ErrVisConfirmRequired)
		os.Exit(constants.ExitVisConfirmReq)
	}

	if strings.TrimSpace(strings.ToLower(line)) != "yes" {
		fmt.Fprint(os.Stderr, constants.ErrVisConfirmRequired)
		os.Exit(constants.ExitVisConfirmReq)
	}
}
