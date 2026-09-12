// Package cmdschedule — types.go defines domain payload models and single reusable Result envelopes.
package cmdschedule

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type (
	// ScheduleExportBundle defines the exported schedule archive payload.
	ScheduleExportBundle struct {
		Task store.SchedulerTask       `json:"task" yaml:"task"`
		Runs []store.ScheduleRunRecord `json:"runs" yaml:"runs"`
	}

	// ScheduleExportOpts defines CLI options for schedule exports.
	ScheduleExportOpts struct {
		TargetName string
		FilePath   string
		Format     string
		ExceptList []string
		IsAll      bool
	}

	// ScheduleExportBundleResult is the canonical single reusable result envelope for bundle slices.
	ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]

	// ScheduleExportBundleSingleResult is the canonical single reusable result envelope for a single bundle.
	ScheduleExportBundleSingleResult = result.Result[ScheduleExportBundle]

	// SchedulerTaskSliceResult is the canonical single reusable result envelope for scheduler task slices.
	SchedulerTaskSliceResult = result.ResultSlice[store.SchedulerTask]
)
