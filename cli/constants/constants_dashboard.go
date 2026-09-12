package constants

// gitmap:cmd top-level
// Dashboard CLI commands.
const (
	CmdDashboard      = "dashboard"
	CmdDashboardAlias = "dash"
)

// Dashboard help text.
const HelpDashboard = "  dashboard (dash)    Generate an interactive HTML dashboard for a repo"

// Dashboard flag descriptions.
const (
	FlagRecent         = "recent"
	FlagDescDashRecent = "Show and highlight recent items only (e.g., last 7 days)"
	FlagDescDashLimit  = "Maximum number of commits to include"
	FlagDescDashSince  = "Only include commits after this date (YYYY-MM-DD)"
	FlagDescDashOpen   = "Open the generated dashboard in the default browser"
	FlagDescNoMerges   = "Exclude merge commits from the output"
	FlagDescDashOutDir = "Output directory for dashboard files"
)

// Dashboard output filenames.
const (
	DashboardJSONFile = "dashboard.json"
	DashboardHTMLFile = "dashboard.html"
	DashboardOutDir   = ".gitmap/output"
)

// Dashboard git log format — pipe-delimited fields:
// full SHA | short SHA | author name | author email | ISO date | subject | parent hashes.
const GitLogDashFormat = "%H|%h|%an|%ae|%aI|%s|%P"

// Dashboard git branch format — pipe-delimited fields:
// refname short | objectname short | creator date ISO.
const GitBranchDashFormat = "%(refname:short)|%(objectname:short)|%(creatordate:iso-strict)"

// Dashboard git tag format — pipe-delimited fields:
// refname short | objectname short | creator date ISO.
const GitTagDashFormat = "%(refname:short)|%(objectname:short)|%(creatordate:iso-strict)"

// Dashboard terminal messages.
const (
	MsgDashCollecting = ColorCyan + "Collecting repository data..." + ColorReset
	MsgDashWriteJSON  = ColorGreen + "✓" + ColorReset + " Wrote " + ColorCyan + "%s" + ColorReset + " (%d commits, %d authors)\n"
	MsgDashWriteHTML  = ColorGreen + "✓" + ColorReset + " Wrote " + ColorCyan + "%s" + ColorReset + "\n"
	MsgDashGenerated  = ColorGreen + "✓" + ColorReset + " Dashboard generated in " + ColorCyan + "%s" + ColorReset + "\n"
	MsgDashOpening    = "Opening dashboard in browser..."
)

// Dashboard error messages.
const (
	ErrDashNotRepo   = "Current directory is not a Git repository."
	ErrDashWriteJSON = "Failed to write dashboard JSON at %s: %v (operation: write)\n"
	ErrDashWriteHTML = "Failed to write dashboard HTML at %s: %v (operation: write)\n"
	ErrDashCollect   = "Failed to collect repository data: %v\n"
)
