package constants

// gitmap:cmd top-level
// Watch CLI commands.
const (
	CmdWatch      = "watch"
	CmdWatchAlias = "w"
)

// Watch help text.
const HelpWatch = "  watch (w)           Live-refresh dashboard of repo status"

// Watch defaults.
const (
	WatchDefaultInterval = 30
	WatchMinInterval     = 5
)

// Watch ANSI control.
const (
	WatchClearScreen = "\033[2J\033[H"
)

// Watch display messages.
const (
	WatchBannerTop    = "╔══════════════════════════════════════╗"
	WatchBannerTitle  = "║          gitmap watch                ║"
	WatchBannerBottom = "╚══════════════════════════════════════╝"
	WatchRefreshFmt   = "gitmap watch — refreshing every %ds (Ctrl+C to stop)"
	WatchLastUpdFmt   = "Last updated: %s"
	WatchHeaderFmt    = "  %s%-22s %-10s %-16s %-6s %-8s %-6s%s\n"
	WatchRowFmt       = "  %-22s %s  %-16s %-6s %-8s %s\n"
	WatchErrorRowFmt  = "  %s%-22s %serror%s\n"
	WatchSummaryFmt   = "Repos: %d | Dirty: %d | Behind: %d | Stash: %d"
	WatchStoppedMsg   = "\ngitmap watch stopped."
)

// Watch table column headers.
var WatchTableColumns = []string{
	"REPO", "STATUS", "BRANCH", "AHEAD", "BEHIND", "STASH",
}

// Watch flag descriptions.
const (
	FlagDescWatchInterval = "Refresh interval in seconds (minimum 5)"
	FlagDescWatchNoFetch  = "Skip git fetch; use local refs only"
	FlagDescWatchJSON     = "Output single snapshot as JSON and exit"
)

// Watch error messages.
const (
	ErrWatchNoRepos = "No repos to watch. Run 'gitmap scan' first."
)
