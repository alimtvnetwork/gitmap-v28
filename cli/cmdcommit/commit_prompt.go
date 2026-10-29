package cmdcommit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"golang.org/x/term"
)

var (
	promptStdinReader   io.Reader = os.Stdin
	promptStderr        io.Writer = os.Stderr
	isTerminalStdinFunc           = func() bool {
		return term.IsTerminal(int(os.Stdin.Fd()))
	}
)

// PromptCommitMessageIfMissing prompts the user via Stderr and Stdin if cleanArgs is empty.
func PromptCommitMessageIfMissing(cleanArgs []string, cmdUsage string, errCode string) ([]string, *apperror.AppError) {
	if len(cleanArgs) > 0 {
		return cleanArgs, nil
	}
	if !isTerminalStdinFunc() {
		return nil, apperror.NewSimple("Usage: gitmap "+cmdUsage+" <message>", errCode)
	}

	return readInteractiveCommitMessage(cmdUsage, errCode)
}

func readInteractiveCommitMessage(cmdUsage string, errCode string) ([]string, *apperror.AppError) {
	msg, err := ReadPromptInput("Enter commit message: ")
	if err != nil || msg == "" {
		return nil, apperror.NewSimple("Usage: gitmap "+cmdUsage+" <message>", errCode)
	}

	return []string{msg}, nil
}

// ReadPromptInput displays a prompt string to os.Stderr and reads a trimmed line from os.Stdin.
func ReadPromptInput(promptText string) (string, error) {
	fmt.Fprint(promptStderr, promptText)
	reader := bufio.NewReader(promptStdinReader)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

// PromptParameterIfMissing asks for a required parameter interactively when omitted.
func PromptParameterIfMissing(value, paramName, promptText string) (string, error) {
	if strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value), nil
	}
	if !isTerminalStdinFunc() {
		return "", fmt.Errorf("required parameter %s is missing in non-interactive environment", paramName)
	}

	return readRequiredPrompt(paramName, promptText)
}

func readRequiredPrompt(paramName, promptText string) (string, error) {
	val, err := ReadPromptInput(promptText)
	if err != nil {
		return "", err
	}
	if val == "" {
		return "", fmt.Errorf("parameter %s cannot be empty", paramName)
	}

	return val, nil
}
