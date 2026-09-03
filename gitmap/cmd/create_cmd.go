package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runCreate(args []string) error {
	checkHelp("create", args)
	if len(args) == 0 && !isInteractiveStdin() {
		return apperror.NewSimple("usage: gitmap create [repo] <name> [flags]", "E1076")
	}
	subArgs, resolveErr := resolveCreateArgs(args)
	if resolveErr != nil {
		return resolveErr
	}
	return executeCreateRepo(subArgs)
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
	if len(args) > 0 && args[0] == "repo" {
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
