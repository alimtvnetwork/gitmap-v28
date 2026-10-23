package cmdpull

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PullEtaEstimator calculates rolling dynamic ETA based on completed historical runs.
type PullEtaEstimator struct {
	commandType string
	targetNode  string
	sampleLimit int
	minSamples  int
	baselineSec float64
}

const sqlSelectEtaSamples = `SELECT PullEtaHistoryId, CommandType, TargetNode, TotalRepos, DurationSeconds, CreatedAt
FROM PullEtaHistory
WHERE CommandType = ? AND TargetNode = ?
ORDER BY CreatedAt DESC LIMIT ?;`

// NewPullEtaEstimator creates an estimator configured with canonical rolling limits.
func NewPullEtaEstimator(cmdType, targetNode string) *PullEtaEstimator {
	return &PullEtaEstimator{
		commandType: cmdType,
		targetNode:  resolveTargetNode(targetNode),
		sampleLimit: constants.PullTaskRollingEtaSampleLimit,
		minSamples:  constants.PullTaskRollingEtaMinSamples,
		baselineSec: constants.PullTaskDefaultBaselineEtaSec,
	}
}

// CalculateRemainingEta calculates remaining seconds for totalRepos given pulledRepos count.
func (e *PullEtaEstimator) CalculateRemainingEta(db *PullTaskSplitDB, totalRepos, pulledRepos int) (int, *apperror.AppError) {
	remaining := totalRepos - pulledRepos
	if remaining <= 0 {
		return 0, nil
	}

	secPerRepo, err := e.CalculateAverageSecondsPerRepo(db)
	if err != nil {
		return 0, err
	}

	etaFloat := secPerRepo * float64(remaining)

	return int(math.Ceil(etaFloat)), nil
}

// CalculateAverageSecondsPerRepo evaluates weighted moving average rate from historical runs.
func (e *PullEtaEstimator) CalculateAverageSecondsPerRepo(db *PullTaskSplitDB) (float64, *apperror.AppError) {
	samples, err := e.queryEtaSamples(db)
	if err != nil {
		return e.baselineSec, err
	}

	hasEnoughSamples := len(samples) >= e.minSamples
	if !hasEnoughSamples {
		return e.baselineSec, nil
	}

	return computeWeightedAverage(samples, e.baselineSec), nil
}

func (e *PullEtaEstimator) queryEtaSamples(db *PullTaskSplitDB) ([]PullEtaHistoryRecord, *apperror.AppError) {
	rows, err := db.Conn().Query(sqlSelectEtaSamples, e.commandType, e.targetNode, e.sampleLimit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_eta.query")
	}

	defer rows.Close()

	return scanEtaHistoryRows(rows)
}

func scanEtaHistoryRows(rows *sql.Rows) ([]PullEtaHistoryRecord, *apperror.AppError) {
	var records []PullEtaHistoryRecord
	for rows.Next() {
		r, err := scanSingleEtaRow(rows)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}

func scanSingleEtaRow(rows *sql.Rows) (PullEtaHistoryRecord, *apperror.AppError) {
	var r PullEtaHistoryRecord
	err := rows.Scan(&r.PullEtaHistoryId, &r.CommandType, &r.TargetNode, &r.TotalRepos, &r.DurationSeconds, &r.CreatedAt)
	if err != nil {
		return r, apperror.WrapSimple(err, "pull_eta.scan")
	}

	return r, nil
}

func computeWeightedAverage(samples []PullEtaHistoryRecord, fallback float64) float64 {
	totalDuration, totalRepos := sumSampleMetrics(samples)

	hasValidData := totalRepos > 0 && totalDuration > 0
	if !hasValidData {
		return fallback
	}

	return totalDuration / float64(totalRepos)
}

func sumSampleMetrics(samples []PullEtaHistoryRecord) (float64, int) {
	var totalDuration float64
	var totalRepos int
	for _, s := range samples {
		totalDuration += s.DurationSeconds
		totalRepos += s.TotalRepos
	}

	return totalDuration, totalRepos
}

// FormatEtaDuration formats seconds into a human-readable duration string.
func FormatEtaDuration(seconds int) string {
	if seconds <= 0 {
		return "0s"
	}

	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	return formatMinutesAndSeconds(seconds/60, seconds%60)
}

func formatMinutesAndSeconds(mins, remSec int) string {
	if remSec == 0 {
		return fmt.Sprintf("%dm", mins)
	}

	return fmt.Sprintf("%dm %ds", mins, remSec)
}

// CalculateRollingEta is a convenience wrapper to calculate remaining ETA in one shot.
func CalculateRollingEta(db *PullTaskSplitDB, cmdType, targetNode string, totalRepos, pulledRepos int) (int, *apperror.AppError) {
	estimator := NewPullEtaEstimator(cmdType, targetNode)

	return estimator.CalculateRemainingEta(db, totalRepos, pulledRepos)
}
