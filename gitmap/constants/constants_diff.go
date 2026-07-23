package constants

// gitmap:cmd top-level
// Diff command identifiers and flag names.
//
// Spec: companion to spec/01-app/97-move-and-merge.md
const (
	CmdDiff      = "diff"
	CmdDiffAlias = "df"
)

// Diff flag names.
const (
	FlagDiffJSON             = "json"
	FlagDiffOnlyConflicts    = "only-conflicts"
	FlagDiffOnlyMissing      = "only-missing"
	FlagDiffIncludeIdentical = "include-identical"
	FlagDiffIncludeVCS       = "include-vcs"
	FlagDiffIncludeNodeMods  = "include-node-modules"
	FlagDiffNoColor          = "no-color"
)

// Diff log prefix and section headers.
const (
	LogPrefixDiff           = "[diff]"
	DiffSectionMissingRight = "Missing on RIGHT (would be added by merge-right / merge-both):"
	DiffSectionMissingLeft  = "Missing on LEFT (would be added by merge-left / merge-both):"
	DiffSectionConflicts    = "Conflicts (different content on both sides):"
	DiffSectionIdentical    = "Identical files (skipped by merge-*):"
	DiffSummaryFmt          = "%s summary: %d missing-on-left, %d missing-on-right, %d conflicts, %d identical\n"
	DiffNothingFmt          = "%s no differences detected.\n"
)

// Diff error messages.
const (
	ErrDiffUsageFmt     = "Usage: gitmap diff LEFT RIGHT [flags]\n"
	ErrDiffSameFolder   = "error: LEFT and RIGHT resolve to the same folder: %s"
	ErrDiffNotFolderFmt = "error: %q is not a directory (gitmap diff requires local folder endpoints; clone URLs first via `gitmap clone`)"
	ErrDiffMissingFmt   = "error: %q does not exist"
)
