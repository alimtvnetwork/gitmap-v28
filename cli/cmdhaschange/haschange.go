// Package cmdhaschange implements `gitmap has-change (hc) <repo>`.
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
package cmdhaschange

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// RunHasChange handles the `has-change` command.
func RunHasChange(args []string) error {
	checkHelp(constants.CmdHasChange, args)

	alias, mode, hasAll, hasFetch, err := parseHasChangeFlags(args)
	if err != nil {
		return err
	}

	if len(alias) == 0 {
		return apperror.NewUsageError(constants.ErrHCUsage)
	}

	return executeHasChange(alias, mode, hasAll, hasFetch)
}

func executeHasChange(alias, mode string, hasAll, hasFetch bool) error {
	target := resolveReleaseAliasPath(alias)
	if hasFetch {
		_ = fetchRemoteIn(target)
	}

	if hasAll {
		return printHasChangeAll(target)
	}

	return printHasChangeOne(target, mode)
}

// parseHasChangeFlags extracts the repo alias and mode flags.
func parseHasChangeFlags(args []string) (string, string, bool, bool, error) {
	fs := flag.NewFlagSet(constants.CmdHasChange, flag.ContinueOnError)
	modeFlag := fs.String(constants.FlagHCMode, constants.HCModeDirty, constants.FlagDescHCMode)
	allFlag := fs.Bool(constants.FlagHCAll, false, constants.FlagDescHCAll)
	fetchFlag := fs.Bool(constants.FlagHCFetch, true, constants.FlagDescHCFetch)
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return "", "", false, false, apperror.WrapSimple(err, "parse has-change flags")
	}

	alias := extractAliasArg(fs.Args())

	return alias, *modeFlag, *allFlag, *fetchFlag, nil
}

func extractAliasArg(args []string) string {
	if len(args) >= 1 {
		return args[0]
	}

	return ""
}

// printHasChangeOne prints true/false for a single dimension.
func printHasChangeOne(target, mode string) error {
	if mode == constants.HCModeDirty {
		return printHasChangeDirty(target)
	}

	if mode == constants.HCModeAhead || mode == constants.HCModeBehind {
		return printHasChangeSync(target, mode)
	}

	return apperror.NewUsageError(fmt.Sprintf(constants.ErrHCBadMode, mode))
}

func printHasChangeDirty(target string) error {
	isDirty, err := isWorkingTreeDirty(target)
	if err != nil {
		return err
	}

	fmt.Println(boolStr(isDirty))

	return nil
}

func printHasChangeSync(target, mode string) error {
	ahead, behind, _, err := readAheadBehind(target)
	if err != nil {
		return err
	}

	hasChange := (mode == constants.HCModeAhead && ahead > 0) || (mode == constants.HCModeBehind && behind > 0)
	fmt.Println(boolStr(hasChange))

	return nil
}

// printHasChangeAll prints structured output for all three dimensions.
func printHasChangeAll(target string) error {
	isDirty, err := isWorkingTreeDirty(target)
	if err != nil {
		return err
	}

	return formatHasChangeAllUpstream(target, isDirty)
}

func formatHasChangeAllUpstream(target string, isDirty bool) error {
	ahead, behind, hasUpstream, err := readAheadBehind(target)
	if err != nil {
		return err
	}

	if !hasUpstream {
		fmt.Printf(constants.MsgHCAllNoUpstream, boolStr(isDirty))

		return nil
	}

	fmt.Printf(constants.MsgHCAllFmt, boolStr(isDirty), boolStr(ahead > 0), boolStr(behind > 0))

	return nil
}

// boolStr returns "true" or "false" for printing.
func boolStr(hasValue bool) string {
	if hasValue {
		return constants.HCTrue
	}

	return constants.HCFalse
}

// readAheadBehind runs `git rev-list --left-right --count HEAD...@{upstream}`
// in target and returns the (ahead, behind, hasUpstream, err) tuple. hasUpstream is false when there
// is no configured upstream.
func readAheadBehind(target string) (int, int, bool, error) {
	cmd := exec.Command(constants.GitBin, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = target
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, false, nil
	}

	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) != 2 {
		return 0, 0, false, apperror.WrapSimple(apperror.NewExecutionError(fmt.Sprintf("unexpected output: %s", string(out))), "parse rev-list count")
	}

	ahead, _ := strconv.Atoi(parts[0])
	behind, _ := strconv.Atoi(parts[1])

	return ahead, behind, true, nil
}

// fetchRemoteIn runs `git fetch` in target. Errors are printed but do not
// fail the command — the user gets a stale ahead/behind result with a warning.
func fetchRemoteIn(target string) error {
	cmd := exec.Command(constants.GitBin, "fetch")
	cmd.Dir = target
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnHCFetchFailed, target, err)

		return apperror.WrapSimple(err, "git fetch in "+target)
	}

	return nil
}
