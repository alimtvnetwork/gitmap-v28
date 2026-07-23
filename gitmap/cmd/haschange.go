// Package cmd — haschange.go implements `gitmap has-change (hc) <repo>`.
//
// Prints "true" or "false" depending on whether the named repo has uncommitted
// changes (default), is ahead of origin, or is behind origin. Use --mode to
// switch dimensions; --all prints structured output covering all three.
//
// Examples:
//
//	gitmap hc gitmap                  -> true | false   (dirty working tree)
//	gitmap hc gitmap --mode=ahead     -> true | false   (local commits not pushed)
//	gitmap hc gitmap --mode=behind    -> true | false   (remote commits not pulled)
//	gitmap hc gitmap --all            -> dirty=true ahead=false behind=true
package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runHasChange handles the `has-change` command.
func runHasChange(args []string) {
	checkHelp(constants.CmdHasChange, args)

	alias, mode, all, fetch := parseHasChangeFlags(args)
	if len(alias) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrHCUsage)
		os.Exit(2)
	}

	target := resolveReleaseAliasPath(alias)

	if fetch {
		fetchRemoteIn(target)
	}

	if all {
		printHasChangeAll(target)

		return
	}

	printHasChangeOne(target, mode)
}

// parseHasChangeFlags extracts the repo alias and mode flags.
func parseHasChangeFlags(args []string) (alias, mode string, all, fetch bool) {
	fs := flag.NewFlagSet(constants.CmdHasChange, flag.ExitOnError)
	modeFlag := fs.String(constants.FlagHCMode, constants.HCModeDirty, constants.FlagDescHCMode)
	allFlag := fs.Bool(constants.FlagHCAll, false, constants.FlagDescHCAll)
	fetchFlag := fs.Bool(constants.FlagHCFetch, true, constants.FlagDescHCFetch)
	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(2)
	}
	rest := fs.Args()
	if len(rest) >= 1 {
		alias = rest[0]
	}

	return alias, *modeFlag, *allFlag, *fetchFlag
}

// printHasChangeOne prints true/false for a single dimension.
func printHasChangeOne(target, mode string) {
	switch mode {
	case constants.HCModeDirty:
		fmt.Println(boolStr(isWorkingTreeDirty(target)))
	case constants.HCModeAhead:
		ahead, _, _ := readAheadBehind(target)
		fmt.Println(boolStr(ahead > 0))
	case constants.HCModeBehind:
		_, behind, _ := readAheadBehind(target)
		fmt.Println(boolStr(behind > 0))
	default:
		fmt.Fprintf(os.Stderr, constants.ErrHCBadMode, mode)
		os.Exit(2)
	}
}

// printHasChangeAll prints structured output for all three dimensions.
func printHasChangeAll(target string) {
	dirty := isWorkingTreeDirty(target)
	ahead, behind, ok := readAheadBehind(target)
	if !ok {
		fmt.Printf(constants.MsgHCAllNoUpstream, boolStr(dirty))

		return
	}
	fmt.Printf(constants.MsgHCAllFmt, boolStr(dirty), boolStr(ahead > 0), boolStr(behind > 0))
}

// boolStr returns "true" or "false" for printing.
func boolStr(b bool) string {
	if b {
		return constants.HCTrue
	}

	return constants.HCFalse
}

// readAheadBehind runs `git rev-list --left-right --count HEAD...@{upstream}`
// in target and returns the (ahead, behind, ok) tuple. ok is false when there
// is no configured upstream.
func readAheadBehind(target string) (int, int, bool) {
	cmd := exec.Command(constants.GitBin, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = target
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, false
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) != 2 {
		return 0, 0, false
	}
	ahead, _ := strconv.Atoi(parts[0])
	behind, _ := strconv.Atoi(parts[1])

	return ahead, behind, true
}

// fetchRemoteIn runs `git fetch` in target. Errors are printed but do not
// fail the command — the user gets a stale ahead/behind result with a warning.
func fetchRemoteIn(target string) {
	cmd := exec.Command(constants.GitBin, "fetch")
	cmd.Dir = target
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnHCFetchFailed, target, err)
	}
}
