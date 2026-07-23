// Package constants — inject/open idempotency tracking (schema v25).
//
// Adds two TEXT columns to the Repo table that record the last time the
// folder was registered with each external tool. Empty string ("") means
// "never injected" — distinct from a missing column on legacy DBs which
// is auto-handled by addColumnIfNotExists().
//
// Both `gitmap inject` and `gitmap open` consult these timestamps to
// avoid spamming Desktop/VS Code on every invocation, and reset them
// when the user passes --force / -f.
package constants

// ALTER statements for additive schema-v25 columns. Idempotent — the
// migration helper swallows "duplicate column" benignly.
const (
	SQLAddRepoLastInjectedDesktopAt = "ALTER TABLE Repo ADD COLUMN LastInjectedDesktopAt TEXT DEFAULT ''"
	SQLAddRepoLastInjectedVSCodeAt  = "ALTER TABLE Repo ADD COLUMN LastInjectedVSCodeAt TEXT DEFAULT ''"
)

// SQL: read both timestamps for a given absolute path. Empty strings
// when the row exists but has never been injected; sql.ErrNoRows when
// the path is unknown to the DB.
const SQLSelectInjectTimestamps = `
SELECT LastInjectedDesktopAt, LastInjectedVSCodeAt
FROM Repo WHERE AbsolutePath = ?`

// SQL: stamp the per-tool timestamp to CURRENT_TIMESTAMP. Caller picks
// the column via fmt.Sprintf — only two known constants are ever
// substituted (ColInjectDesktop / ColInjectVSCode), so there is no
// SQL-injection surface.
const SQLUpdateInjectTimestampFmt = `
UPDATE Repo SET %s = CURRENT_TIMESTAMP, UpdatedAt = CURRENT_TIMESTAMP
WHERE AbsolutePath = ?`

// Column names referenced by SQLUpdateInjectTimestampFmt. Centralized
// so renames cannot drift between the schema and the writer.
const (
	ColInjectDesktop = "LastInjectedDesktopAt"
	ColInjectVSCode  = "LastInjectedVSCodeAt"
)

// InjectKind discriminates the two tool slots in the helper API. Kept
// as a typed constant so callers can't pass a typo'd column name.
type InjectKind int

const (
	InjectKindDesktop InjectKind = iota + 1
	InjectKindVSCode
)

// inject / open --force flag.
const (
	FlagInjectForce      = "force"
	FlagInjectForceShort = "f"
	FlagDescInjectForce  = "Re-register Desktop/VS Code even if already injected"
)

// Skip messages — printed when an action is suppressed because the
// per-tool timestamp is already set and --force was NOT passed.
const (
	MsgInjectSkipDesktopFmt = "  ↳ %sgithub-desktop:%s already injected (%s) — pass %s--force%s to re-register\n"
	MsgInjectSkipVSCodeFmt  = "  ↳ %svscode:%s already injected (%s) — pass %s--force%s to re-open\n"
	MsgInjectForceNotice    = "  " + ColorYellow + "⟳ --force: re-injecting %s into both tools" + ColorReset + "\n"
)
