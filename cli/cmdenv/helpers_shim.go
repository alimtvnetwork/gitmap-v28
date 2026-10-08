package cmdenv

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
	"strings"
)

// ParsePrettyFlag pulls --pretty / --no-pretty (and the --color /
// --no-color synonyms) out of args and returns the cleaned slice + the
// resolved render.PrettyMode. Accepted forms:
//
//	--pretty | --color                 → PrettyOn
//	--pretty=true|on|1|yes|y           → PrettyOn
//	--color=true|on|1|yes|y            → PrettyOn
//	--pretty=false|off|0|no|n          → PrettyOff
//	--color=false|off|0|no|n           → PrettyOff
//	--pretty=auto | --color=auto       → PrettyAuto (explicit reset)
//	--no-pretty | --no-color           → PrettyOff
//
// When the same flag is repeated, the **last** occurrence wins (matches
// stdlib flag.Parse semantics) — and "same flag" spans the synonym
// pair, so `--pretty --no-color` resolves to PrettyOff. When neither
// appears, the returned mode is PrettyAuto so callers can rely on
// Decide()'s default ladder.
//
// Unrecognized values fall through to PrettyAuto and the token is left
// in place so the downstream parser can produce a meaningful error.
func ParsePrettyFlag(args []string) ([]string, render.PrettyModeType) {
	mode := render.PrettyAuto
	out := make([]string, 0, len(args))
	for _, arg := range args {
		token, value, hasValue := splitPrettyToken(arg)
		mode = applyPrettyToken(token, value, hasValue, mode, &out, arg)
	}

	return out, mode
}

func printUsageFooterShort() {
	printGitmapIdentityBlockShort()
}

func splitPrettyToken(arg string) (token, value string, hasValue bool) {
	isMissingPrefix := !hasPrettyPrefix(arg)
	if isMissingPrefix {
		return arg, "", false
	}

	eq := strings.IndexByte(arg, '=')
	hasEqual := eq >= 0
	if hasEqual {
		return arg[:eq], arg[eq+1:], true
	}

	return arg, "", false
}
func applyPrettyToken(
	token,
	value string,
	hasValue bool,
	mode render.PrettyModeType,
	out *[]string,
	arg string,
) render.PrettyModeType {
	switch token {
	case flagPrettyPositive, flagColorPositive:
		return resolvePositivePretty(value, hasValue, mode, out, arg)
	case flagPrettyNegative, flagColorNegative:
		return render.PrettyOff
	default:
		*out = append(*out, arg)

		return mode
	}
}

const (
	flagPrettyPositive = "--pretty"
	flagPrettyNegative = "--no-pretty"
	flagColorPositive  = "--color"
	flagColorNegative  = "--no-color"
)

// prettyFlagPrefixes lists every token recognized by ParsePrettyFlag.
// Centralized so splitPrettyToken's prefix gate stays in sync with the
// switch in ParsePrettyFlag. `--color` / `--no-color` are accepted as
// synonyms for `--pretty` / `--no-pretty` because in this CLI the
// pretty-markdown pipeline is the only thing that emits ANSI color,
// and `--no-color` is the conventional spelling users reach for first
// (it also mirrors the widely-supported NO_COLOR env convention)
var prettyFlagPrefixes = []string{
	flagPrettyPositive, flagPrettyNegative,
	flagColorPositive, flagColorNegative,
}

// ParsePrettyFlag pulls --pretty / --no-pretty (and the --color /
// --no-color synonyms)
func hasPrettyPrefix(arg string) bool {
	for _, prefix := range prettyFlagPrefixes {
		hasMatch := strings.HasPrefix(arg, prefix)
		if hasMatch {
			return true
		}
	}

	return false
}
func resolvePositivePretty(
	value string,
	hasValue bool,
	current render.PrettyModeType,
	out *[]string,
	original string,
) render.PrettyModeType {
	isMissingValue := !hasValue
	if isMissingValue {
		return render.PrettyOn
	}

	switch strings.ToLower(value) {
	case "1", "t", "true", "on", "yes", "y":
		return render.PrettyOn
	case "0", "f", "false", "off", "no", "n":
		return render.PrettyOff
	case "auto", "":
		return render.PrettyAuto
	}

	*out = append(*out, original)

	return current
}
