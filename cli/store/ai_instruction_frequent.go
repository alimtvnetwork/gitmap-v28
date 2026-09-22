package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AiFrequentCommand aggregates metrics for frequently executed AI commands.
type AiFrequentCommand struct {
	CommandText    string `json:"command_text"`
	CategoryCode   string `json:"category_code"`
	RunCount       int    `json:"run_count"`
	LastExecutedAt string `json:"last_executed_at"`
	AvgDurationMs  int    `json:"avg_duration_ms"`
	SuccessCount   int    `json:"success_count"`
}

// QueryFrequentAiCommands retrieves the most frequent AI commands from the split DB.
func (db *AiInstructionSplitDB) QueryFrequentAiCommands(limit int) ([]AiFrequentCommand, error) {
	if db == nil || db.conn == nil {
		return nil, apperror.NewExecutionError("database connection is nil")
	}

	actualLimit := normalizeFrequentLimit(limit)
	q := `SELECT
		CommandText,
		CategoryCode,
		COUNT(*) AS RunCount,
		MAX(ExecutedAt) AS LastExecutedAt,
		CAST(AVG(DurationMs) AS INTEGER) AS AvgDurationMs,
		SUM(CASE WHEN IsSuccess = 1 THEN 1 ELSE 0 END) AS SuccessCount
	FROM AiExecutionHistory
	WHERE TRIM(CommandText) != ''
	GROUP BY CommandText
	ORDER BY RunCount DESC, LastExecutedAt DESC
	LIMIT ?`

	rows, err := db.conn.Query(q, actualLimit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ai_instruction.query_frequent")
	}
	defer rows.Close()

	return scanFrequentCommands(rows)
}

func normalizeFrequentLimit(limit int) int {
	if limit > 0 {
		return limit
	}
	return 25
}

func scanFrequentCommands(rows *sql.Rows) ([]AiFrequentCommand, error) {
	var list []AiFrequentCommand
	for rows.Next() {
		var item AiFrequentCommand
		err := rows.Scan(
			&item.CommandText,
			&item.CategoryCode,
			&item.RunCount,
			&item.LastExecutedAt,
			&item.AvgDurationMs,
			&item.SuccessCount,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "ai_instruction.scan_frequent")
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// GetFrequentAiCommands opens the default AI split DB, queries frequent commands, and closes.
func GetFrequentAiCommands(limit int) ([]AiFrequentCommand, error) {
	db, err := OpenAiInstructionSplitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db.QueryFrequentAiCommands(limit)
}
