package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// QueryCommandStats returns per-command aggregated statistics.
func (db *DB) QueryCommandStats() ([]model.CommandStats, error) {
	rows, err := db.conn.Query(constants.SQLStatsPerCommand)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrStatsQuery, err)
	}
	defer rows.Close()

	return scanStatsRows(rows)
}

// QueryCommandStatsFor returns stats for a single command.
func (db *DB) QueryCommandStatsFor(command string) ([]model.CommandStats, error) {
	rows, err := db.conn.Query(constants.SQLStatsForCommand, command)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrStatsQuery, err)
	}
	defer rows.Close()

	return scanStatsRows(rows)
}

// QueryOverallStats returns the overall summary row.
func (db *DB) QueryOverallStats() (model.OverallStats, error) {
	var s model.OverallStats
	row := db.conn.QueryRow(constants.SQLStatsOverall)

	err := row.Scan(&s.TotalCommands, &s.UniqueCommands,
		&s.TotalSuccess, &s.TotalFail, &s.OverallFailRate, &s.AvgDuration)
	if err != nil {
		return s, fmt.Errorf(constants.ErrStatsQuery, err)
	}

	return s, nil
}

// scanStatsRows reads all rows into CommandStats slices.
func scanStatsRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.CommandStats, error) {
	var results []model.CommandStats

	for rows.Next() {
		var s model.CommandStats
		err := rows.Scan(&s.Command, &s.TotalRuns, &s.SuccessCount,
			&s.FailCount, &s.FailRate, &s.AvgDuration,
			&s.MinDuration, &s.MaxDuration, &s.LastUsed)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrStatsQuery, err)
		}
		results = append(results, s)
	}

	return results, nil
}
