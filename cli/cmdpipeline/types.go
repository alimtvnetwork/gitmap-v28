package cmdpipeline

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	PEFormatProfileSlice     = []PEFormatProfile
	PEFormatProfileSliceWrap = result.WrapSlice[PEFormatProfile]
	PipelineSyncResultWrap   = result.Wrap[*PipelineSyncResult]

	SaveParsedFailedJobsParams struct {
		DB      *pipelinedb.PipelineSplitDb
		Repo    string
		Run     ghRunItem
		Jobs    []FailedJobItem
		RawLogs string
	}
)
