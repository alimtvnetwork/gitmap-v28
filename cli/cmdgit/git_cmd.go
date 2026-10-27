package cmdgit

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// ExecGitInspection executes a read-only git command in targetDir with inherited stdio.
// It enforces safe subprocess environments and disables interactive paging.
func ExecGitInspection(targetDir string, gitArgs ...string) error {
	finalArgs := append([]string{"--no-pager"}, gitArgs...)
	cmd := exec.Command(constants.GitBin, finalArgs...)

	if targetDir != "" {
		cmd.Dir = targetDir
	}

	cmd.Env = gitutil.BuildSafeGitEnv("GIT_PAGER=cat", "PAGER=cat")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 128 {
			return apperror.NewWithDetails(
				"git",
				"E9001",
				"git execution failed: "+strings.Join(gitArgs, " "),
				"cli",
				apperror.ErrorTypeAbort,
				apperror.SeverityWarn,
				map[string]any{"reported": true},
			)
		}

		return apperror.WrapSimple(err, "git "+strings.Join(gitArgs, " "))
	}

	return nil
}

// ParseGitDirAndArgs extracts an optional -C <dir> or --dir <dir> flag from args.
func ParseGitDirAndArgs(args []string) (string, []string) {
	var targetDir string
	var remaining []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-C" && i+1 < len(args) {
			targetDir = args[i+1]
			i++
			continue
		}

		if strings.HasPrefix(arg, "--dir=") {
			targetDir = strings.TrimPrefix(arg, "--dir=")
			continue
		}

		if arg == "--dir" && i+1 < len(args) {
			targetDir = args[i+1]
			i++
			continue
		}

		remaining = append(remaining, arg)
	}

	return targetDir, remaining
}
