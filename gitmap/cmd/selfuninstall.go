package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// selfUninstallOpts holds parsed flags for self-uninstall.
//
// ShellMode mirrors the self-install flag (v3.49.0+): it selects which
// shell-profile families to strip the PATH snippet from. Accepts the
// same singletons (auto|both|zsh|bash|pwsh|fish) and `+`-joined combos
// (e.g. zsh+pwsh) as `gitmap self-install`. Default `auto` strips every
// known profile (safest for full removal).
type selfUninstallOpts struct {
	Confirm     bool
	KeepData    bool
	KeepSnippet bool
	ShellMode   string
}

// runSelfUninstall is the entry point for `gitmap self-uninstall`.
// On Windows the running .exe is locked, so we copy ourselves to a temp
// path and re-exec the hidden self-uninstall-runner from there.
func runSelfUninstall(args []string) {
	checkHelp(constants.CmdSelfUninstall, args)
	opts := parseSelfUninstallFlags(args)
	if !opts.Confirm && !confirmSelfUninstall(opts) {
		fmt.Fprint(os.Stderr, constants.ErrSelfUninstallNoConfirm)
		os.Exit(1)
	}
	if shouldHandoffSelfUninstall() {
		handoffSelfUninstall(opts, args)
		return
	}
	executeSelfUninstall(opts)
}

// parseSelfUninstallFlags reads --confirm / --keep-data / --keep-snippet
// / --shell-mode (canonical) / --profile + --dual-shell (hidden aliases).
//
// The shell-mode resolver/validator is shared with self-install so both
// commands accept identical syntax — see resolveShellMode and
// validateShellMode in selfinstall.go.
func parseSelfUninstallFlags(args []string) selfUninstallOpts {
	fs := flag.NewFlagSet(constants.CmdSelfUninstall, flag.ExitOnError)
	opts := selfUninstallOpts{}
	var (
		shellModeFlag string
		profileFlag   string
		dualShell     bool
	)
	fs.BoolVar(&opts.Confirm, "confirm", false, constants.FlagDescSelfConfirm)
	fs.BoolVar(&opts.KeepData, "keep-data", false, constants.FlagDescSelfKeepData)
	fs.BoolVar(&opts.KeepSnippet, "keep-snippet", false, constants.FlagDescSelfKeepSnippet)
	fs.StringVar(&shellModeFlag, "shell-mode", "", constants.FlagDescSelfShellMode)
	fs.StringVar(&profileFlag, "profile", "", constants.FlagDescSelfProfile)
	fs.BoolVar(&dualShell, "dual-shell", false, constants.FlagDescSelfDualShell)
	fs.Parse(reorderFlagsBeforeArgs(args))
	opts.ShellMode = resolveShellMode(shellModeFlag, profileFlag, dualShell)
	validateShellMode(opts.ShellMode)

	return opts
}

// confirmSelfUninstall prints the target list and prompts for "yes".
func confirmSelfUninstall(opts selfUninstallOpts) bool {
	printSelfUninstallTargets(opts)
	fmt.Print(constants.MsgSelfUninstallConfirmPrompt)
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		return false
	}

	return answer == "yes"
}

// printSelfUninstallTargets prints what self-uninstall will remove. The
// snippet line expands to one entry per resolved profile so the user
// sees exactly which files --shell-mode is about to touch.
func printSelfUninstallTargets(opts selfUninstallOpts) {
	fmt.Print(constants.MsgSelfUninstallHeader)
	fmt.Print(constants.MsgSelfUninstallTargets)
	fmt.Printf(constants.MsgSelfUninstallTargetBin, selfDeployDir())
	fmt.Printf(constants.MsgSelfUninstallTargetData, selfDataDir())
	for _, p := range resolveProfilesForShellMode(opts.ShellMode) {
		fmt.Printf(constants.MsgSelfUninstallTargetSnippet, p)
	}
	fmt.Printf(constants.MsgSelfUninstallTargetCompl, selfDeployDir())
}

// executeSelfUninstall removes each target the user did not opt out of.
// PATH snippets are stripped from every profile resolved by --shell-mode
// (deterministic; mirrors how self-install picked its write targets).
func executeSelfUninstall(opts selfUninstallOpts) {
	if !opts.KeepSnippet {
		for _, p := range resolveProfilesForShellMode(opts.ShellMode) {
			removeProfileSnippet(p)
		}
	}
	removeCompletionSourceLines()
	removeCompletionFiles(selfDeployDir())
	if !opts.KeepData {
		removePathBestEffort(selfDataDir())
	}
	removeDeployArtifacts(selfDeployDir())
	fmt.Print(constants.MsgSelfUninstallDone)
}

// shouldHandoffSelfUninstall reports whether the running binary lives
// inside the directory we are about to delete (Windows only — on Unix
// we can unlink an open file safely).
func shouldHandoffSelfUninstall() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	self, err := os.Executable()
	if err != nil {
		return false
	}
	deploy := selfDeployDir()
	if len(deploy) == 0 {
		return false
	}

	return strings.HasPrefix(filepath.Clean(self), filepath.Clean(deploy))
}

// runSelfUninstallRunner is the hidden command run by the temp-copy
// handoff. It performs the actual removal then deletes itself.
func runSelfUninstallRunner() {
	opts := parseSelfUninstallFlags(os.Args[2:])
	executeSelfUninstall(opts)
	scheduleSelfDelete()
}
