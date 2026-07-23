package constants

// TUI defaults.
const DefaultDashboardRefresh = 30

// gitmap:cmd top-level
// TUI command.
const (
	CmdInteractive      = "interactive"
	CmdInteractiveAlias = "i"
)

// TUI flag constants.
const (
	FlagRefresh     = "refresh"
	FlagDescRefresh = "Dashboard auto-refresh interval in seconds"
)

// TUI help text.
const (
	HelpInteractive      = "  interactive (i)     Launch interactive TUI with repo browser and batch actions"
	HelpInteractiveFlags = "Interactive Flags:"
	HelpRefresh          = "  --refresh <seconds>  Dashboard auto-refresh interval (default: config or 30s)"
)

// TUI view labels.
const (
	TUIViewBrowser      = "Repos"
	TUIViewActions      = "Actions"
	TUIViewGroups       = "Groups"
	TUIViewDashboard    = "Status"
	TUIViewReleases     = "Releases"
	TUIViewTempReleases = "Temp"
	TUIViewZipGroups    = "Zip Groups"
	TUIViewAliases      = "Aliases"
	TUIViewLogs         = "Logs"
)

// TUI status messages.
const (
	TUITitle          = "gitmap interactive"
	TUISearchPrompt   = "Search: "
	TUINoRepos        = "No repositories found. Run 'gitmap scan' first."
	TUINoGroups       = "No groups found. Press 'c' to create one."
	TUINoSelection    = "No repos selected. Use Space to select."
	TUIConfirmDelete  = "Delete group '%s'? (y/n)"
	TUIGroupCreated   = "Group '%s' created"
	TUIGroupDeleted   = "Group '%s' deleted"
	TUIActionPull     = "Pulling %d repo(s)..."
	TUIActionExec     = "Executing across %d repo(s)..."
	TUIActionStatus   = "Checking status of %d repo(s)..."
	TUIActionComplete = "Action complete: %d success, %d failed"
	TUIRefreshing     = "Refreshing..."
	TUIQuitHint       = "q/esc: quit"
	TUITabHint        = "tab: switch view"
	TUISelectHint     = "space: select  enter: detail  /: search"
	TUIBatchHint      = "p: pull  x: exec  s: status  g: add to group"
	TUIGroupHint      = "c: create  d: delete  enter: show members"
	TUIDashHint       = "r: refresh"
	TUIZGHint         = "enter: show items  r: refresh  c: create  d: delete"
	TUIAliasHint      = "r: refresh  c: set alias  d: remove"
)

// TUI zip group messages.
const (
	TUIZGEmpty      = "No zip groups defined. Use 'gitmap z create <name>' to create one."
	TUIZGRefreshed  = "Zip groups refreshed."
	TUIZGCreateHint = "Use CLI: gitmap z create <name>"
)

// TUI alias messages.
const (
	TUIAliasEmpty      = "No aliases defined. Use 'gitmap alias set <alias> <slug>' to create one."
	TUIAliasRefreshed  = "Aliases refreshed."
	TUIAliasCreateHint = "Use CLI: gitmap alias set <alias> <slug>"
	TUIAliasDeleteHint = "Remove alias '%s'? Use CLI: gitmap alias remove %s"
)

// TUI column headers.
const (
	TUIColSlug     = "Slug"
	TUIColBranch   = "Branch"
	TUIColPath     = "Path"
	TUIColType     = "Type"
	TUIColStatus   = "Status"
	TUIColAhead    = "Ahead"
	TUIColBehind   = "Behind"
	TUIColStash    = "Stash"
	TUIColGroup    = "Group"
	TUIColMembers  = "Members"
	TUIColVersion  = "Version"
	TUIColTag      = "Tag"
	TUIColDraft    = "Draft"
	TUIColLatest   = "Latest"
	TUIColSource   = "Source"
	TUIColDate     = "Date"
	TUIColCommand  = "Command"
	TUIColAlias    = "Alias"
	TUIColArgs     = "Args"
	TUIColDuration = "Duration"
	TUIColExit     = "Exit"
)

// TUI log messages.
const (
	TUILogEmpty        = "No command history found. Run some gitmap commands first."
	TUILogHint         = "enter: detail  r: refresh  /: filter"
	TUILogFilterActive = "  Filter: %s (%d matches)"
	TUILogNoMatch      = "  No logs match the current filter."
)

// TUI release messages.
const (
	TUIRelEmpty = "No releases found. Use 'gitmap release' to create one."
	TUIRelHint  = "enter: detail  r: refresh  n: new release"
)

// TUI temp-release messages.
const (
	TUITREmpty     = "No temp-release branches found. Use 'gitmap tr <count> <pattern>' to create."
	TUITRHint      = "enter: detail  g: group by prefix  r: refresh"
	TUIColTRBranch = "Branch"
	TUIColTRPrefix = "Prefix"
	TUIColTRSeq    = "Seq"
	TUIColTRCommit = "Commit"
)

// TUI release trigger messages.
const (
	TUIRelTriggerTitle     = "  New Release"
	TUIRelTriggerCmd       = "gitmap release %s"
	TUIRelTriggerBumpCmd   = "gitmap release %s"
	TUIRelTriggerNavHint   = "  ↑/↓: select  enter: confirm  esc: cancel"
	TUIRelTriggerVerPrompt = "  Version: "
	TUIRelTriggerTypeHint  = "  enter: confirm  esc: back"
	TUIRelTriggerReady     = "  Release Command"
	TUIRelTriggerRunHint   = "  Run this command in your terminal.  esc: back"
)

// TUI errors.
const (
	ErrTUINoTerminal = "interactive mode requires a terminal — use standard commands instead"
	ErrTUIDBOpen     = "failed to open database for interactive mode: %v"
)
