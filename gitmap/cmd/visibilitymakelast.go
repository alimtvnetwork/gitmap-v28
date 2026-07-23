// Package cmd — visibilitymakelast.go: `make-last-public` /
// `make-last-private` (aliases MLPUB / MLPRI). Flips visibility on
// exactly one repo: the highest `-vN` sibling under `<base>` for the
// given owner.
//
// Resolution order:
//  1. If `<base>` itself ends in `-vN`, treat it as an exact repo
//     name and apply directly.
//  2. Else consult OwnerRepoNameIndex for the highest -vN row whose
//     BaseName == `<base>`.
//  3. Else refresh the owner repo list (warming the cache + index)
//     and retry step 2.
//
// Honors -Y / --yes to skip confirmation. Spec follow-up to
// spec/01-app/116-bulk-visibility-mapub-mapri.md.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/visibility"
)

func runMakeLastPublic(args []string) {
	runMakeLast(constants.VisibilityPublic, constants.CmdMakeLastPublic, args)
}
func runMakeLastPrivate(args []string) {
	runMakeLast(constants.VisibilityPrivate, constants.CmdMakeLastPrivate, args)
}

func runMakeLast(target, cmdName string, args []string) {
	checkHelp(cmdName, args)
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, constants.ErrMakeLastMissingArg, cmdName)
		os.Exit(constants.ExitVisBadFlag)
	}
	ownerArg, base, yes := parseMakeLastArgs(args)
	ctx := resolveOwnerOrExit(ownerArg)
	mustEnsureProviderCLI(ctx.Provider, false)
	mustEnsureProviderAuth(ctx.Provider, false)

	repoName, ver := resolveMakeLastRepo(ctx, base)
	if len(repoName) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrMakeLastNoBaseFmt, base, ctx.Owner, ctx.Owner, base)
		os.Exit(constants.ExitVisOK)
	}
	if ver >= 0 {
		fmt.Fprintf(os.Stdout, constants.MsgMakeLastResolvedFmt, base, repoName, ver, base)
	}

	if !yes && !confirmSingle(repoName, target) {
		fmt.Fprint(os.Stderr, constants.MsgBulkAborted)
		os.Exit(constants.ExitVisConfirmReq)
	}

	fmt.Fprintf(os.Stdout, constants.MsgBulkApplyHeaderFmt, target, 1, ctx.Owner)
	fmt.Fprintf(os.Stdout, constants.MsgBulkApplyItemFmt, 1, 1, repoName)
	status := applyOneRepo(ctx, repoName, target, false)
	changed, skipped, failed := tallyStatus(status)
	fmt.Fprintf(os.Stdout, constants.MsgBulkSummaryFmt, changed, skipped, failed, 1)
	os.Exit(bulkExitCode(changed, failed))
}

func parseMakeLastArgs(args []string) (string, string, bool) {
	yes := false
	for _, a := range args[2:] {
		if a == "-Y" || a == "-y" || a == "--yes" {
			yes = true
		}
	}

	return args[0], args[1], yes
}

// resolveMakeLastRepo returns (repoName, version). version is -1
// when the base was already an exact -vN form (no resolution needed).
func resolveMakeLastRepo(ctx ownerContext, base string) (string, int) {
	if _, _, ok := visibility.ParseRepoNameMeta(base); ok {
		return base, -1
	}
	if name, ver, ok := lookupIndexHighest(ctx.Provider, ctx.Owner, base); ok {
		return name, ver
	}
	// Cache miss — force refresh then retry.
	if _, err := listOwnerReposCached(ctx.Provider, ctx.Owner, bulkFlags{CacheTTLSecs: 0, CacheTTLSet: true}); err == nil {
		if name, ver, ok := lookupIndexHighest(ctx.Provider, ctx.Owner, base); ok {
			return name, ver
		}
	}
	// Fallback: scan names in memory (works even when the index
	// table isn't writable for some reason).
	names, err := listOwnerReposCached(ctx.Provider, ctx.Owner, bulkFlags{})
	if err == nil {
		if name, ver, ok := visibility.HighestVersionedMatch(names, base); ok {
			return name, ver
		}
	}

	return "", -1
}

func lookupIndexHighest(provider, owner, base string) (string, int, bool) {
	db, err := openDB()
	if err != nil {
		return "", 0, false
	}

	return db.LookupHighestVersion(provider, owner, base)
}

func confirmSingle(repoName, target string) bool {
	fmt.Fprintf(os.Stdout, "Apply visibility=%s to %s? [y/N]: ", target, repoName)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')

	return strings.EqualFold(strings.TrimSpace(line), "y")
}

func tallyStatus(s applyStatus) (int, int, int) {
	switch s.outcome {
	case "ok":
		return 1, 0, 0
	case "skip":
		return 0, 1, 0
	}

	return 0, 0, 1
}

// keep time imported for parseMakeLastArgs evolution.
var _ = time.Now
