package store

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ClearPullErrors removes pull error records for a specific repo slug or all repos if repoSlug is "all" or empty.
func (s *PullSplitDB) ClearPullErrors(repoSlug string) *apperror.AppError {
	if s.conn == nil {
		return apperror.NewValidationError("database connection is nil")
	}

	if isAllReposQuery(repoSlug) {
		return s.clearAllPullErrors()
	}

	return s.clearPullErrorsBySlug(repoSlug)
}

func (s *PullSplitDB) clearAllPullErrors() *apperror.AppError {
	_, err := s.conn.Exec(`DELETE FROM pull_errors`)
	if err != nil {
		return apperror.WrapSimple(err, "pull_errors.clear_all")
	}

	return nil
}

func (s *PullSplitDB) clearPullErrorsBySlug(repoSlug string) *apperror.AppError {
	query := `DELETE FROM pull_errors WHERE repo_slug = ? OR repo_path LIKE ?`
	_, err := s.conn.Exec(query, repoSlug, "%"+repoSlug+"%")
	if err != nil {
		return apperror.WrapSimple(err, "pull_errors.clear_repo")
	}

	return nil
}

// ClearPullErrorsForRepo removes pull errors for a repo and deletes the per-repo details DB if present.
func (s *PullSplitDB) ClearPullErrorsForRepo(repoSlug, repoPath string) *apperror.AppError {
	if s.conn == nil {
		return apperror.NewValidationError("database connection is nil")
	}

	if err := s.ClearPullErrors(repoSlug); err != nil {
		return err
	}

	removeRepoPullErrorDB(repoPath)

	return nil
}

func removeRepoPullErrorDB(repoPath string) {
	if repoPath == "" {
		return
	}

	repoDBPath := filepath.Join(repoPath, ".gitmap", "pull_errors.db")
	_ = os.Remove(repoDBPath)
}
