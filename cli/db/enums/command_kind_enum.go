package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// CommandKindType is a strongly-typed enum.
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
	CommandKindBashCommand    CommandKindType = 13
	CommandKindShellCommand   CommandKindType = 14
)

// IsValid reports whether the value matches a known variant.
func (c CommandKindType) IsValid() bool {
	switch c {
	case CommandKindPsCommand, CommandKindCmdCommand, CommandKindInstall, CommandKindGitPull, CommandKindGitPush, CommandKindGitCommit, CommandKindGitStatus, CommandKindProjRun, CommandKindProjCreateCICD, CommandKindRestart, CommandKindShutdown, CommandKindLogoff, CommandKindBashCommand, CommandKindShellCommand:
		return true
	default:
		return false
	}
}

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
	case CommandKindBashCommand:
		return "BashCommand"
	case CommandKindShellCommand:
		return "ShellCommand"
	default:
		return fmt.Sprintf("CommandKindType(%d)", int(k))
	}
}

func ParseCommandKind(s string) (CommandKindType, *apperror.AppError) {
	switch strings.ToLower(s) {
	case "pscommand":
		return CommandKindPsCommand, nil
	case "cmdcommand":
		return CommandKindCmdCommand, nil
	case "bashcommand", "bash":
		return CommandKindBashCommand, nil
	case "shellcommand", "shell", "sh":
		return CommandKindShellCommand, nil
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
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid CommandKind: %s", s))
	}
}

// AllCommandKindTypes returns all valid variants.
func AllCommandKindTypes() []CommandKindType {
	return []CommandKindType{
		CommandKindPsCommand,
		CommandKindCmdCommand,
		CommandKindInstall,
		CommandKindGitPull,
		CommandKindGitPush,
		CommandKindGitCommit,
		CommandKindGitStatus,
		CommandKindProjRun,
		CommandKindProjCreateCICD,
		CommandKindRestart,
		CommandKindShutdown,
		CommandKindLogoff,
		CommandKindBashCommand,
		CommandKindShellCommand,
	}
}

// IsPsCommand reports whether the value is CommandKindPsCommand.
func (c CommandKindType) IsPsCommand() bool {
	return c == CommandKindPsCommand
}

// IsCmdCommand reports whether the value is CommandKindCmdCommand.
func (c CommandKindType) IsCmdCommand() bool {
	return c == CommandKindCmdCommand
}

// IsInstall reports whether the value is CommandKindInstall.
func (c CommandKindType) IsInstall() bool {
	return c == CommandKindInstall
}

// IsGitPull reports whether the value is CommandKindGitPull.
func (c CommandKindType) IsGitPull() bool {
	return c == CommandKindGitPull
}

// IsGitPush reports whether the value is CommandKindGitPush.
func (c CommandKindType) IsGitPush() bool {
	return c == CommandKindGitPush
}

// IsGitCommit reports whether the value is CommandKindGitCommit.
func (c CommandKindType) IsGitCommit() bool {
	return c == CommandKindGitCommit
}

// IsGitStatus reports whether the value is CommandKindGitStatus.
func (c CommandKindType) IsGitStatus() bool {
	return c == CommandKindGitStatus
}

// IsProjRun reports whether the value is CommandKindProjRun.
func (c CommandKindType) IsProjRun() bool {
	return c == CommandKindProjRun
}

// IsProjCreateCICD reports whether the value is CommandKindProjCreateCICD.
func (c CommandKindType) IsProjCreateCICD() bool {
	return c == CommandKindProjCreateCICD
}

// IsRestart reports whether the value is CommandKindRestart.
func (c CommandKindType) IsRestart() bool {
	return c == CommandKindRestart
}

// IsShutdown reports whether the value is CommandKindShutdown.
func (c CommandKindType) IsShutdown() bool {
	return c == CommandKindShutdown
}

// IsLogoff reports whether the value is CommandKindLogoff.
func (c CommandKindType) IsLogoff() bool {
	return c == CommandKindLogoff
}

// IsBashCommand reports whether the value is CommandKindBashCommand.
func (c CommandKindType) IsBashCommand() bool {
	return c == CommandKindBashCommand
}

// IsShellCommand reports whether the value is CommandKindShellCommand.
func (c CommandKindType) IsShellCommand() bool {
	return c == CommandKindShellCommand
}

// HasKind reports whether the value is non-empty and valid.
func (c CommandKindType) HasKind() bool {
	return c.IsValid()
}
