package visibility

import "github.com/alimtvnetwork/gitmap-v28/cli/result"

type (
	PatternSliceWrap = result.WrapSlice[Pattern]
	PatternWrap      = result.Wrap[Pattern]
)
