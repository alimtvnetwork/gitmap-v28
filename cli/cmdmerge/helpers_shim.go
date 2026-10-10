package cmdmerge

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type moveOpts struct {
	yes           bool
	dryRun        bool
	isSkipVSCode  bool
	isSkipDesktop bool
}

// ResolveRepo finds a single repo matching target string across DB and JSON records.
func ResolveRepo(db *store.DB, target string) (*model.ScanRecord, error) {
	all := loadUnifiedCandidates(db)
	t := strings.TrimSpace(target)
	if len(t) == 0 || t == "." {
		return resolveByPWD(all)
	}

	if rec := resolveByPath(t, all); rec != nil {
		return rec, nil
	}

	if rec := resolveByAlias(db, t, all); rec != nil {
		return rec, nil
	}

	if rec := resolveByTarget(t, all); rec != nil {
		return rec, nil
	}

	return nil, apperror.NewNotFoundError(fmt.Sprintf("no repository matched %q", target))
}

func calculateDestPath(srcPath, destTarget string) (string, error) {
	cleanSrc := filepath.Clean(filepath.FromSlash(fsutil.TrimTrailingSlashes(srcPath)))
	repoName := filepath.Base(cleanSrc)
	t := fsutil.TrimTrailingSlashes(destTarget)

	if t == ".." {
		parent := filepath.Dir(cleanSrc)
		grandParent := filepath.Dir(parent)

		return filepath.Join(grandParent, repoName), nil
	}

	absDest, err := filepath.Abs(t)
	if err != nil {
		return "", err
	}

	if fsutil.DirExists(absDest) {
		return filepath.Join(absDest, repoName), nil
	}

	return absDest, nil
}

func openDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return store.OpenGlobalDefault()
	}
	_ = db.Migrate()
	if countRegisteredSSHHosts(db) == 0 {
		return resolveGlobalSSHDBFallback(db), nil
	}
	return db, nil
}

func parseMoveFlags(args []string) (moveOpts, []string) {
	opts := moveOpts{}
	var positional []string
	for _, a := range args {
		switch a {
		case "-y", "--yes":
			opts.yes = true
		case "--dry-run":
			opts.dryRun = true
		case "--no-vscode":
			opts.isSkipVSCode = true
		case "--no-desktop":
			opts.isSkipDesktop = true
		default:
			positional = append(positional, a)
		}
	}

	return opts, positional
}

func preflightMove(srcPath, destPath string) error {
	if !fsutil.DirExists(srcPath) {
		return fmt.Errorf("source directory %q does not exist", srcPath)
	}

	if fsutil.EqualPaths(srcPath, destPath) {
		return fmt.Errorf("source and destination are identical")
	}

	if fsutil.FileOrDirExists(destPath) {
		return fmt.Errorf("destination %q already exists", destPath)
	}

	if fsutil.IsSubdirectory(srcPath, destPath) {
		return fmt.Errorf("cannot move repository into its own subdirectory")
	}

	return nil
}

func printMoveDryRun(src, dest string) {
	fmt.Println("  (dry-run) Move operation preview:")
	fmt.Printf("     Source:      %s\n", src)
	fmt.Printf("     Destination: %s\n", dest)
	fmt.Println("     Actions:     Relocate directory, update SQLite Repo path, sync VS Code & GitHub Desktop")
}

func resolveInspectionFiles(fileArgs []string, targetDir string) ([]string, error) {
	if len(fileArgs) > 0 {
		return expandFileGlobs(fileArgs), nil
	}
	return scanDirectoryJSONFiles(targetDir)
}

func countRegisteredSSHHosts(db *store.DB) int {
	var count int
	row := db.SQL().QueryRow("SELECT count(*) FROM ssh_hosts")
	if scanErr := row.Scan(&count); scanErr != nil {
		return 0
	}
	return count
}
func resolveGlobalSSHDBFallback(localDB *store.DB) *store.DB {
	globalDB, err := store.OpenGlobalDefault()
	if err != nil {
		return localDB
	}
	_ = globalDB.Migrate()
	if countRegisteredSSHHosts(globalDB) > 0 {
		localDB.Close()
		return globalDB
	}
	globalDB.Close()
	return localDB
}

func loadUnifiedCandidates(db *store.DB) []model.ScanRecord {
	dbRepos := loadDbCandidates(db)
	jsonRepos := loadJSONCandidates()

	return mergeCandidateRepos(dbRepos, jsonRepos)
}

func loadDbCandidates(db *store.DB) []model.ScanRecord {
	if db == nil {
		return nil
	}

	repos, err := db.ListRepos()
	if err != nil {
		return nil
	}

	return repos
}

func loadJSONCandidates() []model.ScanRecord {
	paths := []string{
		filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile),
		filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
	}
	for _, p := range paths {
		if records, isOk := tryLoadJSONPath(p); isOk {
			return records
		}
	}

	return nil
}

func tryLoadJSONPath(path string) ([]model.ScanRecord, bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, false
	}

	records, err := model.LoadStatusRecords(path)
	if err != nil || len(records) == 0 {
		return nil, false
	}

	return records, true
}

func mergeCandidateRepos(primary, secondary []model.ScanRecord) []model.ScanRecord {
	seen := make(map[string]bool)
	var out []model.ScanRecord

	for _, r := range primary {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	for _, r := range secondary {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	return out
}

func appendUniqueRecords(out, hits []model.ScanRecord, seen map[string]bool) []model.ScanRecord {
	for _, r := range hits {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	return out
}

func resolveByPWD(all []model.ScanRecord) (*model.ScanRecord, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get current directory: %w", err)
	}

	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, pwd) {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("current directory (%s) is not a tracked repository", pwd)
}

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
	wrapper, appErr := dbengine.WrapDb(db.Conn(), dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	cleanPath := store.NormalizeStoragePath(newPath)
	if txErr := runUpdateMoveTx(ctx, wrapper, repoID, cleanPath, newName); txErr != nil {
		return txErr
	}

	return nil
}

func runUpdateMoveTx(ctx context.Context, wrapper *dbengine.DbWrapper, repoID int64, newPath, newName string) *apperror.AppError {
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

func scanForColumnName(rows *sql.Rows, targetCol string) bool {
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt any
		scanErr := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		isMatch := scanErr == nil && strings.EqualFold(name, targetCol)
		if isMatch {
			return true
		}
	}

	return false
}

func expandFileGlobs(patterns []string) []string {
	var collected []string
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err == nil && len(matches) > 0 {
			collected = append(collected, matches...)
			continue
		}
		collected = append(collected, p)
	}
	return dedupeAndFilterJSONFiles(collected)
}

func scanDirectoryJSONFiles(dir string) ([]string, error) {
	var matches []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), ".json") {
			matches = append(matches, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(matches)
	return matches, nil
}

func dedupeAndFilterJSONFiles(files []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, f := range files {
		clean := filepath.Clean(f)
		if seen[clean] {
			continue
		}
		seen[clean] = true
		out = append(out, clean)
	}
	return out
}
