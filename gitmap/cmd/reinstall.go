package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// reinstallOpts collects parsed flags for `gitmap reinstall`.
type reinstallOpts struct {
	Mode string
	Yes  bool
}

// runReinstall is the entry point for `gitmap reinstall`. It picks an
// install path (repo-script vs self-uninstall+self-install), prompts
// for confirmation unless --yes, then delegates.
func runReinstall(args []string) {
	checkHelp(constants.CmdReinstall, args)
	opts := parseReinstallFlags(args)
	mode, detected := resolveReinstallMode(opts.Mode)
	announceReinstallMode(opts.Mode, mode, detected)
	if !opts.Yes && !confirmReinstall() {
		fmt.Fprint(os.Stderr, constants.ErrReinstallAborted)
		os.Exit(1)
	}
	dispatchReinstall(mode)
	fmt.Print(constants.MsgReinstallDone)
}

// parseReinstallFlags reads --mode / --yes / -y.
func parseReinstallFlags(args []string) reinstallOpts {
	fs := flag.NewFlagSet(constants.CmdReinstall, flag.ExitOnError)
	opts := reinstallOpts{}
	fs.StringVar(&opts.Mode, constants.FlagReinstallMode, constants.ReinstallModeAuto, constants.FlagDescReinstallMode)
	fs.BoolVar(&opts.Yes, constants.FlagReinstallYes, false, constants.FlagDescReinstallYes)
	fs.BoolVar(&opts.Yes, "y", false, constants.FlagDescReinstallYes)
	fs.Parse(reorderFlagsBeforeArgs(args))

	return opts
}

// resolveReinstallMode picks the actual mode to run, given the user
// override and the linked-repo state. Returns (mode, autoDetected).
func resolveReinstallMode(override string) (string, bool) {
	switch override {
	case constants.ReinstallModeRepo:
		if len(constants.RepoPath) == 0 {
			fmt.Fprint(os.Stderr, constants.ErrReinstallNoRepoLinked)
			os.Exit(1)
		}
		return constants.ReinstallModeRepo, false
	case constants.ReinstallModeSelf:
		return constants.ReinstallModeSelf, false
	case constants.ReinstallModeAuto:
		if len(constants.RepoPath) > 0 {
			return constants.ReinstallModeRepo, true
		}
		return constants.ReinstallModeSelf, true
	default:
		fmt.Fprintf(os.Stderr, constants.ErrReinstallUnknownMode, override)
		os.Exit(1)
	}

	return constants.ReinstallModeSelf, true
}

// announceReinstallMode prints the resolved mode (and the linked repo
// path when applicable) so the user sees exactly what will happen.
func announceReinstallMode(rawOverride, mode string, detected bool) {
	fmt.Print(constants.MsgReinstallHeader)
	if detected || rawOverride == constants.ReinstallModeAuto {
		fmt.Printf(constants.MsgReinstallModeAuto, mode)
	} else {
		fmt.Printf(constants.MsgReinstallModeForced, mode)
	}
	if mode == constants.ReinstallModeRepo {
		fmt.Printf(constants.MsgReinstallRepoPath, constants.RepoPath)
	}
}

// confirmReinstall prompts for "yes" on stdin. Anything else aborts.
func confirmReinstall() bool {
	fmt.Print(constants.MsgReinstallConfirm)
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		return false
	}

	return strings.EqualFold(strings.TrimSpace(answer), "yes")
}

// dispatchReinstall hands off to the resolved implementation.
func dispatchReinstall(mode string) {
	if mode == constants.ReinstallModeRepo {
		executeReinstallRepo()

		return
	}
	executeReinstallSelf()
}

// executeReinstallRepo runs run.ps1 -reinstall (Windows) or
// run.sh --reinstall (Unix) from the linked source repo.
func executeReinstallRepo() {
	scriptPath, scriptName := pickReinstallScriptPath()
	if _, err := os.Stat(scriptPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReinstallScriptNotFound, scriptName, scriptPath)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgReinstallRunningRepo, scriptName)
	cmd := buildReinstallScriptCmd(scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = constants.RepoPath
	if err := cmd.Run(); err != nil {
		exitCode := 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, constants.ErrReinstallScriptFailed, scriptName, exitCode)
		os.Exit(exitCode)
	}
}

// pickReinstallScriptPath returns the absolute path + display name of
// run.ps1 (Windows) or run.sh (Unix) inside the linked repo.
func pickReinstallScriptPath() (string, string) {
	if runtime.GOOS == constants.OSWindows {
		name := "run.ps1"
		return filepath.Join(constants.RepoPath, name), name
	}
	name := "run.sh"

	return filepath.Join(constants.RepoPath, name), name
}

// buildReinstallScriptCmd assembles the platform-specific script invocation.
func buildReinstallScriptCmd(scriptPath string) *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return exec.Command("pwsh", "-ExecutionPolicy", "Bypass", "-NoProfile",
			"-NoLogo", "-File", scriptPath, "-reinstall")
	}

	return exec.Command("bash", scriptPath, "--reinstall")
}

// executeReinstallSelf runs self-uninstall (with --confirm so it does
// not re-prompt) followed by self-install (with --yes for the same
// reason). Failures in either step abort the reinstall.
func executeReinstallSelf() {
	fmt.Print(constants.MsgReinstallStepUninst)
	runSelfUninstall([]string{"--confirm"})
	fmt.Print(constants.MsgReinstallStepInst)
	runSelfInstall([]string{"--yes"})
}
