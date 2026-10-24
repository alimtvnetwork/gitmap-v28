package db

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/db/enums"

// CommandKindType aliases the strongly-typed enum from the enums package.
type CommandKindType = enums.CommandKindType

const (
	CommandKindPsCommand      = enums.CommandKindPsCommand
	CommandKindCmdCommand     = enums.CommandKindCmdCommand
	CommandKindInstall        = enums.CommandKindInstall
	CommandKindGitPull        = enums.CommandKindGitPull
	CommandKindGitPush        = enums.CommandKindGitPush
	CommandKindGitCommit      = enums.CommandKindGitCommit
	CommandKindGitStatus      = enums.CommandKindGitStatus
	CommandKindProjRun        = enums.CommandKindProjRun
	CommandKindProjCreateCICD = enums.CommandKindProjCreateCICD
	CommandKindRestart        = enums.CommandKindRestart
	CommandKindShutdown       = enums.CommandKindShutdown
	CommandKindLogoff         = enums.CommandKindLogoff
	CommandKindBashCommand    = enums.CommandKindBashCommand
	CommandKindShellCommand   = enums.CommandKindShellCommand
)

// ParseCommandKind aliases the enums package parser so existing callers keep working.
func ParseCommandKind(s string) (CommandKindType, *apperror.AppError) {
	return enums.ParseCommandKind(s)
}

// ResultStatusType aliases the strongly-typed enum from the enums package.
type ResultStatusType = enums.ResultStatusType

const (
	ResultStatusPending      = enums.ResultStatusPending
	ResultStatusSucceeded    = enums.ResultStatusSucceeded
	ResultStatusFailed       = enums.ResultStatusFailed
	ResultStatusSkipped      = enums.ResultStatusSkipped
	ResultStatusDeferred     = enums.ResultStatusDeferred
	ResultStatusRequiresAuth = enums.ResultStatusRequiresAuth
)

// ParseResultStatus aliases the enums package parser so existing callers keep working.
func ParseResultStatus(s string) (ResultStatusType, *apperror.AppError) {
	return enums.ParseResultStatus(s)
}
