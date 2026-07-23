// Package txn — typed Reversible-Action layer (spec/04-generic-cli/29).
//
// This sits ON TOP of the existing TransactionFile byte-snapshot journal:
//
//   - RecordEditFile  → snapshots bytes (TransactionFile row) AND writes
//     an edit_file action row with a BackupRef pointer.
//   - RecordRenamePath → writes a rename_path action row only (no bytes).
//
// On revert (RevertActions), rows are walked Seq DESC and dispatched
// per-Kind. Handlers are pure: they operate on the JSON payload + the
// referenced TransactionFile row, and they are idempotent — replaying a
// reverse on already-reverted state returns nil.
package txn

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// EditFilePayload describes the forward side of edit_file. The reverse
// side is the same struct with BackupRef populated by the journal.
type EditFilePayload struct {
	AbsPath string `json:"absPath"`
	Sha256  string `json:"sha256,omitempty"`
}

// RenamePathPayload describes both sides of rename_path; reverse swaps
// From and To.
type RenamePathPayload struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// RecordEditFile snapshots absPath bytes (reusing SnapshotEdit) AND
// writes an edit_file action row pointing at the byte-backup row.
func (j *Journal) RecordEditFile(absPath string) error {
	if j.id == 0 {
		return nil
	}
	if err := j.SnapshotEdit(absPath); err != nil {
		return err
	}
	backupRef, err := j.lastFileRowID()
	if err != nil {
		return err
	}
	fwd, _ := json.Marshal(EditFilePayload{AbsPath: absPath})
	rev, _ := json.Marshal(EditFilePayload{AbsPath: absPath})

	return j.appendAction(constants.ActionKindEditFile, string(fwd), string(rev), backupRef)
}

// RecordRenamePath logs from→to (no bytes). Reverse swaps the pair.
func (j *Journal) RecordRenamePath(from, to string) error {
	if j.id == 0 {
		return nil
	}
	fwd, _ := json.Marshal(RenamePathPayload{From: from, To: to})
	rev, _ := json.Marshal(RenamePathPayload{From: to, To: from})

	return j.appendAction(constants.ActionKindRenamePath, string(fwd), string(rev), 0)
}

// appendAction is the shared seq-allocate + insert path.
func (j *Journal) appendAction(kind, fwd, rev string, backupRef int64) error {
	seq, err := j.db.NextActionSeq(j.id)
	if err != nil {
		return err
	}
	_, err = j.db.InsertTransactionAction(model.TransactionActionRecord{
		TransactionID: j.id,
		Seq:           seq,
		Kind:          kind,
		ForwardJSON:   fwd,
		ReverseJSON:   rev,
		BackupRef:     backupRef,
	})

	return err
}

// lastFileRowID returns the most recent TransactionFile row id for this
// transaction; used to wire BackupRef immediately after SnapshotEdit.
func (j *Journal) lastFileRowID() (int64, error) {
	files, err := j.db.ListTransactionFiles(j.id)
	if err != nil {
		return 0, err
	}
	if len(files) == 0 {
		return 0, nil
	}

	return files[len(files)-1].ID, nil
}

// RevertActions replays the typed-action chain in Seq DESC order.
// This is the public entry point used by the v23+ revert path; it is a
// peer to (not a replacement for) the legacy file-only Revert.
func RevertActions(db *store.DB, txnID int64) error {
	row, err := db.FindTransactionByID(txnID)
	if err != nil {
		return wrapNotFound(txnID, err)
	}
	if err := assertCommitted(row); err != nil {
		return err
	}
	actions, err := db.ListTransactionActionsReverse(txnID)
	if err != nil {
		return err
	}

	return applyAllActionReverses(db, actions)
}

// applyAllActionReverses dispatches each action row to its handler.
func applyAllActionReverses(db *store.DB, actions []model.TransactionActionRecord) error {
	for _, a := range actions {
		if a.RevertedAt != 0 {
			continue // idempotent: already reverted
		}
		if err := dispatchActionReverse(db, a); err != nil {
			return err
		}
		if err := db.MarkTransactionActionReverted(a.ID); err != nil {
			return err
		}
	}

	return nil
}

// dispatchActionReverse routes one action to its per-kind handler.
func dispatchActionReverse(db *store.DB, a model.TransactionActionRecord) error {
	switch a.Kind {
	case constants.ActionKindEditFile:
		return reverseEditFileAction(db, a)
	case constants.ActionKindRenamePath:
		return reverseRenamePathAction(a)
	default:
		return fmt.Errorf(constants.ErrActionUnknownKind, a.ID, a.Kind)
	}
}

// reverseEditFileAction restores bytes from the linked TransactionFile row.
func reverseEditFileAction(db *store.DB, a model.TransactionActionRecord) error {
	files, err := db.ListTransactionFiles(a.TransactionID)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.ID == a.BackupRef {
			return reverseRestore(f, RevertOptions{})
		}
	}

	return fmt.Errorf(constants.ErrTxnBackupMissing, a.TransactionID, "<backup row not found>")
}

// reverseRenamePathAction swaps the path back using the reverse payload.
func reverseRenamePathAction(a model.TransactionActionRecord) error {
	var p RenamePathPayload
	if err := json.Unmarshal([]byte(a.ReverseJSON), &p); err != nil {
		return fmt.Errorf(constants.ErrActionPayloadDecode, a.ID, err)
	}
	if _, err := os.Stat(p.From); err != nil {
		return fmt.Errorf(constants.ErrActionLiveConflict, a.ID, a.Kind)
	}

	return os.Rename(p.From, p.To)
}
