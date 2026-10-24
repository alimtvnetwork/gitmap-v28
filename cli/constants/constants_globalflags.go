package constants

// Cross-command CLI flag values.
//
// Each command parses its own `flag.NewFlagSet`, so these strings are
// reused safely (e.g. --json appears for `list-versions`, `list-releases`,
// `amend-list`, etc.). Centralized here so renames stay consistent.
const (
	FlagOpenValue   = "--open"
	FlagJSON        = "--json"
	FlagLimit       = "--limit"
	FlagSource      = "--source"
	FlagAllRepos    = "--all-repos" // v3.20.0: list-releases multi-repo batch view
	FlagCompact     = "--compact"
	FlagGroups      = "--groups"
	FlagFilter      = "--filter" // v5.42.0: `gitmap help --filter <q>` search
	FlagFilterShort = "-f"
	FlagHelp        = "--help"
	FlagHelpShort   = "-h"
	FlagHelpWord    = "help"
	FlagForce       = "--force"
	FlagForceShort  = "-f"
	FlagYes         = "--yes"
	FlagYesShort    = "-y"
	FlagDryRun      = "--dry-run"
	FlagDryRunShort = "-d"
	FlagDryRunAlt   = "-n"
	FlagVerify      = "--verify"
	FlagPlatform    = "--platform"
	FlagTimeout     = "--timeout"
	FlagVerboseShort = "-v"
	FlagNoPush      = "--no-push"
	FlagPassword    = "--pass"
	FlagPasswordShort = "-p"
	FlagPort        = "--port"
)

// IsHelpFlag reports whether s matches any standard help flag variant.
func IsHelpFlag(s string) bool {
	return s == FlagHelp || s == FlagHelpShort || s == FlagHelpWord
}

// IsForceFlag reports whether s matches any standard force / auto-yes flag.
func IsForceFlag(s string) bool {
	return s == FlagForce || s == FlagForceShort || s == FlagYes || s == FlagYesShort
}

// IsDryRunFlag reports whether s matches any standard dry-run flag variant.
func IsDryRunFlag(s string) bool {
	return s == FlagDryRun || s == FlagDryRunShort || s == FlagDryRunAlt
}

// IsJSONFlag reports whether s matches standard json output flags.
func IsJSONFlag(s string) bool {
	return s == FlagJSON || s == "-j"
}

// IsVerboseFlag reports whether s matches standard verbose flags.
func IsVerboseFlag(s string) bool {
	return s == FlagVerbose || s == FlagVerboseShort
}
