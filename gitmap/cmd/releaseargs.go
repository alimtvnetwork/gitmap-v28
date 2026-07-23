package cmd

import "strings"

// reorderFlagsBeforeArgs moves flag-like arguments (starting with "-")
// before positional arguments. Go's flag package stops parsing at the
// first non-flag argument, so "gitmap release v2.55 -y" would silently
// ignore -y. This reorders to "-y v2.55" so all flags are parsed.
//
// Flags that take a value (e.g. --bump patch, -N "note") are kept
// together with their value argument.
func reorderFlagsBeforeArgs(args []string) []string {
	var flags []string
	var positional []string

	// Known flags that consume the next argument as a value.
	valueFlags := map[string]bool{
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
		"--kind": true, "--lang": true,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			// If this flag takes a value, grab the next arg too.
			if valueFlags[arg] && i+1 < len(args) {
				i++
				flags = append(flags, args[i])
			}
		} else {
			positional = append(positional, arg)
		}
	}

	return append(flags, positional...)
}
