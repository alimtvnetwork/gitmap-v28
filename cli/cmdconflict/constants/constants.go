package constants

// Command names and aliases for cmdconflict.
const (
	CmdConflicts      = "conflicts"
	CmdConflictsAlias = "conflict"
	CmdResolve        = "resolve"
	CmdResolveAlias   = "res"
)

// Flag names and argument tokens for CLI parsing.
const (
	FlagTheirs      = "theirs"
	FlagTheirsLong  = "--theirs"
	FlagTheirsShort = "-t"

	FlagOurs      = "ours"
	FlagOursLong  = "--ours"
	FlagOursShort = "-o"

	FlagIDE      = "ide"
	FlagIDELong  = "--ide"
	FlagIDEShort = "-i"
	FlagCodeLong = "--code"

	FlagAbort      = "abort"
	FlagAbortLong  = "--abort"
	FlagAbortShort = "-a"

	FlagRerere     = "rerere"
	FlagRerereLong = "--rerere"

	FlagHelp      = "help"
	FlagHelpLong  = "--help"
	FlagHelpShort = "-h"
)

// Strategy tokens and resolution keywords.
const (
	TokenTheirs  = "theirs"
	TokenOurs    = "ours"
	TokenIDE     = "ide"
	TokenCode    = "code"
	TokenResolve = "resolve"
	TokenAbort   = "abort"
	TokenRerere  = "rerere"
)

// Git CLI executable, subcommands, arguments, and config keys.
const (
	GitBin                    = "git"
	GitFlagDir                = "-C"
	GitSubCmdDiff             = "diff"
	GitArgNameOnly            = "--name-only"
	GitDiffFilterUnmerged     = "--diff-filter=U"
	GitSubCmdStatus           = "status"
	GitArgPorcelain           = "--porcelain"
	GitSubCmdCheckout         = "checkout"
	GitArgTheirs              = "--theirs"
	GitArgOurs                = "--ours"
	GitSubCmdMerge            = "merge"
	GitSubCmdRebase           = "rebase"
	GitSubCmdCherryPick       = "cherry-pick"
	GitArgAbort               = "--abort"
	GitConfigRerereEnabled    = "rerere.enabled"
	GitConfigRerereAutoupdate = "rerere.autoupdate"
	GitArgDoubleDash          = "--"
	GitArgAdd                 = "add"
	GitSubCmdRevParse         = "rev-parse"
	GitArgIsInsideWorkTree    = "--is-inside-work-tree"
	GitTrue                   = "true"
	GitSubCmdConfig           = "config"
	FlagReuseWindow           = "-r"
)

// ANSI terminal color escape codes and UI glyphs.
const (
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"

	GlyphCheck = "✓"
	GlyphWarn  = "⚠"

	TreePrefixBranch = "├──"
	TreePrefixLeaf   = "└──"
)

// Help sections and quick action lines.
const (
	HelpSectionQuickActions = "\n⚡ Quick Actions:\n"
	HelpActionIDE           = "   gitmap resolve %s --ide      Launch VS Code / IDE 3-Way Merge Editor\n"
	HelpActionTheirs        = "   gitmap resolve %s --theirs   Accept incoming remote changes\n"
	HelpActionOurs          = "   gitmap resolve %s --ours     Accept local current changes\n"
	HelpActionAbort         = "   gitmap resolve %s --abort    Abort merge or rebase\n"
	HelpActionCommit        = "   gitmap cpf \"resolve conflicts\" Stage, commit, and push once resolved\n\n"
)

// VS Code 3-Way Merge Editor resolution guide lines.
const (
	GuideDivider  = "--------------------------------------------------------------------------------\n"
	GuideTitle    = "  VS Code 3-Way Merge Editor — Resolution Guide\n"
	GuideStep1    = "  1. Press Ctrl+Shift+G (or Cmd+Shift+G on macOS) to open the Source Control view.\n"
	GuideStep2    = "  2. Under 'Merge Changes', click on any conflicted file.\n"
	GuideStep3    = "  3. Click 'Resolve in Merge Editor' in the bottom-right of the editor tab.\n"
	GuideStep4    = "  4. Check 'Accept Incoming' (remote) or 'Accept Current' (local) for each conflict.\n"
	GuideStep5    = "  5. Click 'Complete Merge' at the bottom-right and save the file (Ctrl+S).\n"
	GuideStep6    = "  6. In GitMap, commit and push your resolution:\n"
	GuideStep6Sub = "     gitmap cpf \"resolve merge conflicts\"\n"
)

// User-facing messages and notification format strings.
const (
	MsgNoConflictsDetectedFmt       = "%s No merge conflicts detected in %s\n"
	MsgNoConflictsAcrossTrackedFmt  = "%s No merge conflicts detected across %d tracked repositories.\n"
	MsgNoConflictsCurrentRepoFmt    = "%s No merge conflicts detected in current repository.\n"
	MsgNoConflictedFilesDetectedFmt = "%s No conflicted files detected in %s\n"
	MsgAcceptedTheirsFmt            = "%s Accepted incoming (theirs) changes on %d conflicted file(s) in %s\n"
	MsgAcceptedOursFmt              = "%s Accepted current (ours) changes on %d conflicted file(s) in %s\n"
	MsgNextCommitInstruction        = "  Next: Run 'gitmap cpf \"resolve merge conflicts\"' to commit and push."
	MsgIDELaunchedFmt               = "%s Launched %s for %s\n"
	MsgIDELaunchFailedFmt           = "%s Could not launch IDE: %v\n"
	MsgAutoLaunchingIDEFmt          = "  Auto-launching %s...\n"
	MsgMergeAbortedFmt              = "%s Merge aborted successfully in %s\n"
	MsgRebaseAbortedFmt             = "%s Rebase aborted successfully in %s\n"
	MsgCherryPickAbortedFmt         = "%s Cherry-pick aborted successfully in %s\n"
	MsgRerereEnabledFmt             = "%s Git rerere (Reuse Recorded Resolution) enabled for %s\n"
	HeaderRepoConflictsFmt          = "\n📁 Repository: %s (%s)\n"
	HeaderConflictedFilesCountFmt   = "   Conflicted Files (%d):\n"
	LineConflictFileFmt             = "   %s %s\n"
)

// Error messages and format strings.
const (
	ErrNoSupportedIDE       = "no supported IDE found on PATH (checked: code, cursor, windsurf, codium)"
	ErrNoAbortInProgressFmt = "no merge, rebase, or cherry-pick in progress to abort in %s"
	ErrResolveUsage         = "not inside a git repository and no repository specified.\nUsage: gitmap resolve <repo> [--theirs | --ours | --ide | --abort]"
	ErrRepoNotFoundFmt      = "repository %q not found or not a valid directory"
	ErrCannotReadRepoStore  = "cannot read repository store"
	ErrCannotListRepos      = "cannot list repositories"
)
