package constants

// Transaction journal & revert (spec/04-generic-cli/28-transaction-revert.md).
//
// Two SQLite tables track every FS-mutating op + the per-file backups
// captured before the mutation. Backup directory layout (mirrored on
// disk):
//
//   <gitmap-binary-dir>/.gitmap/txn/<txn-id>/data/<repo>/<git-sha>/files/<rel>
//
// PascalCase + INTEGER PK AUTOINCREMENT — project-wide schema rules.

const (
	// SQLCreateTransaction defines the master journal row. Status moves
	// pending → committed → reverted (or pending → aborted on failure).
	// "Transaction" is a SQL reserved word in some dialects, so the
	// table name is double-quoted everywhere it appears.
	SQLCreateTransaction = `CREATE TABLE IF NOT EXISTS "Transaction" (
	TransactionId   INTEGER PRIMARY KEY AUTOINCREMENT,
	Kind            TEXT NOT NULL,
	Status          TEXT NOT NULL,
	Argv            TEXT NOT NULL,
	Cwd             TEXT NOT NULL,
	CreatedAt       INTEGER NOT NULL,
	CommittedAt     INTEGER,
	RevertedAt      INTEGER,
	ReverseSummary  TEXT NOT NULL,
	RepoSlug        TEXT NOT NULL,
	GitSha          TEXT NOT NULL
)`

	// SQLCreateTransactionFile records every file gitmap snapshotted
	// before mutating. ON DELETE CASCADE keeps the child rows in sync
	// when PruneOldest drops a parent transaction.
	SQLCreateTransactionFile = `CREATE TABLE IF NOT EXISTS TransactionFile (
	TransactionFileId INTEGER PRIMARY KEY AUTOINCREMENT,
	TransactionId     INTEGER NOT NULL REFERENCES "Transaction"(TransactionId) ON DELETE CASCADE,
	RelPath           TEXT NOT NULL,
	AbsPath           TEXT NOT NULL,
	BackupPath        TEXT NOT NULL,
	ByteSize          INTEGER NOT NULL,
	Sha256            TEXT NOT NULL,
	Action            TEXT NOT NULL
)`

	SQLCreateTransactionStatusIndex = `CREATE INDEX IF NOT EXISTS IdxTransactionStatus
ON "Transaction"(Status, TransactionId DESC)`

	SQLInsertTransaction = `INSERT INTO "Transaction"
(Kind, Status, Argv, Cwd, CreatedAt, ReverseSummary, RepoSlug, GitSha)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	SQLUpdateTransactionStatus = `UPDATE "Transaction"
SET Status = ?, CommittedAt = ?, RevertedAt = ? WHERE TransactionId = ?`

	SQLSelectTransactionByID = `SELECT TransactionId, Kind, Status, Argv, Cwd,
CreatedAt, CommittedAt, RevertedAt, ReverseSummary, RepoSlug, GitSha
FROM "Transaction" WHERE TransactionId = ?`

	SQLSelectTransactionsRecent = `SELECT TransactionId, Kind, Status, Argv, Cwd,
CreatedAt, CommittedAt, RevertedAt, ReverseSummary, RepoSlug, GitSha
FROM "Transaction" ORDER BY TransactionId DESC LIMIT ?`

	SQLSelectLastCommittedTransaction = `SELECT TransactionId FROM "Transaction"
WHERE Status = 'committed' ORDER BY TransactionId DESC LIMIT 1`

	SQLSelectExcessTransactionIDs = `SELECT TransactionId FROM "Transaction"
ORDER BY TransactionId DESC LIMIT -1 OFFSET ?`

	SQLDeleteTransaction = `DELETE FROM "Transaction" WHERE TransactionId = ?`

	SQLInsertTransactionFile = `INSERT INTO TransactionFile
(TransactionId, RelPath, AbsPath, BackupPath, ByteSize, Sha256, Action)
VALUES (?, ?, ?, ?, ?, ?, ?)`

	SQLSelectTransactionFiles = `SELECT TransactionFileId, TransactionId, RelPath,
AbsPath, BackupPath, ByteSize, Sha256, Action
FROM TransactionFile WHERE TransactionId = ? ORDER BY TransactionFileId`
)

// Status / kind / action enums — strings on the wire, constants in code.
const (
	TxnStatusPending   = "pending"
	TxnStatusCommitted = "committed"
	TxnStatusAborted   = "aborted"
	TxnStatusReverted  = "reverted"

	TxnKindClone        = "clone"
	TxnKindMv           = "mv"
	TxnKindMerge        = "merge"
	TxnKindFixRepo      = "fix-repo"
	TxnKindHistoryPurge = "history-purge"
	TxnKindHistoryPin   = "history-pin"

	TxnActionDelete = "delete"
	TxnActionEdit   = "edit"
	TxnActionRename = "rename"

	TxnRetentionCap        = 50
	TxnBackupDirName       = "txn"
	TxnBackupDataDir       = "data"
	TxnBackupFilesDir      = "files"
	TxnUnknownGitShaMarker = "_unknown"

	TxnReplaceModeWarning = "replace mode — file contents NOT recoverable, re-run fix-repo with the prior version"
)

// CLI flags for the extended `gitmap revert` surface.
const (
	FlagRevertListTxn      = "list-txn"
	FlagDescRevertListTxn  = "list the last 50 transactions and exit"
	FlagRevertShowTxn      = "show-txn"
	FlagDescRevertShowTxn  = "print one transaction (including every captured file) by id, then exit"
	FlagRevertTxn          = "txn"
	FlagDescRevertTxn      = "revert the named transaction id"
	FlagRevertLastTxn      = "last-txn"
	FlagDescRevertLastTxn  = "revert the most recent committed transaction"
	FlagRevertLastN        = "last-n-txn"
	FlagDescRevertLastN    = "revert the most recent N committed transactions, newest first"
	FlagRevertPruneTxn     = "prune-txn"
	FlagDescRevertPruneTxn = "force a transaction prune cycle now (drops everything beyond the 50-row cap)"
	FlagRevertForce        = "force"
	FlagDescRevertForce    = "skip the confirm prompt before reverting"
)

// User-facing messages.
const (
	MsgTxnPruned           = "  ✓ pruned %d old transaction(s)\n"
	MsgTxnReverted         = "  ✓ reverted transaction #%d (%s)\n"
	MsgTxnNoCommitted      = "  • no committed transactions found\n"
	MsgTxnConfirmRevert    = "About to revert transaction #%d (%s, %d file(s)).\n"
	MsgTxnConfirmPrompt    = "Type 'yes' to continue: "
	MsgTxnAbortedByUser    = "  • revert canceled by user\n"
	MsgTxnReplaceModeNote  = "  • transaction #%d ran in fix-repo --replace mode and cannot restore file bytes\n"
	MsgTxnLastNHeader      = "About to revert %d transaction(s), newest first:\n"
	MsgTxnLastNRow         = "  #%-4d %-10s %s  %s\n"
	MsgTxnLastNDone        = "  ✓ reverted %d transaction(s)\n"
	MsgTxnLastNNoneFound   = "  • no committed transactions to revert\n"
	ErrRevertLastNBadCount = "revert: --%s requires a positive integer, got %q\n"

	ErrTxnRowNotFound    = "transaction #%d not found\n"
	ErrTxnNotCommitted   = "transaction #%d status is %q, only committed transactions can be reverted\n"
	ErrTxnBackupMissing  = "transaction #%d backup file missing on disk: %q\n"
	ErrTxnBackupShaDrift = "transaction #%d backup sha mismatch for %q (file may have been edited since the backup was taken; pass --force to revert anyway)\n"
	ErrTxnDBOpen         = "transaction journal: failed to open database: %v\n"
	ErrTxnDBWrite        = "transaction journal: write failed: %v\n"
)
