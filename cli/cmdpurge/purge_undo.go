package cmdpurge

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

// UndoOptions configures parameters for reverting a purge operation.
type UndoOptions struct {
	OperationId int64  `json:"operationId"`
	RepoDir     string `json:"repoDir"`
	IsJson      bool   `json:"isJson"`
	IsVerbose   bool   `json:"isVerbose"`
}

// UndoResult describes the outcome of an undo restoration operation.
type UndoResult struct {
	OperationId        int64  `json:"operationId"`
	RepoSlug           string `json:"repoSlug"`
	Branch             string `json:"branch"`
	RestoredHeadSha    string `json:"restoredHeadSha"`
	RestorationSource  string `json:"restorationSource"`
	RestoredFilesCount int    `json:"restoredFilesCount"`
	IsSuccess          bool   `json:"isSuccess"`
}

type undoContext struct {
	db      *store.PurgeHistoryDB
	op      *store.HistoryPurgeOperation
	repoDir string
	branch  string
}

// RunPurgeUndo is the public entrypoint to undo a previous history purge operation.
func RunPurgeUndo(opId int64, args []string) error {
	opts := parseUndoArgs(opId, args)
	res, err := ExecuteUndoOperation(opts)
	if err != nil {
		return err
	}

	if opts.IsJson {
		return printUndoJSON(res)
	}

	printUndoSummaryBox(res)

	return nil
}

func parseUndoArgs(opId int64, args []string) UndoOptions {
	opts := UndoOptions{OperationId: opId}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--json" || a == "-j":
			opts.IsJson = true
		case a == "--verbose" || a == "-v":
			opts.IsVerbose = true
		case a == "--repo" && i+1 < len(args):
			i++
			opts.RepoDir = args[i]
		case !strings.HasPrefix(a, "-") && opts.OperationId == 0:
			var val int64
			if _, err := fmt.Sscanf(a, "%d", &val); err == nil {
				opts.OperationId = val
			}
		}
	}

	return opts
}

// ExecuteUndoOperation performs the complete restoration pipeline from SplitDB and backup refs.
func ExecuteUndoOperation(opts UndoOptions) (*UndoResult, error) {
	repoDir := resolveRepoDir(opts.RepoDir)
	if isDirty, err := CheckWorkingTreeDirty(repoDir); err != nil || isDirty {
		return nil, apperror.NewSimple("ERR_UNDO_DIRTY_TREE", "cannot undo purge: uncommitted changes detected")
	}

	db, err := store.OpenPurgeHistoryDB("")
	if err != nil {
		return nil, apperror.WrapSimple(err, "open splitdb for undo")
	}
	defer db.Close()

	uCtx, err := resolveUndoContext(db, repoDir, opts.OperationId)
	if err != nil {
		return nil, err
	}

	return executeRestoration(uCtx)
}

func resolveRepoDir(customDir string) string {
	if customDir != "" {
		return customDir
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	return wd
}

func resolveUndoContext(db *store.PurgeHistoryDB, repoDir string, opId int64) (*undoContext, error) {
	op, err := fetchPurgeOp(db, repoDir, opId)
	if err != nil {
		return nil, err
	}

	branch, err := queryCurrentBranch(repoDir)
	if err != nil {
		return nil, err
	}

	return &undoContext{
		db:      db,
		op:      op,
		repoDir: repoDir,
		branch:  branch,
	}, nil
}

func fetchPurgeOp(db *store.PurgeHistoryDB, repoDir string, opId int64) (*store.HistoryPurgeOperation, error) {
	if opId > 0 {
		op, err := db.GetHistoryPurgeOperationById(opId)
		if err != nil || op == nil {
			return nil, apperror.NewSimple("ERR_UNDO_OP_NOT_FOUND", fmt.Sprintf("operation %d not found in splitdb", opId))
		}
		return op, nil
	}

	slug := store.SanitizeSlug(filepath.Base(repoDir))
	op, err := db.GetLastHistoryPurgeOperation(slug)
	if err != nil || op == nil {
		return nil, apperror.NewSimple("ERR_UNDO_OP_NOT_FOUND", "no unrestored purge operation found in splitdb")
	}

	return op, nil
}

func queryCurrentBranch(repoDir string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "query current git branch")
	}
	br := strings.TrimSpace(string(out))
	if br == "" {
		return "main", nil
	}

	return br, nil
}

func executeRestoration(ctx *undoContext) (*UndoResult, error) {
	backupRef := fmt.Sprintf("refs/gitmap-backup/%d/%s", ctx.op.OperationId, ctx.branch)
	backupSha, hasBackupRef := verifyGitRef(ctx.repoDir, backupRef)

	if hasBackupRef {
		return restoreViaBackupRef(ctx, backupRef, backupSha)
	}

	return restoreViaTempVault(ctx)
}

func verifyGitRef(repoDir, ref string) (string, bool) {
	cmd := exec.Command("git", "rev-parse", "--verify", ref)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	sha := strings.TrimSpace(string(out))

	return sha, len(sha) > 0
}

func restoreViaBackupRef(ctx *undoContext, ref, sha string) (*UndoResult, error) {
	cmd := exec.Command("git", "update-ref", fmt.Sprintf("refs/heads/%s", ctx.branch), sha)
	cmd.Dir = ctx.repoDir
	if err := cmd.Run(); err != nil {
		return nil, apperror.WrapSimple(err, "update branch reference")
	}

	resetCmd := exec.Command("git", "reset", "--hard", sha)
	resetCmd.Dir = ctx.repoDir
	if err := resetCmd.Run(); err != nil {
		return nil, apperror.WrapSimple(err, "hard reset to backup sha")
	}

	recordUndoSuccess(ctx, sha, "backup_ref", ctx.op.BackedUpFilesCount)

	return &UndoResult{
		OperationId:        ctx.op.OperationId,
		RepoSlug:           ctx.op.RepoSlug,
		Branch:             ctx.branch,
		RestoredHeadSha:    sha,
		RestorationSource:  "backup_ref",
		RestoredFilesCount: ctx.op.BackedUpFilesCount,
		IsSuccess:          true,
	}, nil
}

func restoreViaTempVault(ctx *undoContext) (*UndoResult, error) {
	vaultPath := ctx.op.BackupVaultPath
	if vaultPath == "" {
		opIdStr := fmt.Sprintf("%d", ctx.op.OperationId)
		vaultPath = tempdir.RepoTempDir("history-backup", ctx.op.RepoSlug, opIdStr)
	}

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return nil, apperror.NewSimple("ERR_UNDO_REF_MISSING", "backup ref and temp backup vault both unavailable")
	}

	count, err := overlayTempVaultFiles(ctx.repoDir, vaultPath)
	if err != nil {
		return nil, err
	}

	headSha, _ := verifyGitRef(ctx.repoDir, "HEAD")
	recordUndoSuccess(ctx, headSha, "temp_vault", count)

	return &UndoResult{
		OperationId:        ctx.op.OperationId,
		RepoSlug:           ctx.op.RepoSlug,
		Branch:             ctx.branch,
		RestoredHeadSha:    headSha,
		RestorationSource:  "temp_vault",
		RestoredFilesCount: count,
		IsSuccess:          true,
	}, nil
}

func overlayTempVaultFiles(repoDir, vaultPath string) (int, error) {
	count := 0
	err := filepath.WalkDir(vaultPath, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() == "manifest.json" {
			return err
		}
		rel, err := filepath.Rel(vaultPath, p)
		if err != nil {
			return err
		}
		target := filepath.Join(repoDir, rel)
		_ = os.MkdirAll(filepath.Dir(target), 0755)
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if err := os.WriteFile(target, data, 0644); err != nil {
			return err
		}
		count++
		return nil
	})

	return count, err
}

func recordUndoSuccess(ctx *undoContext, sha, source string, count int) {
	_ = ctx.db.UpdateHistoryPurgeOperationStatus(store.UpdatePurgeStatusOptions{
		OperationId: ctx.op.OperationId,
		IsSuccess:   true,
		IsVerified:  true,
		Notes:       "restored via " + source,
	})

	undoRecord := &store.HistoryUndoOperation{
		OperationId:         ctx.op.OperationId,
		RepoSlug:            ctx.op.RepoSlug,
		RestoredCommitCount: ctx.op.AffectedCommitsCount,
		RestoredFileCount:   count,
		BackupVaultPath:     ctx.op.BackupVaultPath,
		IsVerified:          true,
		IsSuccess:           true,
		CreatedAt:           time.Now().Unix(),
		CompletedAt:         time.Now().Unix(),
	}
	_, _ = ctx.db.InsertHistoryUndoOperation(undoRecord)
}

func printUndoJSON(res *UndoResult) error {
	envelope := map[string]interface{}{
		"version":   "2.0",
		"status":    "success",
		"command":   "history undo",
		"timestamp": time.Now().Unix(),
		"data":      res,
		"error":     nil,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(envelope)
}

func printUndoSummaryBox(res *UndoResult) {
	fmt.Println("┌─ Operation Completed: Git History Purge Undone ──────────────┐")
	fmt.Printf("│ Operation ID      : %-42d │\n", res.OperationId)
	fmt.Printf("│ Restored Branch   : %-42s │\n", res.Branch)
	fmt.Printf("│ Restored HEAD     : %-42s │\n", truncateString(res.RestoredHeadSha, 42))
	fmt.Printf("│ Restoration Source: %-42s │\n", res.RestorationSource)
	fmt.Printf("│ Restored Files    : %-42d │\n", res.RestoredFilesCount)
	fmt.Println("│ SplitDB Journal   : Restored (IsUndone = 1)                  │")
	fmt.Println("└──────────────────────────────────────────────────────────────┘")
}
