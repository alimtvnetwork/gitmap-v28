package cmd

// CLI entry point for `gitmap fix-repo` (alias `fr`). This is the
// Go-native re-implementation of fix-repo.ps1 / fix-repo.sh. Behavior,
// exit codes, and config schema match the scripts 1:1. Spec:
// spec/04-generic-cli/27-fix-repo-command.md.

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// fixRepoOptions captures parsed CLI inputs for one invocation.
type fixRepoOptions struct {
	mode       string // "-2" | "-3" | "-5" | "--all"
	span       int
	isDryRun   bool
	isVerbose  bool
	isStrict   bool // --strict / -Strict: post-rewrite `go test` on touched pkgs
	configPath string
	// restrictNoVersion (v5.39.0+): when true, suppress the v1→v2
	// bare-base sweep so ONLY `{base}-vN` tokens are rewritten. Set
	// via `--restrict no-version` / `-r nv`. See spec
	// 27-fix-repo-command.md §"Restrict modes".
	restrictNoVersion bool
}

// runFixRepo is the dispatcher entry. checkHelp first so `--help`
// works even when other args would fail to parse.
func runFixRepo(args []string) {
	checkHelp(constants.CmdFixRepo, args)
	opts, err := parseFixRepoArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrBadFlagFmt, err.Error())
		cliexit.Exit(constants.FixRepoExitBadFlag)
	}
	identity := resolveFixRepoIdentity()
	loadFixRepoConfig(opts.configPath, identity.root)
	opts.span = computeFixRepoSpan(opts.mode, identity.current)
	targets := computeFixRepoTargets(identity.current, opts.span)
	emitFixRepoHeader(identity, opts.mode, targets)
	if len(targets) == 0 {
		emitFixRepoSummary(0, 0, 0, opts.isDryRun)
		fmt.Print(constants.FixRepoMsgNothing)
		emitFixRepoTips(opts, 0)
		cliexit.Exit(constants.FixRepoExitOk)
	}
	result := runFixRepoSweep(identity, targets, opts)
	if result.backup != nil {
		result.backup.Finalize()
	}
	emitFixRepoSummary(result.scanned, result.changed, result.replacements, opts.isDryRun)
	if !runFixRepoGofmt(result.goFiles, opts) {
		result.failed = true
	}
	if !runFixRepoStrict(identity.root, result.goFiles, opts) {
		// Tests-failed is a distinct exit code so CI scripts can
		// branch on "rewrite produced semantically broken code" vs
		// other write/IO failures. Reported even if gofmt also failed
		// — strict failure is the more actionable diagnosis.
		emitFixRepoTips(opts, result.changed)
		cliexit.Exit(constants.FixRepoExitTestsFailed)
	}
	emitFixRepoTips(opts, result.changed)
	if result.failed {
		cliexit.Exit(constants.FixRepoExitWriteFailed)
	}
	cliexit.Exit(constants.FixRepoExitOk)
}

// computeFixRepoSpan maps the mode flag to an integer span. `--all`
// expands to current-1 so every prior version is rewritten. Any
// `-N` mode (v5.45.0+) is parsed as the integer N so users can pass
// `gitmap fix-repo 4` / `-7` etc. without hitting E_BAD_FLAG.
func computeFixRepoSpan(mode string, current int) int {
	switch mode {
	case constants.FixRepoModeFlag2:
		return 2
	case constants.FixRepoModeFlag3:
		return 3
	case constants.FixRepoModeFlag5:
		return 5
	case "--" + constants.FixRepoFlagAll:
		return current - 1
	}
	if strings.HasPrefix(mode, "-") {
		if n, err := strconv.Atoi(mode[1:]); err == nil && n > 0 {
			return n
		}
	}

	return constants.FixRepoDefaultSpan
}

// computeFixRepoTargets returns the closed range [max(1, current-span)
// .. current-1]. Empty when current ≤ 1 or span ≤ 0.
func computeFixRepoTargets(current, span int) []int {
	if span <= 0 || current <= 1 {
		return nil
	}
	start := current - span
	if start < 1 {
		start = 1
	}
	end := current - 1
	if start > end {
		return nil
	}
	out := make([]int, 0, end-start+1)
	for n := start; n <= end; n++ {
		out = append(out, n)
	}

	return out
}

// emitFixRepoHeader prints the three-line header that matches the
// PowerShell script verbatim. Used by tests + parity scripts.
func emitFixRepoHeader(identity fixRepoIdentity, mode string, targets []int) {
	fmt.Printf(constants.FixRepoMsgHeaderFmt, identity.base, identity.current, mode)
	fmt.Printf(constants.FixRepoMsgTargetsFmt, formatFixRepoTargets(targets))
	fmt.Printf(constants.FixRepoMsgIdentityFmt, identity.host, identity.owner)
	fmt.Println()
}

// emitFixRepoSummary prints the trailing summary block. Mirrors the
// PowerShell script's Write-Summary helper exactly.
func emitFixRepoSummary(scanned, changed, replacements int, isDryRun bool) {
	mode := constants.FixRepoModeWrite
	if isDryRun {
		mode = constants.FixRepoModeDryRun
	}
	fmt.Println()
	fmt.Printf(constants.FixRepoMsgScannedFmt, scanned)
	fmt.Printf(constants.FixRepoMsgChangedFmt, changed, replacements)
	fmt.Printf(constants.FixRepoMsgModeFmt, mode)
}

// formatFixRepoTargets renders a target list as `v1, v2, v3` or
// `(none)` when empty. Kept tiny so emitFixRepoHeader stays short.
func formatFixRepoTargets(targets []int) string {
	if len(targets) == 0 {
		return constants.FixRepoTargetsNone
	}
	out := ""
	for i, n := range targets {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("v%d", n)
	}

	return out
}
