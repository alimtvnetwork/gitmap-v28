package cmdpipeline

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	PEFormatProfileSlice       = []PEFormatProfile
	PEFormatProfileSliceWrap   = result.WrapSlice[PEFormatProfile]
	PipelineSyncResultWrap     = result.Wrap[*PipelineSyncResult]
)
