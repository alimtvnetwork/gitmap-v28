package cmdreplace

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdreplace/enums"
import "slices"

// ReplaceModeType enumerates the invocation shapes the spec accepts.
// ReplaceModeType aliases the strongly-typed enum from the enums package.
type ReplaceModeType = enums.ReplaceModeType

const (
	ReplaceModeTypeUnknown  = enums.ReplaceModeTypeUnknown
	ReplaceModeTypeLiteral  = enums.ReplaceModeTypeLiteral
	ReplaceModeTypeVersionN = enums.ReplaceModeTypeVersionN
	ReplaceModeTypeAll      = enums.ReplaceModeTypeAll
	ReplaceModeTypeAudit    = enums.ReplaceModeTypeAudit
)

// classifyReplaceMode picks the operating mode from positional args and audit flag.
func classifyReplaceMode(positional []string, opts replaceOpts) ReplaceModeType {
	if opts.audit {
		return ReplaceModeTypeAudit
	}

	if len(positional) == constants.ReplaceMinPositionalArgs {
		return classifySingleArgMode(positional[0])
	}

	if len(positional) == constants.ReplaceLiteralArgsCount {
		return ReplaceModeTypeLiteral
	}

	return ReplaceModeTypeUnknown
}

func classifySingleArgMode(arg string) ReplaceModeType {
	if arg == constants.ReplaceSubcmdAll {
		return ReplaceModeTypeAll
	}

	if slices.Contains([]string{"history", "audit"}, arg) {
		return ReplaceModeTypeAudit
	}

	if looksLikeDashN(arg) {
		return ReplaceModeTypeVersionN
	}

	return ReplaceModeTypeUnknown
}
