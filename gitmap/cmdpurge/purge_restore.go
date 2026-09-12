package cmdpurge

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func fetchActivePurgeLog(db *store.DB, repoPath string) (*store.PurgeHistoryLog, error) {
	log, err := db.GetLastPurgeHistoryLog(repoPath)
	if err != nil {
		return nil, apperror.Wrap(err, "failed to get purge log", nil)
	}

	if log == nil {
		return nil, apperror.NewSimple("EXECUTION", "No active purge state found to restore.")
	}

	return log, nil
}

func resetToBranch(branch string) error {
	_, err := runPurgeCmd("git", "reset", "--hard", branch)
	if err != nil {
		return apperror.Wrap(err, "failed to reset to backup branch", nil)
	}

	return nil
}

func restoreFileEntry(tempDir, path string, d fs.DirEntry) error {
	if d.IsDir() {
		return nil
	}

	rel, err := filepath.Rel(tempDir, path)
	if err != nil {
		return apperror.Wrap(err, "failed to get relative path", nil)
	}

	if err := os.MkdirAll(filepath.Dir(rel), 0755); err != nil {
		return apperror.Wrap(err, "failed to create restore dir", nil)
	}

	return copyPurgeFile(path, rel)
}

func restoreFilesFromTemp(tempDir string) error {
	if tempDir == "" {
		return nil
	}

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(tempDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		return restoreFileEntry(tempDir, p, d)
	})
}

func markPurgeRestored(db *store.DB, id int64) error {
	if err := db.MarkPurgeHistoryRestored(id); err != nil {
		return apperror.Wrap(err, "failed to mark purge log restored", nil)
	}

	return nil
}

func doRestore(db *store.DB, repoPath string) error {
	log, err := fetchActivePurgeLog(db, repoPath)
	if err != nil {
		return err
	}

	if err := resetToBranch(log.BackupBranch); err != nil {
		return err
	}

	if err := restoreFilesFromTemp(log.TempDir); err != nil {
		return err
	}

	return markPurgeRestored(db, log.PurgeHistoryLogId)
}
