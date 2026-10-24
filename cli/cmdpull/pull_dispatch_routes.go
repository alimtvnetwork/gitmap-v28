package cmdpull

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
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

func printPullInvocationHeader(isPullAll bool, args []string) {
	header := resolvePullInvocationHeader(isPullAll, args)
	arrow := resolveSubArrow()
	fmt.Printf("\n  %s%s%s %s\n",
		constants.ColorCyan, arrow, constants.ColorReset, header)
}

func resolvePullInvocationHeader(isPullAll bool, args []string) string {
	if isPullAll || hasPullAllArg(args) {
		return formatPullAllHeader()
	}
	targetArg := findFirstPositionalArg(args)
	if targetArg != "" {
		return formatRepoPathHeader(resolveTargetRepoPath(targetArg))
	}
	if isGitRepoCWD() {
		return formatRepoPathHeader(resolveCWDRepoPath())
	}

	return formatPullAllHeader()
}

func formatPullAllHeader() string {
	return fmt.Sprintf("%spull-all%s %s(all workspaces)%s",
		constants.ColorBold, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
}

func formatRepoPathHeader(repoPath string) string {
	return fmt.Sprintf("%s%s%s", constants.ColorBold, repoPath, constants.ColorReset)
}

func resolveCWDRepoPath() string {
	top := getGitRepoCWDTopLevel()
	if top != "" {
		return top
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}

	return filepath.Clean(cwd)
}

func resolveTargetRepoPath(target string) string {
	abs, err := filepath.Abs(target)
	if err == nil && !cloner.IsMissingRepo(abs) {
		return filepath.Clean(abs)
	}
	records := resolveSlugTarget(target)
	if len(records) > 0 && records[0].AbsolutePath != "" {
		return filepath.Clean(records[0].AbsolutePath)
	}

	return target
}

func findFirstPositionalArg(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && !isPullAllToken(strings.ToLower(a)) {
			return a
		}
	}

	return ""
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
