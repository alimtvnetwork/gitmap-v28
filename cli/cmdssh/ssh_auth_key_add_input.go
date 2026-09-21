package cmdssh

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func readSingleArgFile(path string) (string, bool) {
	data, err := os.ReadFile(path)
	hasErr := err != nil
	if hasErr {
		return "", false
	}

	return string(data), true
}

func readArgOrFile(args []string) string {
	isSingle := len(args) == 1
	if !isSingle {
		return strings.Join(args, " ")
	}

	content, isFile := readSingleArgFile(args[0])
	if isFile {
		return content
	}

	return args[0]
}

func promptKeyInput() (string, error) {
	isInteractive := isInteractiveTerminal()
	if !isInteractive {
		return "", apperror.NewValidationError("missing public key argument")
	}

	fmt.Print("Paste your SSH public key: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", apperror.WrapSimple(err, "promptKeyInput.Read")
	}

	return line, nil
}

func readPublicKeyInput(args []string) (string, error) {
	hasArgs := len(args) > 0
	if hasArgs {
		return readArgOrFile(args), nil
	}

	return promptKeyInput()
}
