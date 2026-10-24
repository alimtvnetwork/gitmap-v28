package constants

// Subcommand names and aliases for cmddiff.
const (
	CmdDiff              = "diff"
	CmdDiffAlias         = "df"
	CmdCompare           = "compare"
	CmdCompareAlias      = "comp"
	CmdDiffProfiles      = "diff-profiles"
	CmdDiffProfilesAlias = "dp"
)

// Flag names used in FlagSet parsing.
const (
	FlagDiffJSON             = "json"
	FlagDiffOnlyConflicts    = "only-conflicts"
	FlagDiffOnlyMissing      = "only-missing"
	FlagDiffIncludeIdentical = "include-identical"
	FlagDiffIncludeVCS       = "include-vcs"
	FlagDiffIncludeNodeMods  = "include-node-modules"
	FlagDiffStat             = "stat"
	FlagDiffCached           = "cached"
	FlagDiffStaged           = "staged"
	FlagDiffNameOnly         = "name-only"
	FlagDiffHelp             = "help"
)

// Flag CLI argument strings and formatted menu entries.
const (
	FlagArgStat               = "--stat"
	FlagArgCached             = "--cached"
	FlagArgStaged             = "--staged"
	FlagArgNameOnly           = "--name-only"
	FlagArgJSON               = "--json"
	FlagArgJSONShort          = "-j"
	FlagArgHelp               = "--help"
	FlagArgHelpShort          = "-h"
	FlagHelpEntryCachedStaged = "--cached, --staged"
	FlagHelpEntryJSON         = "-j, --json"
	FlagHelpEntryHelp         = "-h, --help"
)

// Flag descriptions for help text and flag definitions.
const (
	FlagDescJSON             = "emit JSON instead of text"
	FlagDescOnlyConflicts    = "show only conflicting files"
	FlagDescOnlyMissing      = "show only missing-side files"
	FlagDescIncludeIdentical = "include identical files in output"
	FlagDescIncludeVCS       = "include .git/ in the walk"
	FlagDescIncludeNodeMods  = "include node_modules/ in the walk"
	FlagDescStat             = "Show compact change statistics (insertions/deletions)"
	FlagDescCachedStaged     = "Diff staged changes against repository HEAD"
	FlagDescNameOnly         = "Show only names of modified files"
	FlagDescJSONStructured   = "Emit diff metadata in structured JSON format"
	FlagDescHelp             = "Show this diff help menu"
)

// Help menu layout strings, section titles, and tips.
const (
	HelpTitleDiff         = "Diff & Inspection Suite (gitmap diff)"
	HelpUsageDiffTarget   = "gitmap diff [target] [flags]"
	HelpUsageDiffProfiles = "gitmap diff-profiles <profile1> <profile2>"
	HelpSectionRepoDiff   = "Repository Diff Actions"
	HelpSectionProfiles   = "Profile Comparison"
	HelpEntryDiffPath     = "diff [path]"
	HelpDescDiffPath      = "Inspect unstaged file modifications with colored diffs"
	HelpEntryDiffCached   = "diff --cached"
	HelpDescDiffCached    = "Inspect staged changes prepared for commit"
	HelpEntryDiffStat     = "diff --stat"
	HelpDescDiffStat      = "Compact overview of affected files and line deltas"
	HelpEntryDiffProfiles = "diff-profiles <p1> <p2>"
	HelpDescDiffProfiles  = "Compare software packages, aliases, and settings between profiles"
	HelpTipUncommitted    = "Inspect uncommitted modifications across polyglot repositories."
	HelpTipDrift          = "Use 'gitmap diff-profiles' to detect drift between workstation setups."
)

// Error wrap tokens and log prefixes for error management.
const (
	WrapResolveEndpointLeft  = "resolve-endpoint left"
	WrapResolveEndpointRight = "resolve-endpoint right"
	WrapGuardPaths           = "guard-paths"
	WrapDiffFailedPrefix     = "diff failed: "
	WrapReportFailedPrefix   = "report failed: "
	WrapDiffFailed           = "diff failed: [diff]"
	WrapReportFailed         = "report failed: [diff]"
	LogPrefixDiff            = "[diff]"
	ErrDiffUsageFmt          = "Usage: gitmap diff LEFT RIGHT [flags]\n"
	ErrDiffSameFolder        = "error: LEFT and RIGHT resolve to the same folder: %s"
)

// Positional argument limits and slice index constants.
const (
	ExpectedDiffPositionalArgs = 2
	ArgsSubcommandOffset       = 2
	IndexDiffLeft              = 0
	IndexDiffRight             = 1
)

// Exit codes for diff execution.
const (
	ExitUsageError = 2
)
