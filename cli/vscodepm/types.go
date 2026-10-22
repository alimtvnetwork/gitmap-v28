package vscodepm

import "github.com/alimtvnetwork/gitmap-v28/cli/result"

type (
	EntrySliceWrap  = result.WrapSlice[Entry]
	EntryWrap       = result.Wrap[Entry]
	SyncSummaryWrap = result.Wrap[SyncSummary]
)
