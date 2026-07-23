// Package cmd — visibilitybulk.go: spec/01-app/113 §2.2.
//
// Adds positional `make-public|make-private <repo-or-url> <count>`
// form. When count >= 1, the command flips the N most recent
// versions (vN, vN-1, …, vN-count+1) of the base repo on the same
// provider+owner as the current repo's origin.
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// bulkVisibilityRequest captures a parsed bulk invocation.
type bulkVisibilityRequest struct {
	BaseRepo string // base name with no -vN suffix
	StartVer int    // highest version to flip
	Count    int    // how many versions, descending from StartVer
}

// maybeRunBulkVisibility returns true when positional args define a
// bulk request and the bulk path was taken (the function then exits
// the process itself). Returns false to let the caller fall through
// to the legacy single-repo path.
func maybeRunBulkVisibility(positional []string, target string, opts visibilityFlags) bool {
	req, ok := parseBulkRequest(positional)
	if !ok {
		return false
	}

	ctx := mustResolveVisibilityContext()
	mustEnsureProviderCLI(ctx.Provider, opts.verbose)

	runBulkVisibility(ctx, req, target, opts)

	return true
}

// parseBulkRequest accepts either `<repo> <count>` or `<count>`.
// A single positional that parses as int reuses the current repo's
// base; a `<repo> <count>` pair targets a different repo by name.
func parseBulkRequest(positional []string) (bulkVisibilityRequest, bool) {
	switch len(positional) {
	case 1:
		return parseBulkSingleArg(positional[0])
	case 2:
		return parseBulkPairArg(positional[0], positional[1])
	}

	return bulkVisibilityRequest{}, false
}

func parseBulkSingleArg(arg string) (bulkVisibilityRequest, bool) {
	count, err := strconv.Atoi(arg)
	if err != nil || count < 1 {
		return bulkVisibilityRequest{}, false
	}

	base, ver := resolveCurrentRepoBaseAndVersion()

	return bulkVisibilityRequest{BaseRepo: base, StartVer: ver - 1, Count: count}, true
}

func parseBulkPairArg(repoArg, countArg string) (bulkVisibilityRequest, bool) {
	count, err := strconv.Atoi(countArg)
	if err != nil || count < 1 {
		fmt.Fprintf(os.Stderr, constants.ErrVisBulkBadCountFmt, countArg)
		os.Exit(constants.ExitVisBadFlag)
	}

	base, ver := extractBaseAndVersionFromArg(repoArg)
	if len(base) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrVisBulkRepoParseFmt, repoArg)
		os.Exit(constants.ExitVisBadFlag)
	}

	return bulkVisibilityRequest{BaseRepo: base, StartVer: ver - 1, Count: count}, true
}

// resolveCurrentRepoBaseAndVersion reads origin and splits its repo
// identity into (base, version). Falls back to version=1 when the
// current repo is unversioned (caller should typically pass an
// explicit URL in that situation).
func resolveCurrentRepoBaseAndVersion() (string, int) {
	ctx := mustResolveVisibilityContext()

	return extractBaseAndVersionFromArg(ctx.URL)
}

// extractBaseAndVersionFromArg accepts a URL, a `owner/repo[-vN]`
// slug, or a bare `repo[-vN]` and returns (base, version). Version
// defaults to 1 when absent (single-repo flip).
func extractBaseAndVersionFromArg(arg string) (string, int) {
	repo := repoNameFromURL(arg)
	if len(repo) == 0 {
		repo = arg
	}

	parsed := clonenext.ParseRepoName(repo)
	if parsed.HasVersion {
		return parsed.BaseName, parsed.CurrentVersion
	}

	return parsed.BaseName, 1
}

// runBulkVisibility iterates from StartVer down for Count versions,
// flipping each to `target`. Records the worst non-zero exit code
// seen and exits with it at the end so the user sees a real failure
// signal even when most items succeed.
func runBulkVisibility(ctx visibilityContext, req bulkVisibilityRequest,
	target string, opts visibilityFlags) {

	fmt.Printf(constants.MsgVisBulkHeaderFmt, target, req.Count, req.BaseRepo, ctx.Provider)

	worst := constants.ExitVisOK
	owner := ownerFromSlug(ctx.Slug)

	for i := 0; i < req.Count; i++ {
		ver := req.StartVer - i
		if ver < 1 {
			break
		}

		fmt.Printf(constants.MsgVisBulkItemFmt, i+1, req.Count,
			fmt.Sprintf("%s/%s-v%d", owner, req.BaseRepo, ver))

		code := flipOneSlug(ctx.Provider, owner, req.BaseRepo, ver, target, opts)
		if code > worst {
			worst = code
		}
	}

	os.Exit(worst)
}

// flipOneSlug visits a single owner/base-vN slug, flipping to target.
// Returns the exit code that would have been raised had this been a
// single-repo command (0 on success or already-target, ExitVisAuthFailed
// on apply failure).
func flipOneSlug(provider, owner, base string, ver int,
	target string, opts visibilityFlags) int {

	slug := fmt.Sprintf("%s/%s-v%d", owner, base, ver)
	subCtx := visibilityContext{Provider: provider, Slug: slug,
		URL: fmt.Sprintf("https://%s/%s", providerHost(provider), slug)}

	if opts.dryRun {
		fmt.Printf(constants.MsgVisBulkDryFmt, "current", target, slug)

		return constants.ExitVisOK
	}

	current, err := readVisibilitySoft(subCtx, opts.verbose)
	if err != nil {
		fmt.Printf(constants.MsgVisBulkFailFmt, err)

		return constants.ExitVisAuthFailed
	}

	if current == target {
		fmt.Printf(constants.MsgVisBulkSkipFmt, current)

		return constants.ExitVisOK
	}

	return applyAndReport(subCtx, current, target, opts)
}

// applyAndReport runs the provider edit + prints a one-line outcome.
func applyAndReport(ctx visibilityContext, current, target string,
	opts visibilityFlags) int {

	args := applyVisibilityArgs(ctx.Provider, ctx.Slug, target)
	if _, err := runProviderCLICapturingStderr(ctx.Provider, args, opts.verbose); err != nil {
		fmt.Printf(constants.MsgVisBulkFailFmt, err)

		return constants.ExitVisAuthFailed
	}

	fmt.Printf(constants.MsgVisBulkOKFmt, current, target)

	return constants.ExitVisOK
}
