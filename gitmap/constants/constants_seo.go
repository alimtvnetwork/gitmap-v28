package constants

// gitmap:cmd top-level
// SEO-write command constants.
const (
	CmdSEOWrite       = "seo-write"
	CmdSEOWriteAlias  = "sw"
	CmdCreateTemplate = "ct" // gitmap:cmd skip
)

// SEO-write flag names.
const (
	FlagSEOCSV            = "csv"
	FlagSEOURL            = "url"
	FlagSEOService        = "service"
	FlagSEOArea           = "area"
	FlagSEOCompany        = "company"
	FlagSEOPhone          = "phone"
	FlagSEOEmail          = "email"
	FlagSEOAddress        = "address"
	FlagSEOMaxCommits     = "max-commits"
	FlagSEOInterval       = "interval"
	FlagSEOFiles          = "files"
	FlagSEORotateFile     = "rotate-file"
	FlagSEODryRun         = "dry-run"
	FlagSEOTemplate       = "template"
	FlagSEOCreateTemplate = "create-template"
	FlagSEOAuthorName     = "author-name"
	FlagSEOAuthorEmail    = "author-email"
)

// SEO-write flag descriptions.
const (
	FlagDescSEOCSV            = "CSV file with title,description columns"
	FlagDescSEOURL            = "Website URL to glorify in commit messages"
	FlagDescSEOService        = "Service name for template placeholders"
	FlagDescSEOArea           = "Geographic area for template placeholders"
	FlagDescSEOCompany        = "Company name for template placeholders"
	FlagDescSEOPhone          = "Phone number for template placeholders"
	FlagDescSEOEmail          = "Email address for template placeholders"
	FlagDescSEOAddress        = "Physical address for template placeholders"
	FlagDescSEOMaxCommits     = "Stop after N commits (0 = run until Ctrl+C)"
	FlagDescSEOInterval       = "Random delay range in seconds (min-max)"
	FlagDescSEOFiles          = "Glob pattern to select files for staging"
	FlagDescSEORotateFile     = "File to modify during rotation mode"
	FlagDescSEODryRun         = "Preview commit messages without executing"
	FlagDescSEOTemplate       = "Load templates from a custom JSON file"
	FlagDescSEOCreateTemplate = "Generate a sample seo-templates.json and exit"
	FlagDescSEOAuthorName     = "Git author name for commits"
	FlagDescSEOAuthorEmail    = "Git author email for commits"
)

// SEO-write defaults.
const (
	SEODefaultIntervalMin = 60
	SEODefaultIntervalMax = 120
	SEODefaultInterval    = "60-120"
	SEOSeedFile           = "data/seo-templates.json"
	SEOTemplateOutputFile = "seo-templates.json"
)

// SEO-write placeholder tokens.
const (
	PlaceholderService = "{service}"
	PlaceholderArea    = "{area}"
	PlaceholderURL     = "{url}"
	PlaceholderCompany = "{company}"
	PlaceholderPhone   = "{phone}"
	PlaceholderEmail   = "{email}"
	PlaceholderAddress = "{address}"
)

// SEO-write terminal messages.
const (
	MsgSEOHeader          = "seo-write: %d commits planned (interval: %d-%ds)\n"
	MsgSEOHeaderUnlimited = "seo-write: unlimited commits (interval: %d-%ds)\n"
	MsgSEOCommit          = "  [%d/%d] ✓ %q → pushed (file: %s)\n"
	MsgSEOCommitOpen      = "  [%d] ✓ %q → pushed (file: %s)\n"
	MsgSEORotation        = "  [%d/%d] ↻ rotation: %s (append → commit → revert → commit)\n"
	MsgSEORotationOpen    = "  [%d] ↻ rotation: %s (append → commit → revert → commit)\n"
	MsgSEODone            = "  Done: %d commits pushed in %s\n"
	MsgSEODryTitle        = "  [%d] title: %s\n"
	MsgSEODryDesc         = "        desc:  %s\n"
	MsgSEODryAuthor       = "  author: %s\n"
	MsgSEOCreated         = "Created %s with sample templates\n"
	MsgSEOSeeded          = "Seeded %d templates into database\n"
	MsgSEOGraceful        = "\nGraceful shutdown: finishing current commit...\n"
	MsgSEOWaiting         = "  waiting %ds before next commit...\n"
)

// SEO-write error messages.
const (
	ErrSEOURLRequired    = "error: --url is required in template mode\n"
	ErrSEOCSVRead        = "error: failed to read CSV file at %s: %v (operation: read)\n"
	ErrSEOCSVEmpty       = "error: CSV file contains no rows\n"
	ErrSEOTemplateRead   = "error: failed to read template file at %s: %v (operation: read)\n"
	ErrSEOTemplateEmpty  = "error: no templates found\n"
	ErrSEOIntervalFmt    = "error: invalid --interval format, expected min-max (e.g. 60-120)\n"
	ErrSEONoFiles        = "error: no files found matching pattern\n"
	ErrSEORotateNotFound = "error: rotate file not found at %s (operation: resolve, reason: file does not exist)\n"
	ErrSEOGitStage       = "error: git add failed: %v\n"
	ErrSEOGitCommit      = "error: git commit failed: %v\n"
	ErrSEOGitPush        = "error: git push failed: %v\n"
	ErrSEOSeedRead       = "error: failed to read seed file at %s: %v (operation: read)\n"
	ErrSEOCreateWrite    = "error: failed to write template file at %s: %v (operation: write)\n"
	ErrSEODBInsert       = "error: failed to insert template: %v\n"
)

// SEO-write help text.
const (
	HelpSEOWrite       = "  seo-write (sw)      Automated SEO commit scheduler with templates"
	HelpSEOWriteFlags  = "SEO-write flags:"
	HelpSEOCSV         = "  --csv <path>        CSV file with title,description columns"
	HelpSEOURL         = "  --url <url>         Website URL to glorify in commit messages (required)"
	HelpSEOService     = "  --service <name>    Service name for template placeholders"
	HelpSEOArea        = "  --area <name>       Geographic area for template placeholders"
	HelpSEOCompany     = "  --company <name>    Company name for template placeholders"
	HelpSEOPhone       = "  --phone <number>    Phone number for template placeholders"
	HelpSEOEmail       = "  --email <addr>      Email address for template placeholders"
	HelpSEOAddress     = "  --address <addr>    Physical address for template placeholders"
	HelpSEOMaxCommits  = "  --max-commits <N>   Stop after N commits (0 = unlimited, default: 0)"
	HelpSEOInterval    = "  --interval <min-max> Random delay in seconds (default: 60-120)"
	HelpSEOFilesFlag   = "  --files <glob>      Glob pattern to select files for staging"
	HelpSEORotate      = "  --rotate-file <f>   File to modify in rotation mode"
	HelpSEODryRunFlag  = "  --dry-run           Preview commit messages without executing"
	HelpSEOTemplateF   = "  --template <path>   Load templates from a custom JSON file"
	HelpSEOCreateTpl   = "  --create-template   Generate sample seo-templates.json (alias: ct)"
	HelpSEOAuthorName  = "  --author-name <n>   Git author name for commits"
	HelpSEOAuthorEmail = "  --author-email <e>  Git author email for commits"
)

// CommitTemplate table (v15: singular + CommitTemplateId PK).
const TableCommitTemplate = "CommitTemplate"

// Legacy plural retained for migration detection.
const LegacyTableCommitTemplates = "CommitTemplates"

// SQL: create CommitTemplate table.
const SQLCreateCommitTemplate = `CREATE TABLE IF NOT EXISTS CommitTemplate (
	CommitTemplateId INTEGER PRIMARY KEY AUTOINCREMENT,
	Kind             TEXT NOT NULL,
	Template         TEXT NOT NULL,
	CreatedAt        TEXT NOT NULL DEFAULT (datetime('now'))
)`

// SQL: commit-template operations (v15).
const (
	SQLInsertTemplate        = "INSERT INTO CommitTemplate (Kind, Template) VALUES (?, ?)"
	SQLSelectTemplatesByKind = "SELECT CommitTemplateId, Kind, Template, CreatedAt FROM CommitTemplate WHERE Kind = ? ORDER BY CreatedAt"
	SQLCountTemplates        = "SELECT COUNT(*) FROM CommitTemplate"
	SQLDropCommitTemplate    = "DROP TABLE IF EXISTS CommitTemplate"
	SQLDropCommitTemplates   = "DROP TABLE IF EXISTS CommitTemplates" // legacy
)

// Template kinds.
const (
	TemplateKindTitle       = "title"
	TemplateKindDescription = "description"
)
