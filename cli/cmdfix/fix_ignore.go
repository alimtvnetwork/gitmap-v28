package cmdfix

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
)

func isFixIgnoreRequest(args []string) bool {
	for _, a := range args {
		if isIgnoreRoutingKeyword(strings.ToLower(a)) {
			return true
		}
	}
	return false
}

func isIgnoreRoutingKeyword(low string) bool {
	switch low {
	case "ignore", "ignores", "gitignore", "fia", "fias", "fix-ignore-all", "fix-ignores-all":
		return true
	}
	return false
}

func isSSHArgKeyword(low string) bool {
	return low == "ssh" || low == "--ssh" || low == "fias"
}

func isIgnoreFilterKeyword(low string) bool {
	switch low {
	case "ignore", "ignores", "gitignore", "all", "fia", "fias", "ssh", "--ssh":
		return true
	}
	return false
}

func filterIgnoreArgs(args []string) []string {
	var remaining []string
	for _, a := range args {
		low := strings.ToLower(a)
		if isIgnoreFilterKeyword(low) {
			continue
		}
		remaining = append(remaining, a)
	}
	return remaining
}

func isSSHRequested(args []string) bool {
	for _, a := range args {
		if isSSHArgKeyword(strings.ToLower(a)) {
			return true
		}
	}
	return false
}

func runFixIgnoreDispatch(args []string) error {
	hasSSH := isSSHRequested(args)
	remaining := filterIgnoreArgs(args)
	if hasSSH {
		return runFixIgnoreSSH(remaining)
	}
	return runFixIgnoreLocal(remaining)
}

func runFixIgnoreSSH(args []string) error {
	appErr := cmdignore.RunFixIgnoresAllSSH(args)
	if appErr != nil {
		return appErr
	}
	return nil
}

func runFixIgnoreLocal(args []string) error {
	appErr := cmdignore.RunFixIgnoreAll(args)
	if appErr != nil {
		return appErr
	}
	return nil
}
