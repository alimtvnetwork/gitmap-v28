package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
)

// runGitHubDesktop registers the current working directory with GitHub
// Desktop in one shot. Unlike `desktop-sync` (which walks last-scan output),
// this command requires no prior scan — it just verifies cwd is a git repo
// and invokes the GitHub Desktop CLI on it.
func runGitHubDesktop(args []string) error {
	checkHelp(constants.CmdGitHubDesktop, args)

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.Wrap(err, constants.ErrGHDesktopCwd, nil)
	}

	target := resolveGHDesktopTarget(cwd, args)
	isNonGitRepo := !isGitRepo(target)
	if isNonGitRepo {
		return apperror.New(constants.ErrGHDesktopNotRepo, "E9000", nil)
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
		apperror.New(constants.MsgDesktopNotFound, "E9000", nil)
		return
	}

	fmt.Printf(constants.MsgGHDesktopRegister, target)
	cmd := exec.Command(cli, target)
	_, runErr := cmd.CombinedOutput()
	if runErr != nil {
		apperror.New(constants.ErrGHDesktopInvoke, "E9000", nil)
		return
	}

	fmt.Printf(constants.MsgGHDesktopDone, target)
}
