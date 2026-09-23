package cmdpull

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PullSessionTelemetry contains the overall session context to be stored in PullRun.
type PullSessionTelemetry struct {
	CommandType   string
	WorkingDir    string
	TotalRepos    int
	PulledRepos   int
	SkippedRepos  int
	SuccessCount  int
	FailedCount   int
	IsEfficient   bool
	Duration      time.Duration
	GitMapVersion string
	Notes         string
	Comments      string
}

// ConvertToStoreRunRecord maps session telemetry to a store.PullRunRecord.
func (t PullSessionTelemetry) ConvertToStoreRunRecord() *store.PullRunRecord {
	return &store.PullRunRecord{
		CommandType:   t.CommandType,
		WorkingDir:    t.WorkingDir,
		TotalRepos:    t.TotalRepos,
		PulledRepos:   t.PulledRepos,
		SkippedRepos:  t.SkippedRepos,
		SuccessCount:  t.SuccessCount,
		FailedCount:   t.FailedCount,
		IsEfficient:   t.IsEfficient,
		DurationMs:    t.Duration.Milliseconds(),
		GitMapVersion: t.GitMapVersion,
		Notes:         t.Notes,
		Comments:      t.Comments,
	}
}
