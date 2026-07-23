package txn

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// RevertOptions tunes Revert behavior.
type RevertOptions struct {
	Force bool // skip pre-revert sha256 verification of backup blobs
}

// Revert applies the inverse of one committed transaction:
//   - delete  → restore file bytes from backup
//   - edit    → overwrite current bytes with backup
//   - rename  → swap the path back (no bytes copied)
func Revert(db *store.DB, id int64, opts RevertOptions) error {
	row, err := db.FindTransactionByID(id)
	if err != nil {
		return wrapNotFound(id, err)
	}
	if err := assertCommitted(row); err != nil {
		return err
	}
	files, err := db.ListTransactionFiles(id)
	if err != nil {
		return err
	}
	if err := applyAllReverses(files, opts); err != nil {
		return err
	}

	return db.MarkTransactionReverted(id)
}

// wrapNotFound formats the user-facing not-found error.
func wrapNotFound(id int64, err error) error {
	if errors.Is(err, store.ErrTransactionNotFound) {
		return fmt.Errorf(constants.ErrTxnRowNotFound, id)
	}

	return err
}

// assertCommitted refuses to revert anything that isn't in committed state.
func assertCommitted(r model.TransactionRecord) error {
	if r.Status == constants.TxnStatusCommitted {
		return nil
	}

	return fmt.Errorf(constants.ErrTxnNotCommitted, r.ID, r.Status)
}

// applyAllReverses walks the file list newest-first to undo nested mutations.
func applyAllReverses(files []model.TransactionFileRecord, opts RevertOptions) error {
	for i := len(files) - 1; i >= 0; i-- {
		if err := applyReverse(files[i], opts); err != nil {
			return err
		}
	}

	return nil
}

// applyReverse dispatches by Action to the right inverse-op helper.
func applyReverse(f model.TransactionFileRecord, opts RevertOptions) error {
	switch f.Action {
	case constants.TxnActionRename:
		return reverseRename(f)
	case constants.TxnActionDelete, constants.TxnActionEdit:
		return reverseRestore(f, opts)
	default:
		return fmt.Errorf("transaction revert: unknown action %q", f.Action)
	}
}

// reverseRename moves the post-mutation path back to its original location.
func reverseRename(f model.TransactionFileRecord) error {
	if err := os.MkdirAll(filepath.Dir(f.BackupPath), 0o755); err != nil {
		return fmt.Errorf("transaction revert mkdir: %w", err)
	}
	if err := os.Rename(f.AbsPath, f.BackupPath); err != nil {
		return fmt.Errorf("transaction revert rename %q→%q: %w",
			f.AbsPath, f.BackupPath, err)
	}

	return nil
}

// reverseRestore copies BackupPath bytes back over AbsPath.
func reverseRestore(f model.TransactionFileRecord, opts RevertOptions) error {
	if err := assertBackupSha(f, opts); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(f.AbsPath), 0o755); err != nil {
		return fmt.Errorf("transaction revert mkdir: %w", err)
	}

	return copyFile(f.BackupPath, f.AbsPath)
}

// assertBackupSha re-hashes the backup blob and compares to the journal sha.
func assertBackupSha(f model.TransactionFileRecord, opts RevertOptions) error {
	if opts.Force || len(f.Sha256) == 0 {
		return nil
	}
	got, err := hashFile(f.BackupPath)
	if err != nil {
		return fmt.Errorf(constants.ErrTxnBackupMissing, f.TransactionID, f.BackupPath)
	}
	if got != f.Sha256 {
		return fmt.Errorf(constants.ErrTxnBackupShaDrift, f.TransactionID, f.BackupPath)
	}

	return nil
}

// hashFile streams sha256 over a backup blob.
func hashFile(path string) (string, error) {
	in, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer in.Close()
	h := sha256.New()
	if _, err := io.Copy(h, in); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// copyFile is a small helper used by the restore path.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("transaction revert open: %w", err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("transaction revert create: %w", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("transaction revert copy: %w", err)
	}

	return nil
}
