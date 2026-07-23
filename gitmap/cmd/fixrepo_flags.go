package cmd

// Flag parsing for `gitmap fix-repo`. Accepts both GNU long-form
// (`--dry-run`) and PowerShell single-dash forms (`-DryRun`) so the
// command is a drop-in replacement for fix-repo.ps1.

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// isFixRepoSpanArg reports whether a is a span-style argument: a bare
// positive integer like "4" or a dash-prefixed integer like "-4" /
// "-7". v5.45.0+ accepts these as a generalization of the canonical
// -2/-3/-5 modes so `gitmap fix-repo 4` and `gitmap fix-repo -4` are
// no longer rejected with E_BAD_FLAG. The span value is the integer
// itself; computeFixRepoSpan handles it as a raw "-N" mode token.
func isFixRepoSpanArg(a string) bool {
	s := a
	if strings.HasPrefix(s, "-") {
		s = s[1:]
	}
	if s == "" {
		return false
	}
	n, err := strconv.Atoi(s)
	return err == nil && n > 0
}

// normalizeFixRepoSpanArg coerces a bare digit ("4") into the
// canonical dash form ("-4") so downstream span computation has a
// single shape to switch on.
func normalizeFixRepoSpanArg(a string) string {
	if strings.HasPrefix(a, "-") {
		return a
	}
	return "-" + a
}

// fixRepoModeFlags lists every accepted mode-token in lower case.
// `--all` is also accepted via the long-flag branch below.
var fixRepoModeFlags = []string{
	constants.FixRepoModeFlag2,
	constants.FixRepoModeFlag3,
	constants.FixRepoModeFlag5,
}

// parseFixRepoArgs walks args once, collecting modes, booleans, and
// the optional --config value. Returns a typed error so the caller
// can wrap it in the standard E_BAD_FLAG message.
func parseFixRepoArgs(args []string) (fixRepoOptions, error) {
	out := fixRepoOptions{mode: constants.FixRepoModeFlag2}
	modes := []string{}
	unknown := []string{}
	i := 0
	for i < len(args) {
		consumed, err := consumeOneFixRepoArg(args, i, &out, &modes, &unknown)
		if err != nil {
			return out, err
		}
		i += consumed
	}

	return finalizeFixRepoOpts(out, modes, unknown)
}

// consumeOneFixRepoArg processes args[i] and returns how many tokens
// were consumed (1 or 2). Split out so parseFixRepoArgs stays inside
// the function-length budget.
func consumeOneFixRepoArg(args []string, i int, out *fixRepoOptions,
	modes, unknown *[]string,
) (int, error) {
	a := args[i]
	if isFixRepoModeArg(a) {
		*modes = append(*modes, normalizeFixRepoMode(a))

		return 1, nil
	}
	if isFixRepoSpanArg(a) {
		*modes = append(*modes, normalizeFixRepoSpanArg(a))

		return 1, nil
	}
	if isFixRepoDryRunArg(a) {
		out.isDryRun = true

		return 1, nil
	}
	if isFixRepoVerboseArg(a) {
		out.isVerbose = true

		return 1, nil
	}
	if isFixRepoStrictArg(a) {
		out.isStrict = true

		return 1, nil
	}
	consumed, ok, err := consumeFixRepoRestrictArg(args, i, out)
	if err != nil {
		return 0, err
	}
	if ok {
		return consumed, nil
	}
	consumed, ok, err = consumeFixRepoConfigArg(args, i, out)
	if err != nil {
		return 0, err
	}
	if ok {
		return consumed, nil
	}
	*unknown = append(*unknown, a)

	return 1, nil
}

// finalizeFixRepoOpts validates the collected mode/unknown lists and
// returns the canonical fixRepoOptions.
func finalizeFixRepoOpts(out fixRepoOptions, modes, unknown []string) (fixRepoOptions, error) {
	if len(modes) > 1 {
		return out, fmt.Errorf("multiple mode flags: %s", strings.Join(modes, " "))
	}
	if len(unknown) > 0 {
		return out, fmt.Errorf("unknown flag(s): %s\n%s", strings.Join(unknown, " "), fixRepoFlagHint())
	}
	if len(modes) == 1 {
		out.mode = modes[0]
	}

	return out, nil
}

// isFixRepoModeArg reports whether a is one of -2 / -3 / -5 / -All / --all.
func isFixRepoModeArg(a string) bool {
	for _, m := range fixRepoModeFlags {
		if a == m {
			return true
		}
	}
	low := strings.ToLower(a)

	return low == "-all" || low == "--all"
}

// normalizeFixRepoMode collapses -All / --all variants to "--all"
// and leaves numeric modes (-2/-3/-5) untouched.
func normalizeFixRepoMode(a string) string {
	low := strings.ToLower(a)
	if low == "-all" || low == "--all" {
		return "--" + constants.FixRepoFlagAll
	}

	return a
}

// isFixRepoDryRunArg matches both --dry-run and -DryRun (any case).
func isFixRepoDryRunArg(a string) bool {
	low := strings.ToLower(a)

	return low == "--"+constants.FixRepoFlagDryRun || low == "-"+constants.FixRepoFlagDryRun
}

// isFixRepoVerboseArg matches both --verbose and -Verbose (any case).
func isFixRepoVerboseArg(a string) bool {
	low := strings.ToLower(a)

	return low == "--"+constants.FixRepoFlagVerbose || low == "-"+constants.FixRepoFlagVerbose
}

// isFixRepoStrictArg matches both --strict and -Strict (any case).
// The flag has no value — presence flips opts.isStrict.
func isFixRepoStrictArg(a string) bool {
	low := strings.ToLower(a)

	return low == "--"+constants.FixRepoFlagStrict || low == "-"+constants.FixRepoFlagStrict
}

// consumeFixRepoConfigArg handles `--config <p>`, `-Config <p>`,
// `--config=<p>`, and `-Config=<p>`. Returns (consumed, matched, err).
func consumeFixRepoConfigArg(args []string, i int, out *fixRepoOptions) (int, bool, error) {
	a := args[i]
	low := strings.ToLower(a)
	bareLong := "--" + constants.FixRepoFlagConfig
	bareShort := "-" + constants.FixRepoFlagConfig
	if low == bareLong || low == bareShort {
		if i+1 >= len(args) {
			return 0, true, errors.New("--config requires a path")
		}
		out.configPath = args[i+1]

		return 2, true, nil
	}
	prefixLong := bareLong + "="
	prefixShort := bareShort + "="
	if strings.HasPrefix(low, prefixLong) {
		out.configPath = a[len(prefixLong):]

		return 1, true, nil
	}
	if strings.HasPrefix(low, prefixShort) {
		out.configPath = a[len(prefixShort):]

		return 1, true, nil
	}

	return 0, false, nil
}

// consumeFixRepoRestrictArg handles `--restrict <val>` / `-restrict <val>`
// / `-r <val>` plus `=<val>` forms. Accepts `no-version` or its short
// alias `nv` (case-insensitive). Returns (consumed, matched, err).
func consumeFixRepoRestrictArg(args []string, i int, out *fixRepoOptions) (int, bool, error) {
	a := args[i]
	low := strings.ToLower(a)
	bareLong := "--" + constants.FixRepoFlagRestrict
	bareShort := "-" + constants.FixRepoFlagRestrict
	bareTiny := "-" + constants.FixRepoFlagRestrictShort
	if low == bareLong || low == bareShort || low == bareTiny {
		if i+1 >= len(args) {
			return 0, true, fmt.Errorf("%s requires a value (no-version|nv)", a)
		}

		return 2, true, applyRestrictValue(out, args[i+1])
	}
	for _, p := range []string{bareLong + "=", bareShort + "=", bareTiny + "="} {
		if strings.HasPrefix(low, p) {
			return 1, true, applyRestrictValue(out, a[len(p):])
		}
	}

	return 0, false, nil
}

// applyRestrictValue normalizes the restrict value and sets the
// matching option, or returns an error for unknown values.
func applyRestrictValue(out *fixRepoOptions, v string) error {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case constants.FixRepoRestrictNoVersion, constants.FixRepoRestrictNoVersionShort:
		out.restrictNoVersion = true

		return nil
	}

	return fmt.Errorf("unknown --restrict value %q (want: no-version|nv)", v)
}
