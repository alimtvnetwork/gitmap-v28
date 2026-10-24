package cmdmerge

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func resolveByPath(target string, all []model.ScanRecord) *model.ScanRecord {
	cleanTarget := fsutil.NormalizeSlashes(fsutil.TrimTrailingSlashes(target))
	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, cleanTarget) {
			return &r
		}
	}

	if found := findByCleanTargetAbs(cleanTarget, all); found != nil {
		return found
	}

	for _, r := range all {
		rClean := fsutil.NormalizeSlashes(r.AbsolutePath)
		if strings.EqualFold(filepath.Base(rClean), filepath.Base(cleanTarget)) || strings.EqualFold(r.Slug, filepath.Base(cleanTarget)) {
			return &r
		}
	}

	return nil
}

func findByCleanTargetAbs(cleanTarget string, all []model.ScanRecord) *model.ScanRecord {
	abs, err := filepath.Abs(cleanTarget)
	if err != nil {
		return nil
	}

	return findByAbsolutePath(abs, all)
}

func findByAbsolutePath(abs string, all []model.ScanRecord) *model.ScanRecord {
	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, abs) {
			return &r
		}
	}

	return nil
}

func resolveByAlias(db *store.DB, target string, all []model.ScanRecord) *model.ScanRecord {
	aliasRow, err := db.ResolveAlias(target)
	if err != nil {
		return nil
	}

	for _, r := range all {
		if r.ID == aliasRow.RepoID || fsutil.EqualPaths(r.AbsolutePath, aliasRow.AbsolutePath) {
			return &r
		}
	}

	return nil
}

func resolveByTarget(target string, all []model.ScanRecord) *model.ScanRecord {
	for _, r := range all {
		if matchesRepoTarget(r, target) {
			return &r
		}
	}

	return nil
}

func matchesRepoTarget(r model.ScanRecord, target string) bool {
	t := strings.TrimSpace(target)
	if len(t) == 0 {
		return false
	}

	if matchesRepoNameOrSlug(r, t) {
		return true
	}

	return matchesRepoPath(r, t)
}

func matchesRepoNameOrSlug(r model.ScanRecord, t string) bool {
	if strings.EqualFold(r.Slug, t) || strings.EqualFold(r.RepoName, t) {
		return true
	}

	slugBase := filepath.Base(r.Slug)
	repoBase := filepath.Base(r.RepoName)

	return strings.EqualFold(slugBase, t) || strings.EqualFold(repoBase, t)
}

func matchesRepoPath(r model.ScanRecord, t string) bool {
	if fsutil.EqualPaths(r.AbsolutePath, t) {
		return true
	}

	dirBase := filepath.Base(fsutil.NormalizeSlashes(r.AbsolutePath))
	if strings.EqualFold(dirBase, t) {
		return true
	}

	abs, err := filepath.Abs(t)

	return err == nil && fsutil.EqualPaths(r.AbsolutePath, abs)
}

// PrintRepoSuggestions queries the database for suggestions and prints them.
func PrintRepoSuggestions(db *store.DB, target string) {
	if db == nil {
		return
	}

	suggs, _ := db.GetRepoSuggestions(target)
	if len(suggs) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "Did you mean:\n")
	for _, s := range suggs {
		fmt.Fprintf(os.Stderr, "  %s\n", s)
	}
}

func confirmMovePrompt(slug, src, dest string) bool {
	fmt.Printf("Move %s\n  From: %s\n  To:   %s\nProceed? [y/N] ", slug, src, dest)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	ans := strings.ToLower(strings.TrimSpace(line))

	return ans == "y" || ans == "yes"
}

func executeMove(db *store.DB, rec model.ScanRecord, destPath string, opts moveOpts) {
	if err := fsutil.SafeRename(rec.AbsolutePath, destPath); err != nil {
		fmt.Printf("mv: physical move failed: %v\n", err)

		return
	}

	newRepoName := filepath.Base(destPath)
	if err := updateRepoInDB(db, rec.ID, destPath, newRepoName); err != nil {
		fmt.Printf("mv: database update error: %v\n", err)
	}

	syncExternalMove(rec.AbsolutePath, destPath, newRepoName, opts)
	printMoveSuccess(rec.Slug, rec.AbsolutePath, destPath)
}

func syncExternalMove(oldPath, newPath, newName string, opts moveOpts) {
	if !opts.isSkipVSCode {
		_ = vscodepm.UpdateRootPath(oldPath, newPath, newName)
	}

	if !opts.isSkipDesktop {
		_ = desktop.UpdateRepoPath(oldPath, newPath)
	}
}

func printMoveSuccess(slug, oldPath, newPath string) {
	fmt.Println()
	fmt.Printf("  ✔ Moved %q successfully:\n", slug)
	fmt.Printf("     From: %s\n", oldPath)
	fmt.Printf("     To:   %s\n", newPath)
	fmt.Println()
}

func updateRepoInDB(db *store.DB, repoID int64, newPath, newName string) error {
	ctx := context.Background()
	wrapper, appErr := dbengine.WrapDB(db.Conn(), dbengine.DBSQLite)
	if appErr != nil {
		return appErr
	}

	cleanPath := store.NormalizeStoragePath(newPath)
	if txErr := runUpdateMoveTx(ctx, wrapper, repoID, cleanPath, newName); txErr != nil {
		return txErr
	}

	return nil
}

func runUpdateMoveTx(ctx context.Context, wrapper *dbengine.DBWrapper, repoID int64, newPath, newName string) *apperror.AppError {
	return wrapper.WithImmediateTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeMoveDBTx(ctx, tx, repoID, newPath, newName)
	})
}

func executeMoveDBTx(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, newPath, newName string) *apperror.AppError {
	if err := updateRepoRow(ctx, tx, repoID, newPath, newName); err != nil {
		return err
	}

	return updateAliasRowIfPresent(ctx, tx, repoID, newPath)
}

func updateRepoRow(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, newPath, newName string) *apperror.AppError {
	query := "UPDATE Repo SET AbsolutePath = ?, RepoName = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE RepoId = ?"
	if _, err := tx.Exec(ctx, query, newPath, newName, repoID); err == nil {
		return nil
	}

	fallbackQuery := "UPDATE Repo SET AbsolutePath = ?, RepoName = ? WHERE Id = ?"
	if _, fallbackErr := tx.Exec(ctx, fallbackQuery, newPath, newName, repoID); fallbackErr != nil {
		return apperror.WrapSimple(fallbackErr, "mv: update repo row")
	}

	return nil
}

func updateAliasRowIfPresent(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, newPath string) *apperror.AppError {
	hasColumn := hasAliasPathColumn(ctx, tx)
	if !hasColumn {
		return nil
	}

	query := "UPDATE Alias SET AbsolutePath = ? WHERE RepoId = ?"
	if _, err := tx.Exec(ctx, query, newPath, repoID); err != nil {
		return apperror.WrapSimple(err, "mv: update alias row")
	}

	return nil
}

func hasAliasPathColumn(ctx context.Context, tx *dbengine.TxWrapper) bool {
	rows, appErr := tx.Query(ctx, "PRAGMA table_info(Alias)")
	if appErr != nil {
		return false
	}

	defer rows.Close()

	return scanForColumnName(rows, "AbsolutePath")
}
