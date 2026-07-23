package constants

// Replace command flag names (see spec/04-generic-cli/15-replace-command.md).
const (
	ReplaceFlagYes     = "yes"
	ReplaceFlagYesS    = "y"
	ReplaceFlagDryRun  = "dry-run"
	ReplaceFlagQuiet   = "quiet"
	ReplaceFlagQuietS  = "q"
	ReplaceFlagExt     = "ext"
	ReplaceFlagExtCase = "ext-case"
	ReplaceSubcmdAudit = "--audit"
	ReplaceSubcmdAll   = "all"
)

// Replace --ext flag description and value separator.
const (
	FlagDescReplaceExt     = "Comma-separated extension allow-list (e.g. \".go,.md\"). Leading dot optional."
	FlagDescReplaceExtCase = "Extension match casing: sensitive | insensitive (default insensitive)."
	ReplaceExtSep          = ","
)

// Allowed values for --ext-case. The default is `insensitive` so the
// historical behavior is preserved when the flag is omitted.
const (
	ReplaceExtCaseSensitive   = "sensitive"
	ReplaceExtCaseInsensitive = "insensitive"
	ReplaceExtCaseDefault     = ReplaceExtCaseInsensitive
)

// Error format for an unknown --ext-case value. Kept here so the
// message stays in sync with the constants above.
const ErrReplaceBadExtCase = "replace: --ext-case must be %q or %q (got %q)\n"

// Replace command messages and errors. All user-facing text lives here
// to honor the no-magic-strings rule.
const (
	MsgReplaceScanning     = "replace: scanning %d files in %s\n"
	MsgReplaceFileMatch    = "replace: %s: %d matches (%s -> %s)\n"
	MsgReplaceFileMatchOne = "replace: %s: 1 match (%s -> %s)\n"
	MsgReplaceSummary      = "replace: %d files, %d replacements\n"
	MsgReplaceConfirmLit   = "Apply %d replacements across %d files? [y/N]: "
	MsgReplaceConfirmVer   = "Apply replacements for versions v%d..v%d -> v%d? [y/N]: "
	MsgReplaceApplied      = "replace: applied %d replacements across %d files\n"
	MsgReplaceAborted      = "replace: aborted by user\n"
	MsgReplaceAlreadyAtV1  = "replace: already at v1, nothing to upgrade\n"
	MsgReplaceAuditMatch   = "%s:%d: %s\n"
	MsgReplaceAuditClean   = "replace: no older version references found\n"
	MsgReplaceNoMatches    = "replace: no matches found\n"
	ErrReplaceNeedsArgs    = "replace: expected two quoted strings, -N, all, or --audit\n"
	ErrReplaceEmptyOld     = "replace: old text must not be empty\n"
	ErrReplaceBadN         = "replace: -N must be a positive integer (got %q)\n"
	ErrReplaceNoRemote     = "replace: cannot read git remote: %v\n"
	ErrReplaceVersionParse = "replace: cannot detect version from remote: %q (expected suffix -vN)\n"
	ErrReplaceWrite        = "replace: %s: %v\n"
	ErrReplaceWalk         = "replace: walk error: %v\n"
	ErrReplaceNotInRepo    = "replace: not inside a git repository: %v\n"
)

// Replace exclusion sets. Directory names matched by base name only.
var (
	ReplaceExcludedDirs = []string{".git", ".gitmap", ".release", "node_modules", "vendor"}
	// Path prefixes are matched against the path relative to repo root.
	ReplaceExcludedPrefixes = []string{".gitmap/release", ".gitmap/release-assets"}
)

// ReplaceBinarySniffBytes is the number of bytes scanned for null
// bytes when classifying a file as binary.
const ReplaceBinarySniffBytes = 8192
