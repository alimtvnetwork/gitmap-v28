package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runRescanSubtree is the CLI entry point for
// `gitmap rescan-subtree <absolutePath>`.
//
// Workflow this command supports:
//
//  1. User runs `gitmap scan` on a wide root and inspects the resulting
//     CSV/JSON. Some rows arrive at `depth == --max-depth` — those are
//     boundary discoveries whose nested repos were skipped because the
//     walker hit the cap (see helptext/scan.md, "Depth column and
//     --max-depth").
//  2. User copies the `absolutePath` from one of those at-cap rows.
//  3. `gitmap rescan-subtree <absolutePath>` re-runs `gitmap scan` against
//     that subtree with a deeper default `--max-depth`
//     (constants.RescanSubtreeDefaultMaxDepth) so the previously hidden
//     repos surface in a single command — no need to remember the
//     scan-flag incantation or the cap arithmetic.
//
// Any flag accepted by `gitmap scan` may be passed after the path; it
// is forwarded verbatim to runScan. If the user supplies their own
// `--max-depth`, this command does NOT override it.
//
// Exit codes:
//
//	0 — rescan completed (scan itself decides 0 vs. non-zero on hard errors)
//	2 — bad CLI usage: missing path, --help, or path that is not an
//	    existing directory. Distinct from the scan failure exit code so
//	    shell wrappers can tell "you invoked me wrong" apart from "the
//	    walk itself failed".
func runRescanSubtree(args []string) {
	checkHelp("rescan-subtree", args)

	path, rest, err := splitRescanSubtreeArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	abs, err := resolveRescanSubtreePath(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	scanArgs := buildRescanSubtreeArgs(abs, rest)
	fmt.Printf("  ▶ gitmap rescan-subtree — %s (max-depth=%s)\n",
		abs, extractMaxDepthForLog(scanArgs))
	runScan(scanArgs)
}

// splitRescanSubtreeArgs separates the required <path> positional from
// the trailing flags that will be forwarded to runScan. The path may
// appear anywhere in args (before or after flags) so users who muscle-
// memory `--quiet path` still get a clean parse.
//
// Returns an error message safe to print directly to stderr when the
// path is missing or the user passes more than one positional.
func splitRescanSubtreeArgs(args []string) (string, []string, error) {
	var positionals []string
	var flags []string
	skipNext := false
	for i, a := range args {
		if skipNext {
			skipNext = false
			flags = append(flags, a)
			continue
		}
		if len(a) > 0 && a[0] == '-' {
			flags = append(flags, a)
			// `--flag value` (no `=`) consumes the next token. A flag
			// passed as `--flag=value` or a bare boolean stays single.
			if i+1 < len(args) && !flagHasInlineValue(a) && !isLikelyBoolFlag(a) {
				skipNext = true
			}
			continue
		}
		positionals = append(positionals, a)
	}

	switch len(positionals) {
	case 0:
		return "", nil, fmt.Errorf(
			"  Error: rescan-subtree requires <absolutePath> (e.g. `gitmap rescan-subtree /home/me/work/monorepo`)")
	case 1:
		return positionals[0], flags, nil
	default:
		return "", nil, fmt.Errorf(
			"  Error: rescan-subtree accepts exactly one <absolutePath>, got %d: %v",
			len(positionals), positionals)
	}
}

// flagHasInlineValue reports whether a `-flag=value` token already
// carries its value inline.
func flagHasInlineValue(token string) bool {
	for i := 0; i < len(token); i++ {
		if token[i] == '=' {
			return true
		}
	}
	return false
}

// isLikelyBoolFlag returns true for flags `gitmap scan` defines as bare
// booleans, so we don't accidentally swallow the next positional as
// their "value". Kept narrow on purpose — a wrong guess here only
// affects flag-then-path ordering, never correctness of the scan.
func isLikelyBoolFlag(token string) bool {
	switch token {
	case "--quiet", "-quiet",
		"--github-desktop", "-github-desktop",
		"--open", "-open",
		"--no-vscode-sync", "-no-vscode-sync",
		"--no-auto-tags", "-no-auto-tags",
		"--report-errors", "-report-errors",
		"--no-probe", "-no-probe",
		"--no-probe-wait", "-no-probe-wait":
		return true
	}
	return false
}

// resolveRescanSubtreePath turns the user's path argument into an
// absolute path and verifies it exists and is a directory. Returns a
// stderr-ready error when the directory is missing, is a file, or is
// otherwise unstattable.
func resolveRescanSubtreePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("  Error: cannot resolve %q: %w", path, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf(
				"  Error: rescan-subtree target does not exist: %s\n"+
					"         Did you copy the absolutePath from a row that has since moved?",
				abs)
		}
		return "", fmt.Errorf("  Error: cannot stat %s: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf(
			"  Error: rescan-subtree target is not a directory: %s", abs)
	}
	return abs, nil
}

// buildRescanSubtreeArgs assembles the argv for runScan: forwarded
// flags first (so reorderFlagsBeforeArgs treats them normally), then
// the positional directory, then a synthetic --max-depth iff the user
// did not pass their own. The synthetic value uses
// constants.RescanSubtreeDefaultMaxDepth so the typical "I hit the
// default cap of 4 — go deeper" workflow finishes in one command.
func buildRescanSubtreeArgs(absDir string, forwardedFlags []string) []string {
	out := make([]string, 0, len(forwardedFlags)+3)
	out = append(out, forwardedFlags...)
	if !containsMaxDepthFlag(forwardedFlags) {
		out = append(out,
			"--"+constants.FlagScanMaxDepth,
			strconv.Itoa(constants.RescanSubtreeDefaultMaxDepth))
	}
	out = append(out, absDir)
	return out
}

// containsMaxDepthFlag reports whether the forwarded flag slice already
// includes `--max-depth` (in either `--flag value` or `--flag=value`
// form, with one or two leading dashes).
func containsMaxDepthFlag(flags []string) bool {
	want := constants.FlagScanMaxDepth
	for _, f := range flags {
		if f == "--"+want || f == "-"+want {
			return true
		}
		// `--max-depth=8` / `-max-depth=8`
		prefixDouble := "--" + want + "="
		prefixSingle := "-" + want + "="
		if startsWith(f, prefixDouble) || startsWith(f, prefixSingle) {
			return true
		}
	}
	return false
}

// startsWith is a tiny strings.HasPrefix replacement to keep this file
// dependency-light (the cmd package already pulls strings transitively
// elsewhere; using it here would still be fine, this is just clearer).
func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

// extractMaxDepthForLog returns the --max-depth value from the final
// scan args, purely for the human-facing banner. Defaults to "auto" so
// the banner still reads cleanly if the user supplied a non-trivial
// flag form we don't bother decoding.
func extractMaxDepthForLog(scanArgs []string) string {
	want := constants.FlagScanMaxDepth
	prefixDouble := "--" + want + "="
	prefixSingle := "-" + want + "="
	for i, a := range scanArgs {
		// Space form: `--max-depth N` / `-max-depth N`. The bare token
		// must end the inspection — falling through to the inline-form
		// check would mis-classify a malformed `--max-depth` (no value)
		// as a non-match and let the loop wander into unrelated flags.
		if a == "--"+want || a == "-"+want {
			if i+1 < len(scanArgs) {
				return scanArgs[i+1]
			}
			return "auto"
		}
		if startsWith(a, prefixDouble) {
			return a[len(prefixDouble):]
		}
		if startsWith(a, prefixSingle) {
			return a[len(prefixSingle):]
		}
	}
	return "auto"
}
