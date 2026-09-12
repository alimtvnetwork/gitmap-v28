package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
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
// addition to the DB row.
func runRm(args []string) error {
	checkHelp("rm", args)
	yes, dbOnly, rest := parseRmFlags(args)
	targets := expandRmTargets(rest)
	if len(targets) == 0 {
		fmt.Fprint(os.Stderr, rmUsage)

		return apperror.NewSimple("fatal error", "E9000")
	}

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, "rm: open db:")
	}

	defer db.Close()

	return executeRmTargets(db, targets, yes, dbOnly)
}

func executeRmTargets(db *store.DB, targets []string, yes, dbOnly bool) error {
	matches, missing := resolveRmMatches(db, targets)
	reportMissingRmTargets(db, missing)
	if len(matches) == 0 {
		return apperror.NewSimple("fatal error", "E9000")
	}

	if removeRmMatches(db, matches, yes, dbOnly) {
		cliexit.HandleError(nil, 0)
	}

	return apperror.NewSimple("fatal error", "E9000")
}

func resolveRmMatches(db *store.DB, targets []string) ([]model.ScanRecord, []string) {
	matches, missing := ResolveMultiRepos(db, targets)
	var finalMissing []string
	for _, m := range missing {
		abs, err := filepath.Abs(m)
		if err == nil && fsutil.DirExists(abs) {
			matches = append(matches, model.ScanRecord{
				AbsolutePath: abs,
				Slug:         filepath.Base(abs) + " (untracked)",
			})
			continue
		}

		finalMissing = append(finalMissing, m)
	}

	return matches, finalMissing
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
		for _, p := range strings.Split(a, ",") {
			if p = strings.TrimSpace(p); p != "" {
				p = strings.TrimRight(p, "/\\")
				out = append(out, p)
			}
		}
	}

	return out
}

func removeRmMatches(db *store.DB, matches []model.ScanRecord, yes, dbOnly bool) bool {
	reader := bufio.NewReader(os.Stdin)
	isSuccess := true

	for _, r := range matches {
		if !processSingleRm(db, reader, r, yes, dbOnly) {
			isSuccess = false
		}
	}

	return isSuccess
}

func processSingleRm(db *store.DB, reader *bufio.Reader, r model.ScanRecord, yes, dbOnly bool) bool {
	if !yes && !confirmRemove(reader, r, dbOnly) {
		fmt.Printf("skip: %s\n", r.Slug)

		return true
	}

	if err := removeRepoFully(db, r, dbOnly); err != nil {
		fmt.Fprintf(os.Stderr, "rm: %s: %v\n", r.Slug, err)

		return false
	}

	fmt.Printf("removed: %s (%s)\n", r.Slug, r.AbsolutePath)

	return true
}

func confirmRemove(r *bufio.Reader, rec model.ScanRecord, dbOnly bool) bool {
	action := "Delete folder and untrack"
	if dbOnly {
		action = "Untrack from database"
	}

	fmt.Printf("%s %s\n  %s ? [y/N] ", action, rec.Slug, rec.AbsolutePath)
	line, _ := r.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))

	return ans == "y" || ans == "yes"
}

func removeRepoFully(db *store.DB, r model.ScanRecord, dbOnly bool) error {
	if err := removeRepoDisk(r.AbsolutePath, dbOnly); err != nil {
		return err
	}

	if err := removeRepoDB(db, r); err != nil {
		return fmt.Errorf("db delete: %w", err)
	}

	_ = vscodepm.RemoveEntry(r.AbsolutePath)
	_ = desktop.RemoveRepo(r.AbsolutePath)

	return nil
}

func removeRepoDB(db *store.DB, r model.ScanRecord) error {
	ctx := context.Background()
	wrapper, appErr := dbengine.WrapDb(db.Conn(), dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	if txErr := runRemoveRepoTx(ctx, wrapper, r); txErr != nil {
		return txErr
	}

	return nil
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

func removeRepoDisk(absPath string, dbOnly bool) error {
	if dbOnly {
		return nil
	}

	if err := fsutil.SafeRemoveAll(absPath); err != nil {
		return fmt.Errorf("remove dir: %w", err)
	}

	return nil
}
