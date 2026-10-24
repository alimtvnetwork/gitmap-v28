package constants

// Command name literals for cmdstatus.
const (
	CmdStatus      = "status"
	CmdStatusAlias = "st"
)

// CLI flag identifiers for cmdstatus.
const (
	FlagGroup      = "group"
	FlagGroupShort = "g"
	FlagAll        = "all"
	FlagDirty      = "dirty"
	FlagOnlyDirty  = "only-dirty"
	FlagAhead      = "ahead"
	FlagBehind     = "behind"
	FlagTable      = "table"
	FlagCompact    = "compact"
	FlagJSON       = "json"
	FlagJSONShort  = "j"
	FlagHelp       = "help"
	FlagHelpShort  = "h"
	FlagDirtyShort = "d"
)

// Flag CLI arguments and aliases.
const (
	FlagArgGroup       = "--group"
	FlagArgGroupShort  = "-g"
	FlagArgAll         = "--all"
	FlagArgDirty       = "--dirty"
	FlagArgOnlyDirty   = "--only-dirty"
	FlagArgAhead       = "--ahead"
	FlagArgBehind      = "--behind"
	FlagArgTable       = "--table"
	FlagArgCompact     = "--compact"
	FlagArgJSON        = "--json"
	FlagArgJSONShort   = "-j"
	FlagArgHelp        = "--help"
	FlagArgHelpShort   = "-h"
	FlagArgDirtyShort  = "-d"
	FlagHelpEntryDirty = "-d, --dirty"
	FlagHelpEntryJSON  = "-j, --json"
	FlagHelpEntryHelp  = "-h, --help"
)

// Flag descriptions for help and flag parsing.
const (
	FlagDescDirty       = "Display only repositories with uncommitted, unstaged, or unpushed changes"
	FlagDescFilterDirty = "Filter to only repositories with uncommitted working changes"
	FlagDescFilterAhead = "Filter to repositories with commits ahead of remote tracking"
	FlagDescFilterBehind = "Filter to repositories with remote commits waiting to pull"
	FlagDescJSON        = "Output status records in structured JSON array"
	FlagDescTable       = "Render rich ANSI status table with colored columns"
	FlagDescCompact     = "Single-line per repository compact output"
	FlagDescHelp        = "Show this status help menu"
)

// Table column header names for cmdstatus table presentation.
const (
	TableColRepository = "Repository"
	TableColBranch     = "Branch"
	TableColStatus     = "Status"
	TableColSync       = "Sync"
	TableColStash      = "Stash"
	TableColFiles      = "Files"
)

// UI string tokens, layout formatting, and status banners.
const (
	StatusNotFoundText         = "✖ not found"
	StatusCleanSuccessMsg      = "✨ All repositories are !clean"
	StatusTableColGap          = "   "
	StatusTableEmptyDash       = "-"
	StatusBatchProgressLabel   = "Status"
	StatusHelpTitle            = "Workspace Repository Status (gitmap status)"
	StatusHelpTipAudit         = "Quickly audit modified repositories across your entire multi-repo workspace."
	StatusHelpTipJSON          = "Pass '--json' to pipe dirty repo inventory directly to AI fspath."
	StatusSectionFilterModes   = "Filter Modes"
	StatusSectionPresentation  = "Presentation & Output"
	StatusHelpEntryAll         = "all repos"
	StatusHelpDescAll          = "Scan status across all repositories in current workspace"
	StatusHelpDescDirty        = "Show only modified, untracked, or staged repositories"
	StatusHelpEntryAheadBehind = "--ahead / --behind"
	StatusHelpDescAheadBehind  = "Isolate branches out of sync with upstream remotes"
	StatusHelpEntryTableMode   = "Table Mode"
	StatusHelpDescTableMode    = "Interactive branch, dirty count, ahead/behind matrix"
	StatusHelpEntryJSONMode    = "JSON Mode"
	StatusHelpDescJSONMode     = "Machine-readable status telemetry for automation scripts"
	UsageStatus                = "gitmap status [flags]"
	UsageStatusAlias           = "gitmap st [flags]"
	UsageStatusDirty           = "gitmap status --dirty"
)

// Remediation UI strings for missing repositories.
const (
	RemediationHeader       = "  To resolve missing repositories:"
	RemediationStep1        = "  1. Relocate to a new folder:"
	RemediationStep1Cmd     = "     $ gitmap scan-folder update <slug> <new-path>"
	RemediationStep2        = "  2. Untrack from database:"
	RemediationStep2CmdFmt  = "     $ gitmap rm %s\n"
	RemediationStep3        = "  3. Clear or reset tracking database:"
	RemediationStep3Cmd     = "     $ gitmap db-reset  (or: $ gitmap storage reset-errors)"
	MissingRepoBannerFmt    = "  ▲ %d missing repository(ies) detected:"
	MissingRepoItemBullet   = "     • "
	DirtySummaryTipFmt      = "  Tip: %d repository(ies) dirty. Run 'gitmap fix' or 'gitmap pull --fix' to remediate."
)
