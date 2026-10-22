package cmdagy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// AgyProcessInfoWrap wraps AgyProcessInfo in a Result container.
	AgyProcessInfoWrap = result.Wrap[AgyProcessInfo]

	// AgyConvIDWrap wraps string conversation ID in a Result container.
	AgyConvIDWrap = result.Wrap[string]
)
