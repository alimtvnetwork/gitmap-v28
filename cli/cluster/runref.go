package cluster

import (
	"database/sql"
	"fmt"
	"time"
)

// RunRefGenerator generates a RUN-<YYYYMMDD>-<NNN> identifier using the current date
// and a daily counter read from the ClusterRun table.
func RunRefGenerator(db *sql.DB) (string, error) {
	now := time.Now().UTC()
	dateStr := now.Format("20060102")

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startStr := startOfDay.Format(time.RFC3339)

	var count int
	query := `SELECT COUNT(*) FROM ClusterRun WHERE StartedAt >= ?`
	err := db.QueryRow(query, startStr).Scan(&count)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("RUN-%s-%03d", dateStr, count+1), nil
}
