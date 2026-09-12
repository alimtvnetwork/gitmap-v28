package constants

// gitmap:cmd top-level
// Diff-profiles CLI commands.
const (
	CmdDiffProfiles      = "diff-profiles"
	CmdDiffProfilesAlias = "dp"
)

// Diff-profiles help text.
const HelpDiffProfiles = "  diff-profiles (dp)  Compare repos across two database profiles"

// Diff-profiles messages.
const (
	MsgDPHeader         = "Comparing profiles: %s ↔ %s\n"
	MsgDPOnlyInHeader   = "\nONLY IN %s:"
	MsgDPOnlyInRowFmt   = "  %-20s %s\n"
	MsgDPDiffHeader     = "\nDIFFERENT:"
	MsgDPDiffNameFmt    = "  %s\n"
	MsgDPDiffDetailFmt  = "    %-10s %s\n"
	MsgDPSameFmt        = "\nSAME: %d repos (use --all to show)\n"
	MsgDPSameAllHeader  = "\nSAME:"
	MsgDPSameRowFmt     = "  %-20s %s\n"
	MsgDPSummaryFmt     = "\nSummary: %d only-left | %d only-right | %d different | %d same\n"
	MsgDPEmpty          = "Both profiles have no repos."
	ErrDPUsage          = "usage: gitmap diff-profiles <profileA> <profileB> [--all] [--json]\n"
	ErrDPProfileMissing = "profile not found: %s\n"
	ErrDPOpenFailed     = "failed to open profile database '%s': %v\n"
)
