package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
)

// Flag tokens for pretty-mode parsing. Two synonym pairs are accepted so
// the flag is comfortable to type either way: `--pretty[=bool]` follows
// the standard Go flag style, while `--no-pretty` matches the convention
// many CLIs use for explicit negation.
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
// (it also mirrors the widely-supported NO_COLOR env convention).
var prettyFlagPrefixes = []string{
	flagPrettyPositive, flagPrettyNegative,
	flagColorPositive, flagColorNegative,
}

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
func ParsePrettyFlag(args []string) ([]string, render.PrettyMode) {
	mode := render.PrettyAuto
	out := make([]string, 0, len(args))
	for _, a := range args {
		token, value, hasValue := splitPrettyToken(a)
		switch token {
		case flagPrettyPositive, flagColorPositive:
			mode = resolvePositivePretty(value, hasValue, mode, &out, a)
		case flagPrettyNegative, flagColorNegative:
			mode = render.PrettyOff
		default:
			out = append(out, a)
		}
	}

	return out, mode
}

// splitPrettyToken splits "--pretty=value" into ("--pretty", "value", true)
// and "--pretty" into ("--pretty", "", false). Anything else returns the
// original token in slot 0 with hasValue=false so the caller can passthrough.
// Recognizes every prefix in prettyFlagPrefixes (pretty + color synonyms).
func splitPrettyToken(arg string) (token, value string, hasValue bool) {
	if !hasPrettyPrefix(arg) {
		return arg, "", false
	}
	if eq := strings.IndexByte(arg, '='); eq >= 0 {
		return arg[:eq], arg[eq+1:], true
	}

	return arg, "", false
}

// hasPrettyPrefix reports whether arg begins with any token managed by
// ParsePrettyFlag. Extracted so the prefix list lives in one place.
func hasPrettyPrefix(arg string) bool {
	for _, p := range prettyFlagPrefixes {
		if strings.HasPrefix(arg, p) {
			return true
		}
	}

	return false
}

// resolvePositivePretty maps a "--pretty[=value]" occurrence to a
// PrettyMode. Falls back to keeping the original token in `out` when
// the value is unrecognized so flag.Parse downstream can report it.
func resolvePositivePretty(value string, hasValue bool, current render.PrettyMode, out *[]string, original string) render.PrettyMode {
	if !hasValue {
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
