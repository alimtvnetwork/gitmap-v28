package txn

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// SnapshotEdit copies the current bytes of absPath into the backup tree
// and records an "edit" row. Call BEFORE the in-place mutation.
func (j *Journal) SnapshotEdit(absPath string) error {
	return j.snapshotFile(absPath, constants.TxnActionEdit)
}

// SnapshotDelete copies the file aside and records a "delete" row.
// Call BEFORE removing absPath from disk.
func (j *Journal) SnapshotDelete(absPath string) error {
	return j.snapshotFile(absPath, constants.TxnActionDelete)
}

// RecordRename logs a directory- or file-rename so revert can swap it back.
// Bytes are NOT copied — the rename is the inverse of a rename.
func (j *Journal) RecordRename(from, to string) error {
	if j.id == 0 {
		return nil
	}
	rec := model.TransactionFileRecord{
		TransactionID: j.id,
		RelPath:       j.relTo(from),
		AbsPath:       to,   // post-mutation location
		BackupPath:    from, // pre-mutation location (the inverse)
		Action:        constants.TxnActionRename,
	}

	return j.db.InsertTransactionFile(rec)
}

// snapshotFile is the shared copy-then-record path for edit + delete.
func (j *Journal) snapshotFile(absPath, action string) error {
	if j.id == 0 {
		return nil
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("transaction snapshot stat %q: %w", absPath, err)
	}
	backup := j.backupPath(absPath)
	sum, err := copyFileWithSha(absPath, backup)
	if err != nil {
		return err
	}

	return j.recordSnapshot(absPath, backup, action, info.Size(), sum)
}

// recordSnapshot writes the TransactionFile row for one snapshotted file.
func (j *Journal) recordSnapshot(abs, backup, action string, size int64, sum string) error {
	return j.db.InsertTransactionFile(model.TransactionFileRecord{
		TransactionID: j.id,
		RelPath:       j.relTo(abs),
		AbsPath:       abs,
		BackupPath:    backup,
		ByteSize:      size,
		Sha256:        sum,
		Action:        action,
	})
}

// backupPath maps an absolute source path into the per-txn backup tree.
func (j *Journal) backupPath(abs string) string {
	rel := j.relTo(abs)
	rel = strings.TrimPrefix(rel, string(os.PathSeparator))

	return filepath.Join(j.txnRoot(), constants.TxnBackupDataDir,
		safeSlug(j.meta.RepoSlug), j.meta.GitSha,
		constants.TxnBackupFilesDir, rel)
}

// relTo returns abs relative to the meta cwd, falling back to the basename.
func (j *Journal) relTo(abs string) string {
	if len(j.meta.Cwd) == 0 {
		return filepath.Base(abs)
	}
	rel, err := filepath.Rel(j.meta.Cwd, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return filepath.Base(abs)
	}

	return rel
}

// safeSlug ensures repo-slugs that are empty don't collapse the path.
func safeSlug(s string) string {
	if len(s) == 0 {
		return "_unscoped"
	}

	return s
}

// copyFileWithSha streams src → dst (creating parents) and returns its sha256.
func copyFileWithSha(src, dst string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", fmt.Errorf("transaction backup mkdir: %w", err)
	}
	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("transaction backup open src: %w", err)
	}
	defer in.Close()

	return writeBackupAndHash(in, dst)
}

// writeBackupAndHash writes the backup file and returns the streamed sha256.
func writeBackupAndHash(in io.Reader, dst string) (string, error) {
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("transaction backup create: %w", err)
	}
	defer out.Close()
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(out, h), in); err != nil {
		return "", fmt.Errorf("transaction backup copy: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
