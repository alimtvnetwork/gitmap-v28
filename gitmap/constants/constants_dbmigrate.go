package constants

// gitmap:cmd top-level
// `gitmap db-migrate` (alias `dbm`) — explicitly run database schema
// migrations. Safe to invoke at any time; designed to be run automatically
// after `gitmap update` and on first launch after a fresh install.
const (
	CmdDBMigrate      = "db-migrate"
	CmdDBMigrateAlias = "dbm"

	HelpDBMigrate = "  db-migrate (dbm)    Run pending database schema migrations (safe, idempotent)"

	FlagDBMigrateVerbose = "verbose"
	FlagDescDBMigrateV   = "Print every migration step (otherwise summary only)"

	MsgDBMigrateRunning = "▸ Running gitmap database migrations...\n"
	MsgDBMigrateDoneFmt = "  ✓ Migrations complete (%d tables ensured, %d steps applied, %d warnings).\n"
	MsgDBMigrateStepFmt = "    • %s\n"
	MsgDBMigrateNoWork  = "  ✓ Database schema is already up to date — nothing to migrate.\n"
	ErrDBMigrateFailFmt = "Error: gitmap db-migrate failed: %v\n"

	// Auto-run hook for the post-update flow.
	MsgDBMigratePostUpdate = "▸ Running database migrations after update...\n"
	WarnDBMigratePostFail  = "  ⚠ Post-update migration failed: %v\n      Run `gitmap db-migrate --verbose` to retry.\n"
)
