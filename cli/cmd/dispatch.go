package cmd

import (
	"os"
	"strings"
)

func init() {
	hookPathDispatch()
}

func hookPathDispatch() {
	if len(os.Args) >= 2 && isPathCommand(os.Args[1]) {
		if err := RunPathCLI(os.Args[2:]); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func isPathCommand(cmd string) bool {
	clean := strings.ToLower(strings.TrimSpace(cmd))
	return clean == "path"
}

// DispatchPath handles programmatic invocation of path commands.
func DispatchPath(args []string) error {
	return RunPathCLI(args)
}
