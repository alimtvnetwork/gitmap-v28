// Package cmd — cfrppriorversion.go: spec/01-app/113 §2.3.
//
// After cfrp publishes vN, probe v(N-1), v(N-2), … on the same
// provider+owner. Any that are currently `public` are offered up
// for privatization. With `-y` we auto-confirm; otherwise prompt.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
)

// runCFRPPriorVersionPrivatize is invoked after the make-public step
// in cfrp succeeds. autoYes mirrors the cfrp `-y` flag.
func runCFRPPriorVersionPrivatize(absPath string, autoYes bool) {
	base, current := resolvePriorScanIdentity(absPath)
	if len(base) == 0 || current < 2 {
		return
	}

	ctx, ok := resolvePriorScanProvider(absPath)
	if !ok {
		return
	}

	fmt.Printf(constants.MsgCFRPPriorHeaderFmt, base, constants.CFRPPriorMaxLookback)

	publicSlugs := scanPriorPublicSlugs(ctx, base, current)
	if len(publicSlugs) == 0 {
		fmt.Print(constants.MsgCFRPPriorNoneFound)

		return
	}

	fmt.Printf(constants.MsgCFRPPriorFoundFmt, len(publicSlugs),
		strings.Join(publicSlugs, ", "))

	if !autoYes && !promptPrivatize(len(publicSlugs)) {
		fmt.Print(constants.MsgCFRPPriorSkipped)

		return
	}

	privatizeSlugs(ctx, publicSlugs)
}

// resolvePriorScanIdentity returns (baseName, currentVer). Returns
// ("", 0) when the repo identity is unversioned or unresolvable —
// the caller treats that as "nothing to scan".
func resolvePriorScanIdentity(absPath string) (string, int) {
	remoteURL, err := gitutil.RemoteURL(absPath)
	if err != nil {
		return "", 0
	}

	repoName := repoNameFromURL(remoteURL)
	parsed := clonenext.ParseRepoName(repoName)
	if !parsed.HasVersion {
		return "", 0
	}

	return parsed.BaseName, parsed.CurrentVersion
}

// resolvePriorScanProvider builds a visibilityContext for the just-
// cloned repo so we know which provider CLI to invoke against the
// owner's prior-version slugs.
func resolvePriorScanProvider(absPath string) (visibilityContext, bool) {
	remoteURL, err := gitutil.RemoteURL(absPath)
	if err != nil || len(remoteURL) == 0 {
		return visibilityContext{}, false
	}

	provider := classifyProvider(remoteURL)
	slug := parseOwnerRepo(remoteURL)
	if len(provider) == 0 || len(slug) == 0 {
		return visibilityContext{}, false
	}

	if _, lookErr := exec.LookPath(providerCLI(provider)); lookErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrVisCLIMissingFmt, providerCLI(provider))

		return visibilityContext{}, false
	}

	return visibilityContext{URL: remoteURL, Provider: provider, Slug: slug}, true
}

// scanPriorPublicSlugs walks vN-1, vN-2, … and collects the slugs
// whose current visibility is public. Stops after two consecutive
// 404/error responses to avoid hammering the API for non-existent
// versions.
func scanPriorPublicSlugs(ctx visibilityContext, base string, current int) []string {
	owner := ownerFromSlug(ctx.Slug)
	var publicSlugs []string
	misses := 0

	for i := 1; i <= constants.CFRPPriorMaxLookback; i++ {
		ver := current - i
		if ver < 1 {
			break
		}

		slug := fmt.Sprintf("%s/%s-v%d", owner, base, ver)
		sub := visibilityContext{Provider: ctx.Provider, Slug: slug, URL: ctx.URL}

		state, err := readVisibilitySoft(sub, false)
		if err != nil || len(state) == 0 {
			misses++
			if misses >= 2 {
				break
			}

			continue
		}

		misses = 0
		if state == constants.VisibilityPublic {
			publicSlugs = append(publicSlugs, slug)
		}
	}

	return publicSlugs
}

// promptPrivatize asks for [y/N] confirmation. Default is no — we
// only privatize when the user types exactly `y` or `yes`.
func promptPrivatize(count int) bool {
	fmt.Printf(constants.MsgCFRPPriorPromptFmt, count)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	answer := strings.TrimSpace(strings.ToLower(line))

	return answer == "y" || answer == "yes"
}

// privatizeSlugs runs `gh repo edit --visibility private` for each
// slug, printing one-line results. Failures are non-fatal so a
// transient API hiccup on v(N-3) doesn't strand v(N-2) and v(N-1).
func privatizeSlugs(parent visibilityContext, slugs []string) {
	opts := visibilityFlags{}

	for i, slug := range slugs {
		fmt.Printf(constants.MsgVisBulkItemFmt, i+1, len(slugs), slug)

		sub := visibilityContext{Provider: parent.Provider, Slug: slug, URL: parent.URL}
		_ = applyAndReport(sub, constants.VisibilityPublic,
			constants.VisibilityPrivate, opts)
	}
}
