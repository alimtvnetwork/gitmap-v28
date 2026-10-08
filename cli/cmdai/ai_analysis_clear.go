package cmdai

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	clearCmd = &cobra.Command{
		Use:     "clear",
		Aliases: []string{"prune", "truncate", "purge"},
		Short:   "Prune or clear recorded AI analysis sessions and reclaim disk space",
		RunE:    runClearCmd,
	}

	clearFlags    AiClearOptions
	clearOlderStr string
)

func runClearCmd(cmd *cobra.Command, args []string) error {
	dur, durErr := resolveClearDuration(clearOlderStr)
	if durErr != nil {
		return durErr
	}
	if dur > 0 {
		clearFlags.OlderThan = dur
	}

	count, err := ExecuteAiClear(clearFlags)
	if err != nil {
		return err
	}

	fmt.Printf("Cleared %d AI analysis session(s) and executed database VACUUM.\n", count)

	return nil
}

func resolveClearDuration(raw string) (time.Duration, error) {
	clean := strings.TrimSpace(raw)
	isEmpty := clean == ""
	if isEmpty {
		return 0, nil
	}

	return ParseRetentionDuration(clean)
}

// ExecuteAiClear prunes or truncates AI analysis sessions based on options.
func ExecuteAiClear(opts AiClearOptions) (int64, error) {
	db, err := OpenAiAnalysisSplitDB(opts.RepoRoot)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	if opts.All {
		return executePurgeAll(db, opts.Force)
	}

	hasTaskID := strings.TrimSpace(opts.TaskId) != ""
	if hasTaskID {
		return executeDeleteTask(db, strings.TrimSpace(opts.TaskId))
	}

	hasOlderThan := opts.OlderThan > 0
	if hasOlderThan {
		return executePruneOlderThan(db, opts.OlderThan)
	}

	return 0, apperror.NewValidation("cmd.ai.clear", "E3011", "specify --task, --older-than, or --all --force")
}

func executePurgeAll(db *sql.DB, force bool) (int64, error) {
	if !force {
		return 0, apperror.NewValidation("cmd.ai.clear", "E3011", "--force flag is required when clearing with --all")
	}

	delLines := store.ExecWrapper(db, "DELETE FROM AiAnalysisLine;")
	if delLines.IsFailure {
		return 0, apperror.WrapSimple(delLines.Error, "cmdai.executePurgeAll.lines")
	}

	delTasks := store.ExecWrapper(db, "DELETE FROM AiAnalysisTask;")
	if delTasks.IsFailure {
		return 0, apperror.WrapSimple(delTasks.Error, "cmdai.executePurgeAll.tasks")
	}

	_ = store.ExecWrapper(db, "VACUUM;")

	affected, _ := delTasks.Data.RowsAffected()

	return affected, nil
}

func executeDeleteTask(db *sql.DB, taskID string) (int64, error) {
	query := "DELETE FROM AiAnalysisTask WHERE task_id = ?"
	res := store.ExecWrapper(db, query, taskID)
	if res.IsFailure {
		return 0, apperror.WrapSimple(res.Error, "cmdai.executeDeleteTask")
	}

	affected, err := res.Data.RowsAffected()
	if err != nil {
		return 0, apperror.WrapSimple(err, "cmdai.executeDeleteTask.rows")
	}

	hasAffected := affected > 0
	if hasAffected {
		_ = store.ExecWrapper(db, "VACUUM;")
	}

	return affected, nil
}

func executePruneOlderThan(db *sql.DB, dur time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-dur).Format(time.RFC3339)
	query := "DELETE FROM AiAnalysisTask WHERE created_at < ?"
	res := store.ExecWrapper(db, query, cutoff)
	if res.IsFailure {
		return 0, apperror.WrapSimple(res.Error, "cmdai.executePruneOlderThan")
	}

	affected, err := res.Data.RowsAffected()
	if err != nil {
		return 0, apperror.WrapSimple(err, "cmdai.executePruneOlderThan.rows")
	}

	hasAffected := affected > 0
	if hasAffected {
		_ = store.ExecWrapper(db, "VACUUM;")
	}

	return affected, nil
}

// ParseRetentionDuration parses strings like "14d", "30d", "2w", or "24h".
func ParseRetentionDuration(raw string) (time.Duration, error) {
	clean := strings.ToLower(strings.TrimSpace(raw))
	isEmpty := clean == ""
	if isEmpty {
		return 0, apperror.NewValidation("cmd.ai.clear", "E3011", "duration string is empty")
	}

	isDays := strings.HasSuffix(clean, "d")
	if isDays {
		return parseDayDuration(clean[:len(clean)-1])
	}

	isWeeks := strings.HasSuffix(clean, "w")
	if isWeeks {
		return parseWeekDuration(clean[:len(clean)-1])
	}

	dur, err := time.ParseDuration(clean)
	if err != nil {
		return 0, apperror.NewValidation("cmd.ai.clear", "E3011", fmt.Sprintf("invalid duration %q: %v", raw, err))
	}

	return dur, nil
}

func parseDayDuration(valStr string) (time.Duration, error) {
	days, err := strconv.Atoi(strings.TrimSpace(valStr))
	isValid := err == nil && days > 0
	if !isValid {
		return 0, apperror.NewValidation("cmd.ai.clear", "E3011", fmt.Sprintf("invalid days count %q", valStr))
	}

	return time.Duration(days) * 24 * time.Hour, nil
}

func parseWeekDuration(valStr string) (time.Duration, error) {
	weeks, err := strconv.Atoi(strings.TrimSpace(valStr))
	isValid := err == nil && weeks > 0
	if !isValid {
		return 0, apperror.NewValidation("cmd.ai.clear", "E3011", fmt.Sprintf("invalid weeks count %q", valStr))
	}

	return time.Duration(weeks) * 7 * 24 * time.Hour, nil
}

func initClearCommands() {
	clearCmd.Flags().StringVarP(&clearFlags.TaskId, "task", "t", "", "Delete specific task session")
	clearCmd.Flags().StringVar(&clearOlderStr, "older-than", "", "Prune sessions older than duration (e.g. 14d, 30d, 24h)")
	clearCmd.Flags().BoolVar(&clearFlags.All, "all", false, "Purge all AI analysis records")
	clearCmd.Flags().BoolVar(&clearFlags.Force, "force", false, "Confirm destructive purge operations")

	AiCmd.AddCommand(clearCmd)
}
