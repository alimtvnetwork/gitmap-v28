package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// replaceOpts holds parsed flags for `gitmap replace`. Kept tiny so
// every handler can pass it by value.
type replaceOpts struct {
	yes        bool
	dryRun     bool
	quiet      bool
	audit      bool
	exts       []string // pre-normalized extensions (with leading dot)
	extCaseIns bool     // true = lowercase comparison, false = byte-exact
}

// parseReplaceFlags consumes flag tokens (in any position) and returns
// the remaining positional arguments. We hand-roll the audit detection
// because `--audit` is documented as a subcommand, not a typed flag,
// and we want the same token to work whether the user writes it before
// or after positional args.
func parseReplaceFlags(args []string) (replaceOpts, []string, error) {
	opts, rest := stripAuditToken(args)

	fs, raw := defineReplaceFlags()
	flags, positional := splitReplaceFlagsAndArgs(rest)
	if err := fs.Parse(flags); err != nil {
		return opts, nil, errors.New("replace: " + err.Error())
	}

	opts.yes = *raw.yes || *raw.yesShort
	opts.dryRun = *raw.dry
	opts.quiet = *raw.quiet || *raw.quietShort
	opts.extCaseIns = resolveExtCase(*raw.extCase)
	opts.exts = normalizeExtList(*raw.ext, opts.extCaseIns)

	return opts, positional, nil
}

// rawReplaceFlags holds pointers returned from FlagSet.* registrations
// so parseReplaceFlags can stay under the 15-line ceiling.
type rawReplaceFlags struct {
	yes, yesShort, dry, quiet, quietShort *bool
	ext, extCase                          *string
}

// defineReplaceFlags registers every flag on a fresh FlagSet.
func defineReplaceFlags() (*flag.FlagSet, rawReplaceFlags) {
	fs := flag.NewFlagSet(constants.CmdReplace, flag.ContinueOnError)
	r := rawReplaceFlags{
		yes:        fs.Bool(constants.ReplaceFlagYes, false, "Skip the confirmation prompt"),
		yesShort:   fs.Bool(constants.ReplaceFlagYesS, false, "Alias for --yes"),
		dry:        fs.Bool(constants.ReplaceFlagDryRun, false, constants.FlagDescDryRun),
		quiet:      fs.Bool(constants.ReplaceFlagQuiet, false, "Suppress per-file diff lines"),
		quietShort: fs.Bool(constants.ReplaceFlagQuietS, false, "Alias for --quiet"),
		ext:        fs.String(constants.ReplaceFlagExt, "", constants.FlagDescReplaceExt),
		extCase:    fs.String(constants.ReplaceFlagExtCase, "", constants.FlagDescReplaceExtCase),
	}
	return fs, r
}

// resolveExtCase maps the raw --ext-case string to the boolean the
// walker uses. Empty (flag omitted) keeps the default behavior. An
// unknown value is fatal — silent fallback would mask user typos.
func resolveExtCase(raw string) bool {
	v := strings.ToLower(strings.TrimSpace(raw))
	// ReplaceExtCaseDefault aliases ReplaceExtCaseInsensitive, so the
	// empty string and both names share one branch (collapsing them
	// avoids a duplicate-case compile error).
	switch v {
	case "", constants.ReplaceExtCaseInsensitive:
		return true
	case constants.ReplaceExtCaseSensitive:
		return false
	default:
		fmt.Fprintf(os.Stderr, constants.ErrReplaceBadExtCase,
			constants.ReplaceExtCaseSensitive,
			constants.ReplaceExtCaseInsensitive, raw)
		os.Exit(1)
		return true
	}
}

// stripAuditToken pulls --audit out of args before flag.Parse sees it.
func stripAuditToken(args []string) (replaceOpts, []string) {
	out := make([]string, 0, len(args))
	var opts replaceOpts
	for _, a := range args {
		if a == constants.ReplaceSubcmdAudit {
			opts.audit = true
			continue
		}
		out = append(out, a)
	}
	return opts, out
}

// splitReplaceFlagsAndArgs partitions args so flag.Parse only sees
// flag-like tokens. Anything else (including `-N` and `all`) is a
// positional. We treat tokens beginning with `--` or `-X`/`-XX`
// (letters) as flags. A flag's value (`--ext .go,.md`) sticks to its
// flag token via the previous-token lookback.
func splitReplaceFlagsAndArgs(args []string) (flags, positional []string) {
	expectValue := false
	for _, a := range args {
		if expectValue {
			flags = append(flags, a)
			expectValue = false
			continue
		}
		if isReplaceFlag(a) {
			flags = append(flags, a)
			expectValue = needsValue(a)
			continue
		}
		positional = append(positional, a)
	}
	return
}

// needsValue reports whether a flag token consumes the next arg as its
// value. Boolean replace flags do not; --ext and --ext-case do. We
// don't handle `--ext=x` because flag.Parse splits that itself.
func needsValue(token string) bool {
	if strings.Contains(token, "=") {
		return false
	}
	name := strings.TrimLeft(token, "-")
	return name == constants.ReplaceFlagExt || name == constants.ReplaceFlagExtCase
}

// isReplaceFlag returns true for `--xxx`, `-y`, `-q`. It deliberately
// returns false for `-N` digit forms so they remain positional.
func isReplaceFlag(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	if s[1] == '-' {
		return true
	}
	c := s[1]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// normalizeExtList parses --ext into a deduplicated list of extensions
// with leading dots. When caseInsensitive is true the entries are
// lowercased so matchesExtFilter can do a fast direct compare; when
// false the user's casing is preserved byte-for-byte.
func normalizeExtList(raw string, caseInsensitive bool) []string {
	if raw == "" {
		return nil
	}
	seen := make(map[string]struct{}, 4)
	out := make([]string, 0, 4)
	for _, piece := range strings.Split(raw, constants.ReplaceExtSep) {
		ext := normalizeOneExt(piece, caseInsensitive)
		if ext == "" {
			continue
		}
		if _, dup := seen[ext]; dup {
			continue
		}
		seen[ext] = struct{}{}
		out = append(out, ext)
	}
	return out
}

// normalizeOneExt trims spaces, optionally lowercases, and ensures a
// leading dot. Returns "" for inputs that should be dropped (empty
// piece, or the lone string ".").
func normalizeOneExt(piece string, caseInsensitive bool) string {
	piece = strings.TrimSpace(piece)
	if caseInsensitive {
		piece = strings.ToLower(piece)
	}
	if piece == "" || piece == "." {
		return ""
	}
	if piece[0] != '.' {
		piece = "." + piece
	}
	return piece
}
