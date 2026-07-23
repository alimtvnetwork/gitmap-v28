// Package constants centralizes every user-facing string in the CLI.
// This file owns the `reinstall` top-level command surface.
//
// `reinstall` does NOT collide with `install` / `uninstall` (which are
// owned by the third-party tool installer in constants_install.go). It
// is a separate top-level verb that manages the gitmap binary itself,
// auto-detecting the install method:
//
//   - If a source repo is linked (constants.RepoPath != ""), shell out to
//     run.ps1 -reinstall (Windows) or run.sh --reinstall (Unix). This
//     preserves the user's existing build → deploy → setup pipeline.
//   - Otherwise, run self-uninstall + self-install back-to-back.
//
// The auto-detect is documented in helptext/reinstall.md and overridable
// with --mode {auto|repo|self}.
package constants

// gitmap:cmd top-level
// Reinstall command.
const (
	CmdReinstall = "reinstall"
)

// Reinstall flag names.
const (
	FlagReinstallMode = "mode"
	FlagReinstallYes  = "yes"
)

// Reinstall flag descriptions.
const (
	FlagDescReinstallMode = "Override auto-detect: 'auto' (default), 'repo' (force run.ps1/sh), 'self' (force self-uninstall+self-install)"
	FlagDescReinstallYes  = "Auto-confirm without prompting"
)

// Reinstall mode values.
const (
	ReinstallModeAuto = "auto"
	ReinstallModeRepo = "repo"
	ReinstallModeSelf = "self"
)

// Reinstall messages.
const (
	MsgReinstallHeader      = "\n  ╔══════════════════════════════════════╗\n  ║         gitmap reinstall             ║\n  ╚══════════════════════════════════════╝\n\n"
	MsgReinstallModeAuto    = "  → Mode: auto (detected: %s)\n"
	MsgReinstallModeForced  = "  → Mode: %s (forced via --mode)\n"
	MsgReinstallRepoPath    = "  → Source repo: %s\n"
	MsgReinstallRunningRepo = "  → Running %s -reinstall ...\n"
	MsgReinstallRunningSelf = "  → Running self-uninstall, then self-install ...\n"
	MsgReinstallStepUninst  = "\n  [1/2] self-uninstall\n"
	MsgReinstallStepInst    = "\n  [2/2] self-install\n"
	MsgReinstallDone        = "\n  ✓ Reinstall complete.\n\n"
	MsgReinstallConfirm     = "  → Proceed with reinstall? (yes/N): "
)

// Reinstall errors.
const (
	ErrReinstallUnknownMode    = "Error: unknown --mode value %q (allowed: auto|repo|self)\n"
	ErrReinstallNoRepoLinked   = "Error: --mode=repo requested but no source repo is linked (RepoPath empty). Run 'gitmap set-source-repo' first or use --mode=self.\n"
	ErrReinstallScriptNotFound = "Error: reinstall script %s not found at %s (operation: stat, reason: file does not exist)\n"
	ErrReinstallScriptFailed   = "Error: reinstall script %s exited with code %d\n"
	ErrReinstallSelfUninstall  = "Error: self-uninstall step failed: %v\n"
	ErrReinstallSelfInstall    = "Error: self-install step failed: %v\n"
	ErrReinstallAborted        = "  ✗ Reinstall canceled by user.\n"
)
