package clonepick

import "github.com/alimtvnetwork/gitmap-v28/cli/result"

type (
	StringSliceWrap = result.StringSliceResult
	StringWrap      = result.StringResult
	PlanWrap        = result.Wrap[Plan]
)
