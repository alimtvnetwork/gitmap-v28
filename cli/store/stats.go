package store

import (
	"fmt"
	"math"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// RoundAvgDuration converts float64 duration to rounded int64.
func RoundAvgDuration(duration float64) int64 {
	return int64(math.Round(duration))
}

// QueryCommandStats returns per-command aggregated statistics.
func (db *DB) QueryCommandStats() ([]model.CommandStats, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLStatsPerCommand).Destruct()

	if err != nil {
		return nil, fmt.Errorf(constants.ErrStatsQuery, err)
	}

	defer rows.Close()

	return scanStatsRows(rows)
}

// QueryCommandStatsFor returns stats for a single command.
func (db *DB) QueryCommandStatsFor(command string) ([]model.CommandStats, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLStatsForCommand, command).Destruct()

	if err != nil {
		return nil, fmt.Errorf(constants.ErrStatsQuery, err)
	}

	defer rows.Close()

	return scanStatsRows(rows)
}

// QueryOverallStats returns the overall summary row.
func (db *DB) QueryOverallStats() (model.OverallStats, error) {
	var s model.OverallStats
	var avgDuration float64
	row := QueryRowWrapper(db.conn, constants.SQLStatsOverall)
	err := row.Scan(&s.TotalCommands, &s.UniqueCommands,
		&s.TotalSuccess, &s.TotalFail, &s.OverallFailRate, &avgDuration)

	if err != nil {
		return s, fmt.Errorf(constants.ErrStatsQuery, err)
	}

	s.AvgDuration = RoundAvgDuration(avgDuration)

	return s, nil
}

// scanSingleStatRow scans a single command stats row.
func scanSingleStatRow(rows interface{ Scan(dest ...any) error }) (model.CommandStats, error) {
	var s model.CommandStats
	var avgDuration float64
	err := rows.Scan(&s.Command, &s.TotalRuns, &s.SuccessCount,
		&s.FailCount, &s.FailRate, &avgDuration,
		&s.MinDuration, &s.MaxDuration, &s.LastUsed)

	if err != nil {
		return s, fmt.Errorf(constants.ErrStatsQuery, err)
	}

	s.AvgDuration = RoundAvgDuration(avgDuration)

	return s, nil
}

type statsRowScanner interface {
	Next() bool
	Scan(dest ...any) error
}

// scanStatsRows reads all rows into CommandStats slices.
func scanStatsRows(rows statsRowScanner) ([]model.CommandStats, error) {
	var results []model.CommandStats

	for rows.Next() {
		stat, err := scanSingleStatRow(rows)

		if err != nil {
			return nil, err
		}

		results = append(results, stat)
	}

	return results, nil
}
