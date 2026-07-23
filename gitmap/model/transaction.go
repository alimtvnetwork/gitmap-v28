// Package model — transaction.go defines the SQLite-backed transaction
// journal records used by the revert/undo subsystem. See
// spec/04-generic-cli/28-transaction-revert.md for the wire contract.
package model

// TransactionRecord is one row of the "Transaction" master journal.
//
// One row is written per state-mutating gitmap command (clone, mv,
// merge-*, fix-repo, history-*, ...). Status flows
// pending → committed (or pending → aborted on failure;
// committed → reverted by `gitmap revert --txn <id>`).
type TransactionRecord struct {
	ID             int64  `json:"id"`
	Kind           string `json:"kind"`
	Status         string `json:"status"`
	Argv           string `json:"argv"`
	Cwd            string `json:"cwd"`
	CreatedAt      int64  `json:"createdAt"`
	CommittedAt    int64  `json:"committedAt,omitempty"`
	RevertedAt     int64  `json:"revertedAt,omitempty"`
	ReverseSummary string `json:"reverseSummary"`
	RepoSlug       string `json:"repoSlug,omitempty"`
	GitSha         string `json:"gitSha,omitempty"`
}

// TransactionFileRecord captures one file gitmap snapshotted to the
// per-transaction backup directory before mutating it on disk.
//
// Action is one of TxnAction* constants ("delete", "edit", "rename").
type TransactionFileRecord struct {
	ID            int64  `json:"id"`
	TransactionID int64  `json:"transactionId"`
	RelPath       string `json:"relPath"`
	AbsPath       string `json:"absPath"`
	BackupPath    string `json:"backupPath"`
	ByteSize      int64  `json:"byteSize"`
	Sha256        string `json:"sha256"`
	Action        string `json:"action"`
}
