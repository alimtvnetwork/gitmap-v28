package constants

// gitmap:cmd top-level
// Amend command.
const (
	CmdAmend          = "amend"
	CmdAmendAlias     = "am"
	CmdAmendList      = "amend-list"
	CmdAmendListAlias = "al"
)

// Amend flag names.
const (
	FlagAmendName      = "name"
	FlagAmendEmail     = "email"
	FlagAmendBranch    = "branch"
	FlagAmendDryRun    = "dry-run"
	FlagAmendForcePush = "force-push"
)

// Amend flag descriptions.
const (
	FlagDescAmendName      = "New author name for commits"
	FlagDescAmendEmail     = "New author email for commits"
	FlagDescAmendBranch    = "Target branch (default: current branch)"
	FlagDescAmendDryRun    = "Preview which commits would be amended"
	FlagDescAmendForcePush = "Auto-run git push --force-with-lease after amend"
)

// Amend modes.
const (
	AmendModeAll   = "all"
	AmendModeRange = "range"
	AmendModeHead  = "head"
)

// Amend audit directory.
const (
	AmendAuditDir        = ".gitmap/amendments"
	AmendAuditFilePrefix = "amend-"
)

// Amend terminal messages.
const (
	MsgAmendHeader      = "amend: rewriting %d commits from %s..%s (branch: %s)\n"
	MsgAmendHeaderAll   = "amend: rewriting %d commits on branch: %s\n"
	MsgAmendAuthor      = "  author: %q -> %q\n"
	MsgAmendProgress    = "  [%d/%d] %s - %s\n"
	MsgAmendDone        = "\nDone: %d commits amended\n"
	MsgAmendAuditFile   = "  Audit log: %s\n"
	MsgAmendAuditDB     = "  Database:  1 record saved to Amendments table\n"
	MsgAmendForcePush   = "  Force push: completed\n"
	MsgAmendWarnPush    = "Warning: Run 'git push --force-with-lease' to update the remote\n"
	MsgAmendDryHeader   = "amend (dry-run): %d commits would be rewritten\n"
	MsgAmendDryLine     = "  [%d] %s - %s (author: %s <%s>)\n"
	MsgAmendDrySkip     = "  No changes applied (dry-run mode)\n"
	MsgAmendCheckout    = "  Switching to branch: %s\n"
	MsgAmendReturn      = "  Returning to branch: %s\n"
	MsgAmendWarnRewrite = "Warning: This rewrites Git history and requires force-push.\n"
)

// Amend error messages.
const (
	ErrAmendNoFlags     = "error: at least one of --name or --email is required\n"
	ErrAmendCheckout    = "error: failed to checkout branch %s: %v\n"
	ErrAmendListCommits = "error: failed to list commits: %v\n"
	ErrAmendFilter      = "error: git filter-branch failed: %v\n"
	ErrAmendForcePush   = "error: force push failed: %v\n"
	ErrAmendAuditWrite  = "error: failed to write audit file: %v\n"
	ErrAmendCommitAmend = "error: git commit --amend failed: %v\n"
	ErrAmendNoCommits   = "error: no commits found in the specified range\n"
)

// Amend help text.
const (
	HelpAmend      = "  amend (am) [hash]   Rewrite author name/email on commits"
	HelpAmendList  = "  amend-list (al)     Show stored amendments from database (--limit N, --json, --branch)"
	HelpAmendFlags = "Amend flags:"
	HelpAmendName  = "  --name <name>       New author name"
	HelpAmendEmail = "  --email <email>     New author email"
	HelpAmendBr    = "  --branch <branch>   Target branch (default: current)"
	HelpAmendDry   = "  --dry-run           Preview which commits would be amended"
	HelpAmendForce = "  --force-push        Auto force-push after amend"
)

// Amend-list flag.
const FlagAmendListBranch = "--branch"

// Amend-list terminal messages.
const (
	MsgAmendListEmpty     = "No amendments found."
	MsgAmendListHeader    = "Amendments: %d record(s)\n"
	MsgAmendListSeparator = "──────────────────────────────────────────────────────────────────────────────────"
	MsgAmendListColumns   = "BRANCH          MODE    COMMITS  PREV AUTHOR              NEW AUTHOR               PUSHED  DATE"
	MsgAmendListRowFmt    = "%-15s %-7s %7d  %-12s %-12s %-12s %-12s %-6s  %s\n"
)

// Amend-list error messages.
const ErrAmendListFailed = "error: failed to list amendments: %v\n"

// Amendment table (v15: singular + AmendmentId PK).
const TableAmendment = "Amendment"

// Legacy plural retained for migration detection.
const LegacyTableAmendments = "Amendments"

// SQL: create Amendment table (v15).
const SQLCreateAmendment = `CREATE TABLE IF NOT EXISTS Amendment (
	AmendmentId   INTEGER PRIMARY KEY AUTOINCREMENT,
	Branch        TEXT NOT NULL,
	FromCommit    TEXT NOT NULL,
	ToCommit      TEXT NOT NULL,
	TotalCommits  INTEGER NOT NULL,
	PreviousName  TEXT DEFAULT '',
	PreviousEmail TEXT DEFAULT '',
	NewName       TEXT DEFAULT '',
	NewEmail      TEXT DEFAULT '',
	Mode          TEXT NOT NULL,
	ForcePushed   INTEGER DEFAULT 0,
	CreatedAt     TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: amendment operations (v15).
const (
	SQLInsertAmendment = `INSERT INTO Amendment (Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	SQLSelectAllAmendments = `SELECT AmendmentId, Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed, CreatedAt
		FROM Amendment ORDER BY CreatedAt DESC`

	SQLSelectAmendmentsByBranch = `SELECT AmendmentId, Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed, CreatedAt
		FROM Amendment WHERE Branch = ? ORDER BY CreatedAt DESC`

	SQLDropAmendment  = "DROP TABLE IF EXISTS Amendment"
	SQLDropAmendments = "DROP TABLE IF EXISTS Amendments" // legacy
)
