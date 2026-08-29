package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
)

func ParseSubCommands(tokens []string) ([]cluster.ClusterSubCommand, error) {
	var subCmds []cluster.ClusterSubCommand
	isEmptyTokens := len(tokens) == 0
	if isEmptyTokens {
		return nil, nil
	}

	var currentTokens []string

	commitCurrent := func() error {
		isEmpty := len(currentTokens) == 0
		if isEmpty {
			return nil
		}

		cmdToken := strings.ToLower(currentTokens[0])
		var kind db.CommandKindType
		var rawArgParts []string

		isGit := cmdToken == "git"
		isProj := cmdToken == "proj"
		hasMinTokens := len(currentTokens) >= 2

		if isGit && !hasMinTokens {
			return fmt.Errorf("missing git sub-command")
		}
		if isProj && !hasMinTokens {
			return fmt.Errorf("missing proj sub-command")
		}

		if isGit {
			subToken := strings.ToLower(currentTokens[1])
			switch subToken {
			case "pull":
				kind = db.CommandKindGitPull
			case "push":
				kind = db.CommandKindGitPush
			case "commit":
				kind = db.CommandKindGitCommit
			case "status":
				kind = db.CommandKindGitStatus
			default:
				return fmt.Errorf("unknown git sub-command: %s", subToken)
			}
			rawArgParts = currentTokens[2:]
		} else if isProj {
			subToken := strings.ToLower(currentTokens[1])
			switch subToken {
			case "run":
				kind = db.CommandKindProjRun
			case "create-cicd":
				kind = db.CommandKindProjCreateCICD
			default:
				return fmt.Errorf("unknown proj sub-command: %s", subToken)
			}
			rawArgParts = currentTokens[2:]
		} else {
			switch cmdToken {
			case "ps":
				kind = db.CommandKindPsCommand
			case "cmd":
				kind = db.CommandKindCmdCommand
			case "install":
				kind = db.CommandKindInstall
			case "restart":
				kind = db.CommandKindRestart
			case "shutdown":
				kind = db.CommandKindShutdown
			case "logoff":
				kind = db.CommandKindLogoff
			default:
				return fmt.Errorf("unknown sub-command token: %s", cmdToken)
			}
			rawArgParts = currentTokens[1:]
		}

		subCmds = append(subCmds, cluster.ClusterSubCommand{
			Kind:   kind,
			RawArg: strings.Join(rawArgParts, " "),
		})
		currentTokens = nil
		return nil
	}

	for _, token := range tokens {
		isComma := token == ","
		var err error

		if isComma {
			err = commitCurrent()
		}
		if err != nil {
			return nil, err
		}
		if isComma {
			continue
		}

		hasCommaSuffix := strings.HasSuffix(token, ",")
		stripped := strings.TrimSuffix(token, ",")
		isStrippedEmpty := stripped == ""

		if hasCommaSuffix && !isStrippedEmpty {
			currentTokens = append(currentTokens, stripped)
		}

		if hasCommaSuffix {
			err = commitCurrent()
		}

		if err != nil {
			return nil, err
		}
		if hasCommaSuffix {
			continue
		}

		currentTokens = append(currentTokens, token)
	}

	err := commitCurrent()
	hasErr := err != nil
	if hasErr {
		return nil, err
	}

	return subCmds, nil
}
