// Package txn is the SQLite-backed transaction journal that gives every
// state-mutating gitmap command a recorded, revertable trail.
//
// Spec: spec/04-generic-cli/28-transaction-revert.md
//
// Lifecycle:
//
//	t, _ := txn.Begin(db, txn.Meta{Kind: TxnKindMv, Argv: os.Args, ...})
//	t.RecordRename(absFrom, absTo)            // bytes-free rename
//	t.SnapshotEdit(absPath)                    // before mutating in place
//	t.SnapshotDelete(absPath)                  // before unlinking
//	if err != nil { t.Abort(); return err }
//	t.Commit()
//
// Backups land at:
//
//	<binaryDir>/.gitmap/txn/<txnId>/data/<repoSlug>/<gitSha>/files/<rel>
package txn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// Meta is the per-command context every transaction needs at Begin time.
type Meta struct {
	Kind           string   // TxnKind* constant
	Argv           []string // os.Args verbatim, joined for the journal row
	Cwd            string   // process working directory at command start
	ReverseSummary string   // one-line "what revert will do"
	RepoSlug       string   // empty when not repo-scoped
	GitSha         string   // resolved HEAD; "" → TxnUnknownGitShaMarker
}

// Journal is a live transaction handle.
type Journal struct {
	db        *store.DB
	id        int64
	binaryDir string
	meta      Meta
}

// Begin inserts a pending row and returns a handle the caller commits or aborts.
// The returned *Journal is never nil even when SQLite errors — callers can
// safely defer Abort() and unconditionally call snapshot helpers, which
// degrade to no-ops when id == 0.
func Begin(db *store.DB, m Meta) (*Journal, error) {
	rec := model.TransactionRecord{
		Kind:           m.Kind,
		Argv:           strings.Join(m.Argv, " "),
		Cwd:            m.Cwd,
		ReverseSummary: m.ReverseSummary,
		RepoSlug:       m.RepoSlug,
		GitSha:         resolveGitSha(m.GitSha),
	}
	id, err := db.InsertTransaction(rec)
	if err != nil {
		return &Journal{db: db, meta: m}, err
	}

	return &Journal{db: db, id: id, binaryDir: store.BinaryDataDir(), meta: m}, nil
}

// ID returns the journal row id (0 when Begin's INSERT failed).
func (j *Journal) ID() int64 { return j.id }

// Commit flips the row to committed and prunes anything beyond the cap.
func (j *Journal) Commit() error {
	if j.id == 0 {
		return nil
	}
	if err := j.db.MarkTransactionCommitted(j.id); err != nil {
		return err
	}

	return j.pruneExcess()
}

// Abort flips the row to aborted and removes any partial backup dir.
func (j *Journal) Abort() error {
	if j.id == 0 {
		return nil
	}
	_ = os.RemoveAll(j.txnRoot())

	return j.db.MarkTransactionAborted(j.id)
}

// pruneExcess keeps only the newest TxnRetentionCap rows + their backups.
func (j *Journal) pruneExcess() error {
	dropped, err := j.db.PruneOldestTransactions(constants.TxnRetentionCap)
	if err != nil {
		return err
	}
	for _, id := range dropped {
		_ = os.RemoveAll(j.txnRootFor(id))
	}

	return nil
}

// resolveGitSha defaults blank shas to the unknown-marker constant.
func resolveGitSha(s string) string {
	if len(s) == 0 {
		return constants.TxnUnknownGitShaMarker
	}

	return s
}

// txnRoot is the on-disk backup root for this live transaction.
func (j *Journal) txnRoot() string { return j.txnRootFor(j.id) }

// txnRootFor builds the backup root for an arbitrary transaction id.
func (j *Journal) txnRootFor(id int64) string {
	return filepath.Join(filepath.Dir(j.binaryDir),
		constants.GitMapSubdir, constants.TxnBackupDirName,
		fmt.Sprintf("%d", id))
}
