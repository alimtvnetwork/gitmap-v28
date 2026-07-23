package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
)

// runGitHubDesktop registers the current working directory with GitHub
// Desktop in one shot. Unlike `desktop-sync` (which walks last-scan output),
// this command requires no prior scan — it just verifies cwd is a git repo
// and invokes the GitHub Desktop CLI on it.
func runGitHubDesktop(args []string) {
	checkHelp(constants.CmdGitHubDesktop, args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGHDesktopCwd, err)
		os.Exit(1)
	}

	target := resolveGHDesktopTarget(cwd, args)
	if !isGitRepo(target) {
		fmt.Fprintf(os.Stderr, constants.ErrGHDesktopNotRepo, target)
		os.Exit(1)
	}

	registerGHDesktop(target)
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
		fmt.Fprintln(os.Stderr, constants.MsgDesktopNotFound)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgGHDesktopRegister, target)
	cmd := exec.Command(cli, target)
	output, runErr := cmd.CombinedOutput()
	if runErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGHDesktopInvoke, runErr, output)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgGHDesktopDone, target)
}
