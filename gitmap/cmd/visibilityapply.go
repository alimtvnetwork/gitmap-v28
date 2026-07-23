// Package cmd — visibilityapply.go: read / write / verify visibility
// via the host provider's CLI (`gh repo view --json visibility` and
// `gh repo edit --visibility ...`; analogous `glab` commands).
//
// We deliberately shell out instead of speaking the REST API
// directly — it lets users authenticate once via their normal
// `gh auth login` / `glab auth login` flow and keeps gitmap free
// of token storage.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// mustReadCurrentVisibility runs the provider CLI to fetch the
// current visibility. Returns "public" or "private"; exits
// ExitVisAuthFailed on any auth/API failure.
func mustReadCurrentVisibility(ctx visibilityContext, verbose bool) string {
	args := readVisibilityArgs(ctx.Provider, ctx.Slug)
	out, err := runProviderCLI(ctx.Provider, args, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrVisReadCurrentFmt, providerCLI(ctx.Provider), err)
		os.Exit(constants.ExitVisAuthFailed)
	}

	return parseVisibilityOutput(ctx.Provider, out)
}

// readVisibilityArgs builds the argv for "show me the current
// visibility" on the given provider.
func readVisibilityArgs(provider, slug string) []string {
	if provider == constants.ProviderGitLab {
		// glab repo view <slug> -F json
		return []string{"repo", "view", slug, "-F", "json"}
	}

	return []string{"repo", "view", slug, "--json", "visibility"}
}

// parseVisibilityOutput normalises the JSON-ish output from each CLI
// to the lowercase "public"/"private" tokens we use internally.
//
// gh emits {"visibility":"PUBLIC"} (uppercase); glab emits
// {"visibility":"public"} (lowercase). We lowercase + substring-match
// to be tolerant of either form and of any extra fields.
func parseVisibilityOutput(_, out string) string {
	lower := strings.ToLower(out)
	switch {
	case strings.Contains(lower, `"visibility":"public"`):
		return constants.VisibilityPublic
	case strings.Contains(lower, `"visibility":"private"`):
		return constants.VisibilityPrivate
	case strings.Contains(lower, `"visibility":"internal"`):
		// GitLab "internal" = signed-in users can see it. Treat as
		// private for our purposes (it is NOT publicly visible).
		return constants.VisibilityPrivate
	}

	return ""
}

// applyVisibilityOrExit runs `gh repo edit --visibility <target>`
// (or the glab equivalent). Captures stderr so we can surface a
// useful error message on failure.
func applyVisibilityOrExit(ctx visibilityContext, target string, verbose bool) {
	args := applyVisibilityArgs(ctx.Provider, ctx.Slug, target)
	stderr, err := runProviderCLICapturingStderr(ctx.Provider, args, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrVisApplyFailedFmt, err, stderr)
		os.Exit(constants.ExitVisAuthFailed)
	}
}

// applyVisibilityArgs builds the argv for "set visibility = target".
// `gh` requires --accept-visibility-change-consequences when going
// public→private to acknowledge the loss of stars/forks visibility.
func applyVisibilityArgs(provider, slug, target string) []string {
	if provider == constants.ProviderGitLab {
		// glab repo edit <slug> --visibility <target>
		return []string{"repo", "edit", slug, "--visibility", target}
	}

	return []string{
		"repo", "edit", slug,
		"--visibility", target,
		"--accept-visibility-change-consequences",
	}
}

// verifyVisibilityOrExit re-reads visibility after the apply step
// and exits ExitVisVerifyFailed if it does not match the target.
func verifyVisibilityOrExit(ctx visibilityContext, target string, verbose bool) {
	current := mustReadCurrentVisibility(ctx, verbose)
	if current != target {
		fmt.Fprintf(os.Stderr, constants.ErrVisVerifyFailedFmt, current, target)
		os.Exit(constants.ExitVisVerifyFailed)
	}
	fmt.Printf(constants.MsgVisVerifyOK, current)
}

// runProviderCLI invokes `gh`/`glab` and returns combined stdout.
// Verbose mode echoes the command line to stderr first.
func runProviderCLI(provider string, args []string, verbose bool) (string, error) {
	cli := providerCLI(provider)
	if verbose {
		fmt.Fprintf(os.Stderr, constants.MsgVisVerboseExec, cli, strings.Join(args, " "))
	}

	out, err := exec.Command(cli, args...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// runProviderCLICapturingStderr is like runProviderCLI but also
// returns stderr so callers can surface API-level error text.
func runProviderCLICapturingStderr(provider string, args []string, verbose bool) (string, error) {
	cli := providerCLI(provider)
	if verbose {
		fmt.Fprintf(os.Stderr, constants.MsgVisVerboseExec, cli, strings.Join(args, " "))
	}

	var stderr bytes.Buffer
	cmd := exec.Command(cli, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	return stderr.String(), err
}
