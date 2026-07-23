// Package cmd — visibilityallbulk.go: top-level handler for the four
// bulk wildcard visibility commands (make-all-public / make-all-private
// / MAPUB / MAPRI) plus their except-latest counterparts. Owns dispatch,
// flag parsing (-Y / --verbose / --parallel / --cache-ttl / --except-latest),
// owner resolution, repo enumeration (TTL-cached), pattern matching,
// optional except-latest filtering, the optional interactive prompt
// (-Y skips it), the parallel per-repo apply loop, and exit-code
// aggregation.
//
// Heavy lifting is delegated:
//   - ResolveOwnerOnly           → visibilityresolveowner.go
//   - listOwnerReposCached       → visibilityownerlistcache.go
//   - visibility.ParsePatternList / MatchOwnerRepos → gitmap/visibility
//   - splitExceptLatest          → visibilityexceptlatest.go
//   - renderMatchedTable / promptConfirmOrExclude → visibilitybulkprompt.go
//   - applyBulkLoopParallel      → visibilityparallel.go
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §plan + §parallel.
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/visibility"
)

// bulkFlags holds the parsed CLI flags for a bulk visibility run.
type bulkFlags struct {
	Yes          bool
	Verbose      bool
	ExceptLatest bool
	Parallel     int
	CacheTTLSecs int
	CacheTTLSet  bool
}

// runMakeAllPublic / runMakeAllPrivate are the dispatcher entry points.
func runMakeAllPublic(args []string) {
	runMakeAllVisibility(constants.VisibilityPublic, constants.CmdMakeAllPublic, args, false)
}

func runMakeAllPrivate(args []string) {
	runMakeAllVisibility(constants.VisibilityPrivate, constants.CmdMakeAllPrivate, args, false)
}

// Except-latest entry points pre-set the filter and reuse the rest
// of the pipeline so behavior stays in lock-step with the base
// commands.
func runMakeAllPublicExceptLatest(args []string) {
	runMakeAllVisibility(constants.VisibilityPublic, constants.CmdMakeAllPublicExceptLatest, args, true)
}

func runMakeAllPrivateExceptLatest(args []string) {
	runMakeAllVisibility(constants.VisibilityPrivate, constants.CmdMakeAllPrivateExceptLatest, args, true)
}

// runMakeAllVisibility orchestrates the full bulk run.
func runMakeAllVisibility(target, cmdName string, args []string, exceptLatestDefault bool) {
	checkHelp(cmdName, args)
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, constants.ErrMakeAllMissingArgFmt, cmdName)
		os.Exit(constants.ExitVisBadFlag)
	}

	ownerArg, patternsRaw, flags := parseBulkArgs(args)
	if exceptLatestDefault {
		flags.ExceptLatest = true
	}
	ctx := resolveOwnerOrExit(ownerArg)
	mustEnsureProviderCLI(ctx.Provider, flags.Verbose)
	mustEnsureProviderAuth(ctx.Provider, flags.Verbose)

	patterns, err := visibility.ParsePatternList(patternsRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: %v\n", err)
		os.Exit(constants.ExitVisBadFlag)
	}

	matches, ownerTotal := matchOrExitEmpty(ctx, patterns, flags)
	var latestInvert []visibility.MatchedRepo
	inverted := invertVisibility(target)
	if flags.ExceptLatest {
		fmt.Fprint(os.Stdout, constants.MsgBulkExceptLatest)
		matches, latestInvert = splitExceptLatest(matches, os.Stdout, inverted)
		if len(matches) == 0 && len(latestInvert) == 0 {
			fmt.Fprint(os.Stderr, constants.MsgBulkNoMatches)
			os.Exit(constants.ExitVisOK)
		}
	}
	combined := append(append([]visibility.MatchedRepo{}, matches...), latestInvert...)
	audit := beginRunAudit(ctx, target, cmdName, patternsRaw, flags, ownerTotal, combined)

	final := confirmOrAbort(combined, flags.Yes)
	if len(final) == 0 {
		excluded := audit.markExcluded(combined, nil)
		audit.finalize(excluded, 0, 0, 0, constants.ExitVisConfirmReq)
		fmt.Fprint(os.Stderr, constants.MsgBulkAborted)
		os.Exit(constants.ExitVisConfirmReq)
	}
	excludedCount := audit.markExcluded(combined, final)
	mainFinal, invertFinal := partitionByName(final, latestInvert)

	changed, skipped, failed := 0, 0, 0
	if len(mainFinal) > 0 {
		fmt.Fprintf(os.Stdout, constants.MsgBulkApplyHeaderFmt, target, len(mainFinal), ctx.Owner)
		c, s, f := applyBulkLoopParallel(ctx, target, mainFinal, flags, audit)
		changed, skipped, failed = changed+c, skipped+s, failed+f
	}
	if len(invertFinal) > 0 {
		fmt.Fprintf(os.Stdout, constants.MsgBulkInvertHeaderFmt, inverted, len(invertFinal), ctx.Owner)
		c, s, f := applyBulkLoopParallel(ctx, inverted, invertFinal, flags, audit)
		changed, skipped, failed = changed+c, skipped+s, failed+f
	}
	fmt.Fprintf(os.Stdout, constants.MsgBulkSummaryFmt, changed, skipped, failed, len(final))
	exit := bulkExitCode(changed, failed)
	audit.finalize(excludedCount, changed, skipped, failed, exit)
	os.Exit(exit)
}

// partitionByName splits `final` into (mainSet, invertSet) using the
// names present in `invertSource` as the membership test for invertSet.
func partitionByName(final, invertSource []visibility.MatchedRepo) ([]visibility.MatchedRepo, []visibility.MatchedRepo) {
	if len(invertSource) == 0 {
		return final, nil
	}
	isInvert := make(map[string]bool, len(invertSource))
	for _, m := range invertSource {
		isInvert[m.RepoName] = true
	}
	main := make([]visibility.MatchedRepo, 0, len(final))
	inv := make([]visibility.MatchedRepo, 0, len(invertSource))
	for _, m := range final {
		if isInvert[m.RepoName] {
			inv = append(inv, m)
		} else {
			main = append(main, m)
		}
	}

	return main, inv
}

// parseBulkArgs splits owner / pattern-list / flags. Accepts the
// legacy -Y/-y/--yes/--verbose plus the new --parallel=N, --cache-ttl=N,
// and --except-latest/-XL flags anywhere after the first two positional
// args. Unknown flags are ignored (forwards-compatible with future spec
// additions).
func parseBulkArgs(args []string) (string, string, bulkFlags) {
	flags := bulkFlags{Parallel: constants.DefaultBulkParallelism}
	for _, a := range args[2:] {
		switch {
		case a == "-Y" || a == "-y" || a == "--yes":
			flags.Yes = true
		case a == "--verbose":
			flags.Verbose = true
		case a == constants.FlagBulkExceptLatest || a == constants.FlagBulkExceptLatestShort:
			flags.ExceptLatest = true
		case strings.HasPrefix(a, constants.FlagBulkParallel+"="):
			if n, err := strconv.Atoi(strings.TrimPrefix(a, constants.FlagBulkParallel+"=")); err == nil && n > 0 {
				if n > constants.MaxBulkParallelism {
					n = constants.MaxBulkParallelism
				}
				flags.Parallel = n
			}
		case strings.HasPrefix(a, constants.FlagBulkCacheTTL+"="):
			if n, err := strconv.Atoi(strings.TrimPrefix(a, constants.FlagBulkCacheTTL+"=")); err == nil && n >= 0 {
				flags.CacheTTLSecs = n
				flags.CacheTTLSet = true
			}
		}
	}

	return args[0], args[1], flags
}

// resolveOwnerOrExit wraps ResolveOwnerOnly with Code Red exit handling.
func resolveOwnerOrExit(arg string) ownerContext {
	ctx, err := ResolveOwnerOnly(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrMakeAllResolveFmt, err)
		os.Exit(constants.ExitVisBadProvider)
	}

	return ctx
}

// matchOrExitEmpty lists owner repos (via TTL cache), matches patterns,
// exits 0 with a friendly message when nothing matched. Returns the
// matched subset AND the owner-wide total so the audit layer can
// persist OwnerRepoTotal.
func matchOrExitEmpty(ctx ownerContext, patterns []visibility.Pattern, flags bulkFlags) ([]visibility.MatchedRepo, int) {
	names, err := listOwnerReposCached(ctx.Provider, ctx.Owner, flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: %v\n", err)
		os.Exit(constants.ExitVisAuthFailed)
	}

	matches := visibility.MatchOwnerRepos(names, patterns)
	if len(matches) == 0 {
		patterns, matches = fuzzyFallback(patterns, names)
	}
	fmt.Fprint(os.Stdout, renderMatchedTable(ctx.Owner, len(names), matches))
	if len(matches) == 0 {
		printNearMissHints(patterns, names)
		fmt.Fprint(os.Stderr, constants.MsgBulkNoMatches)
		os.Exit(constants.ExitVisOK)
	}

	return matches, len(names)
}

// fuzzyFallback runs visibility.AutoFixVDigitPatterns and, if it
// added any patterns, re-runs MatchOwnerRepos. Emits a stdout note
// so the user sees exactly which token was auto-fixed.
func fuzzyFallback(patterns []visibility.Pattern, names []string) ([]visibility.Pattern, []visibility.MatchedRepo) {
	extra := visibility.AutoFixVDigitPatterns(patterns, names)
	if len(extra) == 0 {
		return patterns, nil
	}
	for _, p := range extra {
		fmt.Fprintf(os.Stdout, constants.MsgBulkFuzzyAutoFixFmt, p.Raw)
	}
	merged := append(patterns, extra...)

	return merged, visibility.MatchOwnerRepos(names, merged)
}

// printNearMissHints surfaces up to 5 close repo names so users can
// retry with the right token. Stderr-only; no prompt to keep this
// non-interactive path scriptable.
func printNearMissHints(patterns []visibility.Pattern, names []string) {
	hints := visibility.NearMisses(patterns, names, 3, 5)
	if len(hints) == 0 {
		return
	}
	fmt.Fprint(os.Stderr, constants.MsgBulkFuzzyHintHeader)
	for _, h := range hints {
		fmt.Fprintf(os.Stderr, "  - %s\n", h)
	}
}

// confirmOrAbort runs the interactive prompt unless -Y was passed.
func confirmOrAbort(matches []visibility.MatchedRepo, yes bool) []visibility.MatchedRepo {
	if yes {
		return matches
	}
	final, proceed := promptConfirmOrExclude(os.Stdin, os.Stdout, matches)
	if !proceed {
		return nil
	}

	return final
}

// bulkExitCode collapses the tallies into the spec's exit-code matrix.
func bulkExitCode(changed, failed int) int {
	if failed == 0 {
		return constants.ExitVisOK
	}
	if changed == 0 {
		return constants.ExitVisAuthFailed
	}

	return constants.ExitVisBulkPartial
}
