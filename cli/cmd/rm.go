package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// rmUsage describes the `gitmap rm` command.
const rmUsage = `Usage: gitmap rm [-y|--yes] <target>[,<target>...] [<target>...]
       gitmap remove ...
       gitmap del ...

Targets may be:
  - a repo slug/name           (e.g. my-repo)
  - a path                     (./projects/foo, .\macro-ahk, /abs/path)
  - a glob over slug or path   (macro*, gitmap-*)
  - comma-joined combinations  (macro*,gitmap*)

Default: prompts before deleting each repo folder on disk.
With -y/--yes: deletes on-disk folder and DB row without prompting.

Examples:
  gitmap rm my-repo
  gitmap rm macro*
  gitmap rm macro*,gitmap*
  gitmap rm -y macro*
  gitmap rm ./projects/foo ../bar
`

// runRm handles `gitmap rm`. Supports globs, comma-joined targets,
// the -y/--yes auto-confirm flag, and removes the on-disk folder in
// addition to the DB row and .gitmap/output/gitmap.json entry.
func runRm(args []string) error {
	checkHelp("rm", args)
	yes, dbOnly, rest := parseRmFlags(args)
	targets := expandRmTargets(rest)
	if len(targets) == 0 {
		fmt.Fprint(os.Stderr, rmUsage)

		return apperror.NewValidationError("target repository required")
	}

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, "rm: open db:")
	}

	defer db.Close()

	return executeRmTargets(db, targets, yes, dbOnly)
}

func executeRmTargets(db *store.DB, targets []string, isYes, isDbOnly bool) error {
	matches, missing := resolveRmMatches(db, targets)
	reportMissingRmTargets(db, missing)
	if len(matches) == 0 {
		return nil
	}

	_ = removeRmMatches(db, matches, isYes, isDbOnly)
	return nil
}

func buildRmNotFoundError(missing []string) *apperror.AppError {
	if len(missing) == 1 {
		return apperror.NewNotFoundError(fmt.Sprintf("no repository matched %q", missing[0]))
	}

	if len(missing) > 1 {
		return apperror.NewNotFoundError(fmt.Sprintf("no repository matched: %s", strings.Join(missing, ", ")))
	}

	return apperror.NewNotFoundError("no repository matched")
}

func resolveRmMatches(db *store.DB, targets []string) ([]model.ScanRecord, []string) {
	matches, missing := ResolveMultiRepos(db, targets)
	extraMatches, finalMissing := resolveDiskTargets(missing)
	matches = append(matches, extraMatches...)

	return matches, finalMissing
}

func resolveDiskTargets(missing []string) ([]model.ScanRecord, []string) {
	var matches []model.ScanRecord
	var stillMissing []string
	for _, m := range missing {
		rec, isFound := checkUntrackedDir(m)
		if isFound {
			matches = append(matches, rec)
			continue
		}

		stillMissing = append(stillMissing, m)
	}

	return matches, stillMissing
}

func checkUntrackedDir(target string) (model.ScanRecord, bool) {
	abs, err := filepath.Abs(target)
	isDir := err == nil && fsutil.DirExists(abs)
	if !isDir {
		return model.ScanRecord{}, false
	}

	return model.ScanRecord{
		AbsolutePath: abs,
		Slug:         filepath.Base(abs) + " (untracked)",
	}, true
}

func reportMissingRmTargets(db *store.DB, missing []string) {
	for _, m := range missing {
		fmt.Fprintf(os.Stderr, "rm: no repo matched %q\n", m)
		PrintRepoSuggestions(db, m)
	}
}

func parseRmFlags(args []string) (bool, bool, []string) {
	yes, dbOnly := false, false
	var out []string
	for _, a := range args {
		switch a {
		case "-y", "--yes":
			yes = true
		case "--db-only":
			dbOnly = true
		default:
			out = append(out, a)
		}
	}

	return yes, dbOnly, out
}

func expandRmTargets(args []string) []string {
	var out []string
	for _, a := range args {
		tokens := strings.Split(a, ",")
		out = appendValidTargets(out, tokens)
	}

	return out
}

func appendValidTargets(out, tokens []string) []string {
	for _, p := range tokens {
		clean := strings.TrimRight(strings.TrimSpace(p), "/\\")
		if clean != "" {
			out = append(out, clean)
		}
	}

	return out
}

func removeRmMatches(db *store.DB, matches []model.ScanRecord, isYes, isDbOnly bool) bool {
	reader := bufio.NewReader(os.Stdin)
	isSuccess := true
	for _, r := range matches {
		if !processSingleRm(db, reader, r, isYes, isDbOnly) {
			isSuccess = false
		}
	}

	return isSuccess
}

func processSingleRm(db *store.DB, reader *bufio.Reader, r model.ScanRecord, isYes, isDbOnly bool) bool {
	isConfirmed := isYes || confirmRemove(reader, r, isDbOnly)
	if !isConfirmed {
		fmt.Printf("skip: %s\n", r.Slug)

		return true
	}

	if err := removeRepoFully(db, r, isDbOnly); err != nil {
		fmt.Fprintf(os.Stderr, "rm: %s: %v\n", r.Slug, err)

		return false
	}

	fmt.Printf("removed: %s (%s)\n", r.Slug, r.AbsolutePath)

	return true
}

func confirmRemove(r *bufio.Reader, rec model.ScanRecord, isDbOnly bool) bool {
	action := resolveConfirmAction(isDbOnly)
	fmt.Printf("%s %s\n  %s ? [y/N] ", action, rec.Slug, rec.AbsolutePath)
	line, _ := r.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))

	return ans == "y" || ans == "yes"
}

func resolveConfirmAction(isDbOnly bool) string {
	if isDbOnly {
		return "Untrack from database"
	}

	return "Delete folder and untrack"
}

func removeRepoFully(db *store.DB, r model.ScanRecord, isDbOnly bool) *apperror.AppError {
	if err := removeRepoDisk(r.AbsolutePath, isDbOnly); err != nil {
		return err
	}

	if err := removeRepoDB(db, r); err != nil {
		return err
	}

	removeRepoFromJSON(r)
	_ = vscodepm.RemoveEntry(r.AbsolutePath)
	_ = desktop.RemoveRepo(r.AbsolutePath)

	return nil
}

func removeRepoFromJSON(r model.ScanRecord) {
	paths := []string{
		filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile),
		filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
	}
	for _, p := range paths {
		filterAndSaveJSON(p, r)
	}
}

func filterAndSaveJSON(path string, r model.ScanRecord) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	records, err := model.LoadStatusRecords(path)
	if err != nil {
		return
	}

	filtered := filterOutRecord(records, r)
	data, err := json.MarshalIndent(filtered, "", constants.JSONIndent)
	if err != nil {
		return
	}

	_ = os.WriteFile(path, append(data, '\n'), 0644)
}

func filterOutRecord(records []model.ScanRecord, target model.ScanRecord) []model.ScanRecord {
	var out []model.ScanRecord
	for _, rec := range records {
		if !isMatchingScanRecord(rec, target) {
			out = append(out, rec)
		}
	}

	return out
}

func isMatchingScanRecord(a, b model.ScanRecord) bool {
	if fsutil.EqualPaths(a.AbsolutePath, b.AbsolutePath) {
		return true
	}

	if len(b.Slug) > 0 && strings.EqualFold(a.Slug, b.Slug) {
		return true
	}

	if len(b.RepoName) > 0 && strings.EqualFold(a.RepoName, b.RepoName) {
		return true
	}

	return false
}

func removeRepoDB(db *store.DB, r model.ScanRecord) *apperror.AppError {
	if db == nil {
		return nil
	}

	ctx := context.Background()
	wrapper, appErr := dbengine.WrapDb(db.Conn(), dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	return runRemoveRepoTx(ctx, wrapper, r)
}

func runRemoveRepoTx(ctx context.Context, wrapper *dbengine.DbWrapper, r model.ScanRecord) *apperror.AppError {
	return wrapper.WithImmediateTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		repoID := resolveRepoIDForRm(ctx, tx, r)

		return executeRemoveRepoTx(ctx, tx, repoID, r.AbsolutePath)
	})
}

func resolveRepoIDForRm(ctx context.Context, tx *dbengine.TxWrapper, r model.ScanRecord) int64 {
	if r.ID > 0 {
		return r.ID
	}

	var id int64
	row, appErr := tx.QueryRow(ctx, "SELECT RepoId FROM Repo WHERE AbsolutePath = ?", r.AbsolutePath)
	if appErr == nil && row.Scan(&id) == nil {
		return id
	}

	return scanFallbackRepoID(ctx, tx, r.AbsolutePath)
}

func scanFallbackRepoID(ctx context.Context, tx *dbengine.TxWrapper, absPath string) int64 {
	var id int64
	row, appErr := tx.QueryRow(ctx, "SELECT Id FROM Repo WHERE AbsolutePath = ?", absPath)
	if appErr != nil {
		return 0
	}

	if err := row.Scan(&id); err != nil {
		return 0
	}

	return id
}

func removeAliasIfPresent(ctx context.Context, tx *dbengine.TxWrapper, repoID int64) *apperror.AppError {
	if repoID <= 0 {
		return nil
	}

	if _, err := tx.Exec(ctx, "DELETE FROM Alias WHERE RepoId = ?", repoID); err != nil {
		return apperror.WrapSimple(err, "rm: delete alias")
	}

	return nil
}

func executeRemoveRepoTx(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, absPath string) *apperror.AppError {
	if err := removeAliasIfPresent(ctx, tx, repoID); err != nil {
		return err
	}

	query := "DELETE FROM Repo WHERE AbsolutePath = ?"
	if _, err := tx.Exec(ctx, query, absPath); err != nil {
		return apperror.WrapSimple(err, "rm: delete repo")
	}

	return recordRmHistory(ctx, tx, absPath)
}

func recordRmHistory(ctx context.Context, tx *dbengine.TxWrapper, absPath string) *apperror.AppError {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO CommandHistory (Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(ctx, query, "rm", "", absPath, "", now, now, 0, 0, "removed repository", 1)
	if err != nil {
		return apperror.WrapSimple(err, "record rm history")
	}

	return nil
}

func removeRepoDisk(absPath string, isDbOnly bool) *apperror.AppError {
	if isDbOnly {
		return nil
	}

	if err := safeRemoveWithRetry(absPath); err != nil {
		if isProcessLockError(err) {
			fmt.Fprintf(os.Stderr, "rm: warning: some files in %s are locked by another process; untracking repo\n", absPath)
			return nil
		}

		return apperror.WrapSimple(err, "remove dir")
	}

	return nil
}
