package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ReplaceModeType enumerates the invocation shapes the spec accepts.
type ReplaceModeType int

const (
	ReplaceModeTypeUnknown ReplaceModeType = iota
	ReplaceModeTypeLiteral
	ReplaceModeTypeVersionN
	ReplaceModeTypeAll
	ReplaceModeTypeAudit
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

	if arg == "history" || arg == "audit" {
		return ReplaceModeTypeAudit
	}

	if looksLikeDashN(arg) {
		return ReplaceModeTypeVersionN
	}

	return ReplaceModeTypeUnknown
}
