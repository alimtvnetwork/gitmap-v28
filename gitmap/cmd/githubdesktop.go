package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
)

// runGitHubDesktop registers the current working directory with GitHub
// Desktop in one shot, or runs optimize-projects / clear.
func runGitHubDesktop(args []string) error {
	checkHelp(constants.CmdGitHubDesktop, args)
	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "optimize-projects", "optimize", "--repeat-fix", "-r", "dedupe":
			return runGitHubDesktopOptimize(args[1:])
		case "clear", "clean":
			return runGitHubDesktopClear(args[1:])
		case "group", "groups", "grp":
			return runGitHubDesktopGroup(args[1:])
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrGHDesktopCwd)
	}

	target := resolveGHDesktopTarget(cwd, args)
	isNonGitRepo := !isGitRepo(target)
	if isNonGitRepo {
		return apperror.NewSimple(constants.ErrGHDesktopNotRepo, "E9000")
	}

	registerGHDesktop(target)
	return nil
}

// resolveGHDesktopTarget returns the absolute path to register: cwd by
// default, or args[0] if the user passed an explicit path.
func resolveGHDesktopTarget(cwd string, args []string) string {
	if len(args) == 0 {
		return cwd
	}

	abs, err := filepath.Abs(args[0])
	if err != nil {
		return args[0]
	}

	return abs
}

// isGitRepo reports whether dir contains a .git directory or file (worktrees
// use a .git file). Returns false on any stat error.
func isGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, constants.ExtGit))

	return err == nil
}

// registerGHDesktop verifies the GitHub Desktop CLI is on PATH, then invokes
// it with the target path. Exits non-zero on missing CLI or invocation error.
func registerGHDesktop(target string) {
	cli := desktop.ResolveCLI()
	if cli == "" {
		apperror.NewSimple(constants.MsgDesktopNotFound, "E9000")
		return
	}

	fmt.Printf(constants.MsgGHDesktopRegister, target)
	cmd := exec.Command(cli, target)
	_, runErr := cmd.CombinedOutput()
	if runErr != nil {
		apperror.NewSimple(constants.ErrGHDesktopInvoke, "E9000")
		return
	}

	fmt.Printf(constants.MsgGHDesktopDone, target)
}
