package cmdpull

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
)

func checkEfficientSubcommand(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	if isEfficientTableSubcmd(args[0]) {
		return true, RunPullAllEfficient(args[1:], true, "pull "+args[0], isShortEfficientSubcmd(args[0]))
	}
	if isEfficientPullSubcmd(args[0]) {
		return true, RunPullAllEfficient(args[1:], false, "pull "+args[0], isShortEfficientSubcmd(args[0]))
	}

	return false, nil
}

func isEfficientTableSubcmd(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "all-efficient-table" || lower == "paet" || lower == "aet"
}

func isEfficientPullSubcmd(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "all-efficient" || lower == "ae" || lower == "pae" || lower == "pull-ae" || lower == "efficient"
}

func isShortEfficientSubcmd(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "ae" || lower == "pae" || lower == "pull-ae" || lower == "paet" || lower == "aet"
}

func isPullAllRootCmd() bool {
	if len(os.Args) <= 1 {
		return false
	}
	first := strings.ToLower(os.Args[1])

	return isPullAllToken(first)
}

func isPullAllTableRootCmd() bool {
	if len(os.Args) <= 1 {
		return false
	}
	first := strings.ToLower(os.Args[1])

	return first == "pull-all-table" || first == "pat"
}

func hasPullAllArg(args []string) bool {
	for _, a := range args {
		if isPullAllToken(strings.ToLower(a)) {
			return true
		}
	}

	return false
}

func isPullAllToken(token string) bool {
	return token == "--all" || token == "-all" || token == "-a" || token == "all" || token == "pa" || token == "ta" || token == "pull-all" || token == "pat" || token == "pull-all-table" || token == "pas" || token == "pull-all-ssh"
}

func printPullInvocationHeader(isPullAll bool) {
	cwd, _ := os.Getwd()
	cmdName := "pull"
	if isPullAll {
		cmdName = "pull-all"
	}
	fmt.Printf("\n  %s→%s %sgitmap %s%s %s(cwd: %s)%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, cmdName, constants.ColorReset,
		constants.ColorDim, cwd, constants.ColorReset)
}

func resolveSubArrow() string {
	if glyphs.Resolve() == glyphs.ModeSafe {
		return "->"
	}

	return "→"
}

func dispatchPullExecution(opts pullOptions) error {
	if isPullCWDEnabled(opts) {
		return executePullCWDWithNotice(opts.isRaw)
	}
	if ShouldFallbackToPullAll(opts) {
		announceNonGitFallbackUnlessJSON(opts.isJSON)
		opts.all = true
	}

	return runPullBatch(opts)
}

func executePullCWDWithNotice(isRaw bool) error {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s %scwd is a git repo — running plain `git pull` here%s\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)

	return runPullCWD(isRaw)
}

func announceNonGitFallbackUnlessJSON(isJSON bool) {
	if !isJSON {
		AnnounceNonGitPullFallback()
	}
}

func handleEmptyBatchRecords(isJSON bool) error {
	if isJSON {
		return renderPullBatchJSONSummary(0, nil, 0)
	}

	return nil
}
