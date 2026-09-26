package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runCreate(args []string) error {
	checkHelp("create", args)
	isNoArgs := len(args) == 0
	if isNoArgs && hasFilesInCurrentDir() {
		args = []string{"."}
		isNoArgs = false
	}
	isNonInteractive := !isInteractiveStdin()
	isHeadlessError := isNoArgs && isNonInteractive
	if isHeadlessError {
		return apperror.NewSimple("usage: gitmap repo-create (repoc) <name> [folder] [slug] [flags]", "E1076")
	}

	subArgs, resolveErr := resolveCreateArgs(args)
	if resolveErr != nil {
		return resolveErr
	}

	return executeCreateRepo(subArgs, false)
}

func runCreateLocal(args []string) error {
	checkHelp("create", args)
	isNoArgs := len(args) == 0
	isNonInteractive := !isInteractiveStdin()
	isHeadlessError := isNoArgs && isNonInteractive
	if isHeadlessError {
		return apperror.NewSimple("usage: gitmap create-local-repo (clr) <name> [folder] [slug] [flags]", "E1076")
	}

	subArgs, resolveErr := resolveCreateArgs(args)
	if resolveErr != nil {
		return resolveErr
	}

	return executeCreateRepo(subArgs, true)
}

func resolveCreateArgs(args []string) ([]string, error) {
	subArgs := normalizeCreateArgs(args)
	if len(subArgs) > 0 {
		return subArgs, nil
	}

	name, promptErr := promptRepoName()
	if promptErr != nil {
		return nil, promptErr
	}

	return []string{name}, nil
}

func normalizeCreateArgs(args []string) []string {
	if len(args) > 0 && (args[0] == "repo" || args[0] == "repository") {
		return args[1:]
	}

	return args
}

func promptRepoName() (string, error) {
	fmt.Printf("\n  %s● Create New Git Repository%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Print("Enter repository name: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", apperror.WrapSimple(err, "read repo name:")
	}

	name := strings.TrimSpace(line)
	if len(name) == 0 {
		return "", apperror.NewSimple("repository name cannot be empty", "E1077")
	}

	return name, nil
}

func hasFilesInCurrentDir() bool {
	entries, err := os.ReadDir(".")
	if err != nil {
		return false
	}
	for _, e := range entries {
		name := e.Name()
		if name != ".git" && !strings.HasPrefix(name, ".") {
			return true
		}
	}

	return false
}
