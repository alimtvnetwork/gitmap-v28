package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"golang.org/x/crypto/ssh"
)

func handleRemoteMacroAdd(client *ssh.Client, c db.SSHConnection, args []string) (bool, error) {
	targetName, isAddNoSteps := parseMacroAddTargetName(args)
	if isAddNoSteps {
		return executeMacroAutoDeploy(client, c, targetName)
	}

	return false, nil
}

func executeMacroAutoDeploy(client *ssh.Client, c db.SSHConnection, targetName string) (bool, error) {
	localMacro, err := macro.LoadMacro(targetName)
	if err != nil || localMacro == nil {
		return reportMissingMacroAdvice(c, targetName)
	}

	shell := determineFallbackShell(c.OS)
	isOk := syncSingleMacroToClient(client, c, *localMacro, shell, true)
	if isOk {
		msg := fmt.Sprintf("✓ Deployed local macro %q (%d step(s)) to remote node", localMacro.Name, len(localMacro.Steps))
		printNodeResultOutput(c.Alias, c.IPAddress, msg, nil)
		return true, nil
	}
	deployErr := fmt.Errorf("failed to deploy local macro %q to %s", targetName, c.Alias)
	printNodeResultOutput(c.Alias, c.IPAddress, "", deployErr)
	return true, deployErr
}

func reportMissingMacroAdvice(c db.SSHConnection, targetName string) (bool, error) {
	advice := fmt.Sprintf("macro %q not found locally. To create locally first: 'gitmap macro add %s'\n    Or provide steps: gitmap ssh exec \"gitmap macro add %s <cmd1> [cmd2...]\"", targetName, targetName, targetName)
	printNodeResultOutput(c.Alias, c.IPAddress, advice, nil)

	return true, nil
}

func parseMacroAddTargetName(args []string) (string, bool) {
	tokens := extractNormalizedTokens(args)
	if len(tokens) == 0 {
		return "", false
	}
	if tokens[0] == "gitmap" {
		tokens = tokens[1:]
	}

	return inspectMacroTokens(tokens)
}

func inspectMacroTokens(tokens []string) (string, bool) {
	if len(tokens) < 3 || tokens[0] != "macro" {
		return "", false
	}
	isAddVerb := tokens[1] == "add" || tokens[1] == "create" || tokens[1] == "new"
	if isAddVerb && isValidMacroName(tokens[2]) && hasOnlyFlagTokens(tokens[3:]) {
		return tokens[2], true
	}

	return "", false
}

func isValidMacroName(name string) bool {
	hasDashPrefix := strings.HasPrefix(name, "-")

	return !hasDashPrefix && name != ""
}

func hasOnlyFlagTokens(tokens []string) bool {
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		isValuedFlag := isFlagWithValue(t)
		if isValuedFlag && i+1 < len(tokens) {
			i++
			continue
		}
		isFlag := strings.HasPrefix(t, "-")
		if !isFlag {
			return false
		}
	}

	return true
}

func isFlagWithValue(flag string) bool {
	return flag == "--desc" || flag == "--description" || flag == "--tag"
}

func extractNormalizedTokens(args []string) []string {
	var tokens []string
	for _, a := range args {
		for _, part := range strings.Fields(a) {
			tokens = append(tokens, part)
		}
	}

	return tokens
}
