package cmd

import "strings"

// knownValueFlags records CLI flags that consume the subsequent token as their value.
var knownValueFlags = map[string]bool{
	// release / commit-flow flags
	"--assets": true, "--commit": true, "--branch": true,
	"--bump": true, "--notes": true, "--targets": true,
	"--bundle": true, "--zip-group": true,
	"-N": true, "-Z": true,
	// self-install / self-uninstall value-taking flags
	"--dir": true, "--version": true,
	"--profile": true, "--shell-mode": true,
	// clone-next value-taking flags
	"--csv": true, "--ssh-key": true, "-K": true,
	"--target-dir": true,
	// commit-transfer value-taking flags (spec 106 §8).
	// Without these, "--drop ^WIP --no-provenance" would have
	// --drop swallow "--no-provenance" as its regex value and
	// the negation toggle would silently never fire.
	"--strip": true, "--drop": true,
	"--limit": true, "--since": true,
	// reclone / scan unified manifest flag. Without this entry,
	// `gitmap reclone --manifest path.json --execute` would treat
	// `path.json` as a positional <file>, triggering the
	// manifest-vs-positional conflict and exiting 2.
	"--manifest": true,
	// reclone --scan-root <dir>: redirects auto-pickup root.
	// Same reordering hazard as --manifest — without this entry
	// the directory would land in the positional slot.
	"--scan-root": true,
	// templates list filter flags: without these, `--kind ignore`
	// would split into `--kind` (parsed as a bare bool-style flag,
	// value left empty) and `ignore` (re-classified as positional),
	// which is why TestParseTemplatesListFlagsLowersValues failed.
	"--kind": true, "--lang": true, "--exclude": true, "--except": true,
}

// isKnownValueFlag reports whether the specified flag expects a value argument.
func isKnownValueFlag(flagName string) bool {
	return knownValueFlags[flagName]
}

// appendFlagWithValue appends a flag and its value (if applicable) to the flags slice.
func appendFlagWithValue(args []string, argIdx int, flags []string) ([]string, int) {
	currentArg := args[argIdx]
	flags = append(flags, currentArg)
	hasValue := isKnownValueFlag(currentArg)
	hasNext := argIdx+1 < len(args)
	if hasValue && hasNext {
		argIdx++
		flags = append(flags, args[argIdx])
	}

	return flags, argIdx
}

// partitionFlagsAndArgs separates flags from positional arguments in CLI args.
func partitionFlagsAndArgs(args []string) ([]string, []string) {
	var flags []string
	var positional []string
	for argIdx := 0; argIdx < len(args); argIdx++ {
		currentArg := args[argIdx]
		isFlag := strings.HasPrefix(currentArg, "-")
		if !isFlag {
			positional = append(positional, currentArg)
			continue
		}
		flags, argIdx = appendFlagWithValue(args, argIdx, flags)
	}

	return flags, positional
}

// reorderFlagsBeforeArgs moves flag-like arguments (starting with "-")
// before positional arguments. Go's flag package stops parsing at the
// first non-flag argument, so "gitmap release v2.55 -y" would silently
// ignore -y. This reorders to "-y v2.55" so all flags are parsed.
//
// Flags that take a value (e.g. --bump patch, -N "note") are kept
// together with their value argument.
func reorderFlagsBeforeArgs(args []string) []string {
	flags, positional := partitionFlagsAndArgs(args)

	return append(flags, positional...)
}

