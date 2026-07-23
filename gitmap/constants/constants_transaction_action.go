// Package constants — TransactionAction schema (v23+).
//
// Spec: spec/04-generic-cli/29-reversible-actions-schema.md
//
// TransactionAction is the typed-action layer that sits ON TOP of the
// existing TransactionFile byte-snapshot journal. Each row is a single
// reversible step with a forward+reverse JSON payload and an optional
// BackupRef pointer to a TransactionFile row that holds the byte payload.
//
// Replay rules:
//   - Forward apply: ORDER BY Seq ASC.
//   - Revert: ORDER BY Seq DESC.
//   - Seq is dense, 1-based, unique within a TransactionId.
//   - Reverts are idempotent: applying a reverse to an already-reverted
//     row is a no-op (handler returns nil, RevertedAt is left set).
package constants

const (
	// SQLCreateTransactionAction is the typed action layer.
	// BackupRef is the TransactionFileId of the byte-snapshot row that
	// holds the pre-mutation payload, or 0 when the action does not need
	// one (rename, db row ops, tag flips, workspace writes that the
	// reverse JSON can fully reconstruct).
	SQLCreateTransactionAction = `CREATE TABLE IF NOT EXISTS TransactionAction (
	TransactionActionId INTEGER PRIMARY KEY AUTOINCREMENT,
	TransactionId       INTEGER NOT NULL REFERENCES "Transaction"(TransactionId) ON DELETE CASCADE,
	Seq                 INTEGER NOT NULL,
	Kind                TEXT    NOT NULL,
	ForwardJson         TEXT    NOT NULL,
	ReverseJson         TEXT    NOT NULL,
	BackupRef           INTEGER NOT NULL DEFAULT 0,
	AppliedAt           INTEGER NOT NULL,
	RevertedAt          INTEGER,
	UNIQUE (TransactionId, Seq)
)`

	// SQLCreateTransactionActionTxnSeqIndex is the primary replay index:
	// every revert pass scans (TransactionId ASC|DESC, Seq DESC) so the
	// composite covers both ListByTransaction and the seq-ordered walk.
	SQLCreateTransactionActionTxnSeqIndex = `CREATE INDEX IF NOT EXISTS IdxTransactionActionTxnSeq
ON TransactionAction(TransactionId, Seq DESC)`

	// SQLCreateTransactionActionKindIndex covers analytical queries that
	// filter by Kind across many transactions (e.g. "show every
	// clone_repo action in the last week").
	SQLCreateTransactionActionKindIndex = `CREATE INDEX IF NOT EXISTS IdxTransactionActionKind
ON TransactionAction(Kind, TransactionId DESC)`

	SQLInsertTransactionAction = `INSERT INTO TransactionAction
(TransactionId, Seq, Kind, ForwardJson, ReverseJson, BackupRef, AppliedAt)
VALUES (?, ?, ?, ?, ?, ?, ?)`

	SQLSelectTransactionActions = `SELECT TransactionActionId, TransactionId, Seq,
Kind, ForwardJson, ReverseJson, BackupRef, AppliedAt, RevertedAt
FROM TransactionAction WHERE TransactionId = ? ORDER BY Seq ASC`

	SQLSelectTransactionActionsReverse = `SELECT TransactionActionId, TransactionId, Seq,
Kind, ForwardJson, ReverseJson, BackupRef, AppliedAt, RevertedAt
FROM TransactionAction WHERE TransactionId = ? ORDER BY Seq DESC`

	SQLSelectMaxActionSeq = `SELECT COALESCE(MAX(Seq), 0)
FROM TransactionAction WHERE TransactionId = ?`

	SQLMarkTransactionActionReverted = `UPDATE TransactionAction
SET RevertedAt = ? WHERE TransactionActionId = ?`
)

// ActionKind enum — the closed set of typed reversible actions.
// Every gitmap mutating command must decompose into a sequence of these.
const (
	ActionKindCreatePath     = "create_path"
	ActionKindDeletePath     = "delete_path"
	ActionKindEditFile       = "edit_file"
	ActionKindRenamePath     = "rename_path"
	ActionKindCloneRepo      = "clone_repo"
	ActionKindDBUpsertRepo   = "db_upsert_repo"
	ActionKindDBDeleteRepo   = "db_delete_repo"
	ActionKindTagSet         = "tag_set"
	ActionKindWorkspaceWrite = "workspace_write"
)

// Action error messages — typed for errors.Is dispatch by callers.
const (
	ErrActionUnknownKind   = "transaction action #%d: unknown kind %q"
	ErrActionSeqGap        = "transaction action: seq gap detected for txn #%d (expected %d, got %d)"
	ErrActionInsert        = "transaction action insert: %w"
	ErrActionList          = "transaction action list: %w"
	ErrActionMarkReverted  = "transaction action mark reverted: %w"
	ErrActionPayloadDecode = "transaction action #%d: payload decode failed: %w"
	ErrActionLiveConflict  = "transaction action #%d (%s): live state diverged from forward payload"
)
