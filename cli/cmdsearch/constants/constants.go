package constants

// Command name literals for cmdsearch.
const (
	CmdSearch              = "search"
	CmdRepoSearch          = "repo-search"
	CmdRepoRegex           = "repo-regex"
	CmdRepoSearchJSON      = "repo-search-json"
	CmdRepoSearchRegexJSON = "repo-search-regex-json"
	CmdReplaceRegex        = "replace-regex"
	CmdSearchReplaceAll    = "search-replace-all"
	CmdAUMSearch           = "aum search"
)

// CLI flag identifiers.
const (
	FlagQuery           = "query"
	FlagQueryShort      = "q"
	FlagPattern         = "pattern"
	FlagDir             = "dir"
	FlagDirShort        = "d"
	FlagRoot            = "root"
	FlagPath            = "path"
	FlagPathShort       = "p"
	FlagExt             = "ext"
	FlagExtShort        = "e"
	FlagRegex           = "regex"
	FlagRegexShort      = "r"
	FlagIgnoreCase      = "ignore-case"
	FlagIgnoreCaseShort = "i"
	FlagLimit           = "limit"
	FlagLimitShort      = "l"
	FlagAI              = "ai"
	FlagJSON            = "json"
	FlagJSONShort       = "j"
	FlagHelp            = "help"
	FlagHelpShort       = "h"
)

// Flag CLI arguments, prefixes, and entries.
const (
	FlagArgQuery            = "--query"
	FlagArgQueryShort       = "-q"
	FlagArgQueryPrefix      = "--query="
	FlagArgPattern          = "--pattern"
	FlagArgPatternPrefix    = "--pattern="
	FlagArgDir              = "--dir"
	FlagArgDirShort         = "-d"
	FlagArgRoot             = "--root"
	FlagArgPath             = "--path"
	FlagArgPathShort        = "-p"
	FlagArgExt              = "--ext"
	FlagArgExtShort         = "-e"
	FlagArgRegex            = "--regex"
	FlagArgRegexShort       = "-r"
	FlagArgIgnoreCase       = "--ignore-case"
	FlagArgIgnoreCaseShort  = "-i"
	FlagArgLimit            = "--limit"
	FlagArgLimitShort       = "-l"
	FlagArgLimitPrefix      = "--limit="
	FlagArgAI               = "--ai"
	FlagArgAIShort          = "-ai"
	FlagArgJSON             = "--json"
	FlagArgJSONShort        = "-j"
	FlagArgHelp             = "--help"
	FlagArgHelpShort        = "-h"
	FlagPrefixDash          = "-"
	FlagHelpEntryJSON       = "-j, --json"
	FlagHelpEntryPath       = "-p, --path <dir>"
	FlagHelpEntryExt        = "-e, --ext <exts>"
	FlagHelpEntryRegex      = "-r, --regex"
	FlagHelpEntryIgnoreCase = "-i, --ignore-case"
	FlagHelpEntryLimit      = "-l, --limit <n>"
	FlagHelpEntryHelp       = "-h, --help"
)

// Default execution limits, ignored directories, and system paths.
const (
	DefaultLocalSearchLimit = 20
	DirNodeModules          = "node_modules"
	DirVendor               = "vendor"
	DirBin                  = "bin"
	DirPrefixDot            = "."
	PathRepoSearchOutput    = ".gitmap/output/repo_search"
)

// JSON serialization constants.
const (
	EmptyJSONArray = "[]"
)

// UI header, section titles, and entry labels.
const (
	HelpTitleGlobalSearch   = "Global Indexed Search (gitmap search)"
	SectionSearchOperations = "Search Operations"
	SectionSQLiteCaching    = "High-Speed SQLite Caching"
	HelpEntryQuery          = "<query>"
	HelpEntryHistory        = "history"
	HelpEntryClean          = "clean"
	HelpEntryDH2DCache      = "DH2D SQLite Cache"
	HelpEntryRAMHotCache    = "RAM Hot-Cache"
)

// Help tips and descriptions.
const (
	TipHotQueryCache       = "Hot queries are cached in SQLite with instant DH2D lookup (< 0.05ms)."
	TipHistoryInspect      = "Use 'gitmap search history' to inspect query frequency and hit counts."
	DescJSONFormat         = "Emit results in structured JSON format"
	DescPathFilter         = "Restrict search to target subdirectory"
	DescExtFilter          = "Comma-separated extensions filter (.go, .ts)"
	DescRegexMatch         = "Evaluate search query as regular expression"
	DescIgnoreCase         = "Case-insensitive matching fast-path"
	DescResultLimit        = "Maximum number of search results to return"
	DescHelpMenu           = "Show this search help menu"
	DescStreamingSearch    = "Multi-core streaming keyword search across repository"
	DescSearchHistory      = "Display recent search history and hot query rankings"
	DescSearchClean        = "Purge search history cache and vacuum search database"
	DescDH2DCache          = "Deterministic SQL ID lookup caching verified matches"
	DescRAMHotCache        = "Auto-promotes frequent queries for microsecond retrieval"
	DescFastIndexedSearch  = "Fast indexed keyword and symbol search across repositories using SplitDB."
	DescRepoSearchMetadata = "Search indexed repositories by keyword with structured metadata."
)

// Usage templates.
const (
	UsageSearch            = "  gitmap search <query> [--limit <n>] [--ai]"
	UsageSearchFlags       = "gitmap search <query> [flags]"
	UsageSearchClean       = "  gitmap search clean"
	UsageSearchCleanSimple = "gitmap search clean"
	UsageSearchHistory     = "gitmap search history"
	UsageAUMSearch         = "gitmap aum search <query> [dir] [flags]"
	UsageReplaceRegex      = "  gitmap replace-regex <regex> <replacement>"
	UsageRepoSearch        = "  gitmap repo-search <query> [--limit <n>]"
	UsageRepoRegex         = "Usage: gitmap repo-regex <regex> [--limit <n>]"
)

// Examples.
const (
	ExampleSearch1      = "  gitmap search \"Resolve-Version\""
	ExampleSearch2      = "  gitmap search \"AppError\" --limit 10"
	ExampleSearch3      = "  gitmap search \"type SearchResult struct\" --ai"
	ExampleReplaceRegex = "  gitmap replace-regex \"v6\\.\\d+\\.\\d+\" \"v6.155.6\""
	ExampleRepoSearch1  = "  gitmap repo-search \"pipeline\""
	ExampleRepoSearch2  = "  gitmap repo-search \"constants\" --limit 5"
)

// Feedback messages and notification banners.
const (
	MsgSearchCleanSuccess    = "\n  %s✓ Search cache purged and database vacuumed (%d entries cleared).%s\n\n"
	MsgRepoSearchClear       = "repo-search clear: Cache cleared for current folder."
	MsgRepoRegexHistoryEmpty = "repo-regex history: No history recorded yet."
	MsgRepoRegexClear        = "repo-regex clear: Cache cleared for current folder."
	MsgSearchReplaceAllReset = "search-replace-all reset: Databases cleared."
	MsgSearchReplaceAllExec  = "search-replace-all executed."
	MsgReplaceRegexExec      = "replace-regex executed. (Requires file-writing engine implementation)"
	MsgUntrackedHintTip      = "  Tip: run 'gitmap scan' to index repositories or cd into a tracked repository."
)

// Formatting templates.
const (
	FmtHotCacheMatch       = "⚡ [AUM Hot-Cache %s] Instant match (<0.04ms) — %d result(s)\n"
	FmtSearchCompleted     = "🔍 [AUM Search ID: %s] Completed in %d ms (%d matches)\n"
	FmtFoundHeader         = "Found in %s at position %d"
	FmtSearchingDir        = "Searching for %q in %s...\n"
	FmtMatchingFilesHeader = "\n%sFound %d matching file(s):%s\n"
	FmtBulletItem          = "  • %s%s%s\n"
	FmtMatchingReposHeader = "\n%sFound %d matching repository(ies):%s\n"
	FmtUntrackedNotFound   = "\n%sNo matches found for %q in untracked folder.%s\n"
	FmtRepoWithLocation    = "%s (%s)"
)

// Error message literals.
const (
	ErrDefaultError         = "error"
	ErrRepoSearchFatal      = "repo search fatal error"
	ErrRepoRegexSearchFatal = "repo regex search fatal error"
	ErrRepoNotFound         = "current directory is not a tracked gitmap repository. run 'gitmap scan' first"
	FmtErrFindRepoByPath    = "find repo by path %s failed"
)
