package db

import (
	"fmt"
	"strings"
)

type CommandKindType int

const (
	CommandKindPsCommand      CommandKindType = 1
	CommandKindCmdCommand     CommandKindType = 2
	CommandKindInstall        CommandKindType = 3
	CommandKindGitPull        CommandKindType = 4
	CommandKindGitPush        CommandKindType = 5
	CommandKindGitCommit      CommandKindType = 6
	CommandKindGitStatus      CommandKindType = 7
	CommandKindProjRun        CommandKindType = 8
	CommandKindProjCreateCICD CommandKindType = 9
	CommandKindRestart        CommandKindType = 10
	CommandKindShutdown       CommandKindType = 11
	CommandKindLogoff         CommandKindType = 12
)

func (k CommandKindType) String() string {
	switch k {
	case CommandKindPsCommand:
		return "PsCommand"
	case CommandKindCmdCommand:
		return "CmdCommand"
	case CommandKindInstall:
		return "Install"
	case CommandKindGitPull:
		return "GitPull"
	case CommandKindGitPush:
		return "GitPush"
	case CommandKindGitCommit:
		return "GitCommit"
	case CommandKindGitStatus:
		return "GitStatus"
	case CommandKindProjRun:
		return "ProjRun"
	case CommandKindProjCreateCICD:
		return "ProjCreateCICD"
	case CommandKindRestart:
		return "Restart"
	case CommandKindShutdown:
		return "Shutdown"
	case CommandKindLogoff:
		return "Logoff"
	default:
		return fmt.Sprintf("CommandKindType(%d)", int(k))
	}
}

func ParseCommandKind(s string) (CommandKindType, error) {
	switch strings.ToLower(s) {
	case "pscommand":
		return CommandKindPsCommand, nil
	case "cmdcommand":
		return CommandKindCmdCommand, nil
	case "install":
		return CommandKindInstall, nil
	case "gitpull":
		return CommandKindGitPull, nil
	case "gitpush":
		return CommandKindGitPush, nil
	case "gitcommit":
		return CommandKindGitCommit, nil
	case "gitstatus":
		return CommandKindGitStatus, nil
	case "projrun":
		return CommandKindProjRun, nil
	case "projcreatecicd":
		return CommandKindProjCreateCICD, nil
	case "restart":
		return CommandKindRestart, nil
	case "shutdown":
		return CommandKindShutdown, nil
	case "logoff":
		return CommandKindLogoff, nil
	default:
		return 0, fmt.Errorf("invalid CommandKind: %s", s)
	}
}

type ResultStatusType int

const (
	ResultStatusPending      ResultStatusType = 1
	ResultStatusSucceeded    ResultStatusType = 2
	ResultStatusFailed       ResultStatusType = 3
	ResultStatusSkipped      ResultStatusType = 4
	ResultStatusDeferred     ResultStatusType = 5
	ResultStatusRequiresAuth ResultStatusType = 6
)

func (s ResultStatusType) String() string {
	switch s {
	case ResultStatusPending:
		return "Pending"
	case ResultStatusSucceeded:
		return "Succeeded"
	case ResultStatusFailed:
		return "Failed"
	case ResultStatusSkipped:
		return "Skipped"
	case ResultStatusDeferred:
		return "Deferred"
	case ResultStatusRequiresAuth:
		return "RequiresAuth"
	default:
		return fmt.Sprintf("ResultStatusType(%d)", int(s))
	}
}

func ParseResultStatus(s string) (ResultStatusType, error) {
	switch strings.ToLower(s) {
	case "pending":
		return ResultStatusPending, nil
	case "succeeded":
		return ResultStatusSucceeded, nil
	case "failed":
		return ResultStatusFailed, nil
	case "skipped":
		return ResultStatusSkipped, nil
	case "deferred":
		return ResultStatusDeferred, nil
	case "requiresauth":
		return ResultStatusRequiresAuth, nil
	default:
		return 0, fmt.Errorf("invalid ResultStatus: %s", s)
	}
}
