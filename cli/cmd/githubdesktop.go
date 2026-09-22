package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
)

// runGitHubDesktop registers the current working directory with GitHub
// Desktop in one shot, or runs optimize-projects / clear.
func runGitHubDesktop(args []string) error {
	checkHelp(constants.CmdGitHubDesktop, args)
	if isGHDesktopSubcommand(args) {
		return dispatchGHDesktopSubcommand(args)
	}

	if err := ensureGHDesktopInstalled(args); err != nil {
		return err
	}

	return processGHDesktopTarget(args)
}

func isGHDesktopSubcommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch strings.ToLower(args[0]) {
	case "optimize-projects", "optimize", "--repeat-fix", "-r", "dedupe",
		"clear", "clean", "group", "groups", "grp":
		return true
	}

	return false
}

func dispatchGHDesktopSubcommand(args []string) error {
	switch strings.ToLower(args[0]) {
	case "optimize-projects", "optimize", "--repeat-fix", "-r", "dedupe":
		return runGitHubDesktopOptimize(args[1:])
	case "clear", "clean":
		return runGitHubDesktopClear(args[1:])
	default:
		return runGitHubDesktopGroup(args[1:])
	}
}

// hasGHDesktopInstallFlag reports whether --install or -i was passed.
func hasGHDesktopInstallFlag(args []string) bool {
	for _, arg := range args {
		if arg == constants.FlagGHDesktopInstall || arg == constants.FlagGHDesktopInstallShort {
			return true
		}
	}

	return false
}

// ensureGHDesktopInstalled auto-installs GitHub Desktop if --install/-i is present and CLI is missing.
func ensureGHDesktopInstalled(args []string) error {
	hasInstall := hasGHDesktopInstallFlag(args)
	if !hasInstall {
		return nil
	}

	cli := desktop.ResolveCLI()
	if cli != "" {
		return nil
	}

	return cmdinstall.RunInstall([]string{constants.ToolGitHubDesktop, "--yes"})
}

func processGHDesktopTarget(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrGHDesktopCwd)
	}

	target := resolveGHDesktopTarget(cwd, args)
	isRepo := isGitRepo(target)
	if !isRepo {
		RenderGitHubDesktopHelp()
		return apperror.NewValidationError(fmt.Sprintf(constants.ErrGHDesktopNotRepo, target))
	}

	return registerGHDesktop(target)
}

// resolveGHDesktopTarget returns the absolute path to register: cwd by
// default, or args[0] if the user passed an explicit path.
func resolveGHDesktopTarget(cwd string, args []string) string {
	for _, arg := range args {
		if !isFlagToken(arg) {
			return toAbsPath(arg)
		}
	}

	return cwd
}

func toAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
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
func registerGHDesktop(target string) error {
	cli := desktop.ResolveCLI()
	if cli == "" {
		desktop.PrintInstallSuggestions()
		return desktop.NewMissingCLIError()
	}

	return invokeGHDesktop(cli, target)
}

func invokeGHDesktop(cli, target string) error {
	fmt.Printf(constants.MsgGHDesktopRegister, target)
	cmd := exec.Command(cli, target)
	if _, runErr := cmd.CombinedOutput(); runErr != nil {
		return apperror.NewSimple(constants.ErrGHDesktopInvoke, "E9000")
	}

	fmt.Printf(constants.MsgGHDesktopDone, target)

	return nil
}
