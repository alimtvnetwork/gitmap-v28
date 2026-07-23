package constants

// Self-install / self-uninstall command messages, errors, and defaults.
//
// These constants back the `gitmap self-install` and `gitmap self-uninstall`
// commands, which manage the gitmap binary itself (as opposed to the
// `install`/`uninstall` commands that manage third-party tools).
//
// Spec: spec/01-app/90-self-install-uninstall.md

// Default install directories per platform.
const (
	SelfInstallDefaultWindows = "D:\\gitmap"
	SelfInstallDefaultUnix    = ".local/bin/gitmap" // joined under $HOME at runtime
)

// Embedded script names.
const (
	SelfInstallScriptPwsh = "install.ps1"
	SelfInstallScriptBash = "install.sh"
)

// Remote installer URLs (fallback when embedded scripts are missing).
const (
	SelfInstallRemotePwsh = "https://raw.githubusercontent.com/" +
		"alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.ps1"
	SelfInstallRemoteBash = "https://raw.githubusercontent.com/" +
		"alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.sh"
)

// Self-install messages.
const (
	MsgSelfInstallHeader   = "\n  gitmap self-install\n\n"
	MsgSelfInstallPrompt   = "  Install directory [%s]: "
	MsgSelfInstallUsing    = "  Using install directory: %s\n"
	MsgSelfInstallEmbedded = "  Running embedded installer (%s)...\n"
	MsgSelfInstallRemote   = "  Embedded installer unavailable; downloading from %s\n"
	MsgSelfInstallDone     = "  ✓ Install completed.\n"
	MsgSelfInstallReminder = "  Open a new terminal (or reload your profile) to pick up PATH changes.\n"
)

// Self-install errors.
const (
	ErrSelfInstallScriptWrite = "Error: write installer to temp: %v\n"
	ErrSelfInstallScriptRun   = "Error: run installer: %v\n"
	ErrSelfInstallDownload    = "Error: download installer from %s: %v\n"
	ErrSelfInstallNoShell     = "Error: no supported shell found (need PowerShell on Windows or bash on Unix)\n"
	ErrSelfInstallReadStdin   = "Error: read install dir from stdin: %v\n"
)

// Self-uninstall messages.
const (
	MsgSelfUninstallHeader        = "\n  gitmap self-uninstall\n\n"
	MsgSelfUninstallTargets       = "  The following will be removed:\n"
	MsgSelfUninstallTargetBin     = "    - Binary + deploy dir: %s\n"
	MsgSelfUninstallTargetData    = "    - Data dir:            %s\n"
	MsgSelfUninstallTargetSnippet = "    - PATH snippet from:   %s\n"
	MsgSelfUninstallTargetCompl   = "    - Completion files in: %s\n"
	MsgSelfUninstallConfirmPrompt = "\n  Type 'yes' to proceed: "
	MsgSelfUninstallSkipBin       = "  ⚠ Could not resolve own binary location: %v\n"
	MsgSelfUninstallRemovedBin    = "  ✓ Removed binary: %s\n"
	MsgSelfUninstallRemovedDir    = "  ✓ Removed dir:    %s\n"
	MsgSelfUninstallSnippetGone   = "  ✓ PATH snippet removed from %s\n"
	MsgSelfUninstallSnippetMiss   = "  - No PATH snippet found in %s\n"
	MsgSelfUninstallDone          = "\n  ✓ gitmap has been uninstalled. Restart your terminal to clear $env:Path.\n\n"
	MsgSelfUninstallHandoffActive = "  Handing off to %s so the original binary can self-delete...\n"
)

// Self-uninstall errors.
const (
	ErrSelfUninstallNoConfirm    = "Error: refusing to run without --confirm or interactive 'yes'.\n"
	ErrSelfUninstallRemove       = "Error: remove %s: %v\n"
	ErrSelfUninstallSnippetRead  = "Error: read profile %s: %v\n"
	ErrSelfUninstallSnippetWrite = "Error: rewrite profile %s: %v\n"
	ErrSelfUninstallHandoffCopy  = "Error: create handoff copy: %v\n"
)

// Hidden runner subcommand for the self-uninstall handoff (lets the temp
// copy delete the original .exe on Windows where the running file is
// locked).
const CmdSelfUninstallRunner = "self-uninstall-runner" // gitmap:cmd skip

// Flag names shared by self-install / self-uninstall.
//
// FlagSelfShellMode is the canonical name as of v3.48.0. FlagSelfProfile
// (introduced in v3.46.0) and FlagSelfDualShell (v3.43.0) are kept as
// hidden aliases so existing scripts and CI continue to work — the Go
// parser collapses all three onto opts.ShellMode (formerly opts.Profile).
const (
	FlagSelfDir         = "--dir"
	FlagSelfYes         = "--yes"
	FlagSelfConfirm     = "--confirm"
	FlagSelfKeepData    = "--keep-data"
	FlagSelfKeepSnippet = "--keep-snippet"
	FlagSelfFromVersion = "--version"
	FlagSelfShellMode   = "--shell-mode" // canonical (v3.48.0+)
	FlagSelfProfile     = "--profile"    // hidden alias (v3.46.0+)
	FlagSelfDualShell   = "--dual-shell" // hidden alias (v3.43.0+)
	FlagSelfShowPath    = "--show-path"
	FlagSelfForceLock   = "--force-lock"
)

// ShellMode values accepted by --shell-mode. `auto` is the default
// (run detect_active_pwsh + $SHELL heuristics). `both` writes PATH to
// every supported shell's profile file. Singletons restrict writes to
// exactly one family.
//
// Combos: any `+`-joined combination of the singletons (e.g. `zsh+pwsh`,
// `bash+fish`, `zsh+bash+pwsh`) is also accepted. Combos are STRICT —
// only the listed families receive the PATH snippet; ~/.profile and any
// undeclared family are skipped. This is the difference from `both`,
// which writes everything detected.
//
// Aliases retained for back-compat:
//
//	ProfileMode* names mirror ShellMode* names — both are still exported
//	so old code referencing ProfileModeBoth keeps compiling.
const (
	ShellModeAuto = "auto"
	ShellModeBoth = "both"
	ShellModeZsh  = "zsh"
	ShellModeBash = "bash"
	ShellModePwsh = "pwsh"
	ShellModeFish = "fish"

	// Back-compat alias names — same values, different identifiers.
	ProfileModeAuto = ShellModeAuto
	ProfileModeBoth = ShellModeBoth
	ProfileModeZsh  = ShellModeZsh
	ProfileModeBash = ShellModeBash
	ProfileModePwsh = ShellModePwsh
	ProfileModeFish = ShellModeFish

	// ShellModeComboSep is the delimiter for combo modes like "zsh+pwsh".
	ShellModeComboSep = "+"
)

// SelfInstallShellModes lists the singleton (non-combo) values accepted
// by --shell-mode. Combos are validated by splitting on ShellModeComboSep
// and checking each token against this list.
//
// SelfInstallProfileModes is the back-compat alias kept for existing
// references in selfinstall.go and tests.
// Single-line initializer (NOT a `var (...)` block) so the constants-collision
// linter does not misread the indented `ShellMode*,` element lines as new
// top-level var declarations duplicating the const names above.
var SelfInstallShellModes = []string{ShellModeAuto, ShellModeBoth, ShellModeZsh, ShellModeBash, ShellModePwsh, ShellModeFish}

// SelfInstallProfileModes is the back-compat alias kept for existing
// references in selfinstall.go and tests.
var SelfInstallProfileModes = SelfInstallShellModes

// ErrSelfInstallShellModeInvalid fires when --shell-mode gets an unknown
// value (singleton or combo token). Format: %q = full bad value, %s =
// accepted singletons (combos are documented separately in help text to
// avoid an unbounded "valid values" string).
const ErrSelfInstallShellModeInvalid = "Error: --shell-mode %q is not valid. " +
	"Accepted singletons: %s. Combos: any '+'-joined combination of singletons " +
	"(e.g. zsh+pwsh, bash+fish, zsh+bash+pwsh).\n"

// ErrSelfInstallProfileInvalid is the back-compat error for the legacy
// --profile flag. Same wording as before to keep CI grep patterns valid.
const ErrSelfInstallProfileInvalid = "Error: --profile %q is not valid. Accepted: %s\n"

// Flag descriptions.
const (
	FlagDescSelfDir         = "Install directory (prompted with default if omitted)"
	FlagDescSelfYes         = "Skip the install-directory prompt and accept the default"
	FlagDescSelfConfirm     = "Required for self-uninstall to actually remove files"
	FlagDescSelfKeepData    = "Preserve the .gitmap data dir during self-uninstall"
	FlagDescSelfKeepSnippet = "Leave the PATH snippet in shell profile during self-uninstall"
	FlagDescSelfFromVersion = "Pin a specific gitmap version to install (e.g. v3.0.0)"
	FlagDescSelfShellMode   = "Which shell profile(s) to write PATH into: auto|both|zsh|bash|pwsh|fish " +
		"or any '+'-joined combo such as zsh+pwsh (default: auto)"
	FlagDescSelfProfile   = "Deprecated alias for --shell-mode (hidden; still works)"
	FlagDescSelfDualShell = "Deprecated alias for --shell-mode both (hidden; still works)"
	FlagDescSelfShowPath  = "Print detected shell, chosen PATH target, and every profile file written"
	FlagDescSelfForceLock = "Bypass the duplicate-install guard (recover from a stale lock left by a crashed installer)"
)

// Self-install duplicate-install guard.
const (
	// SelfInstallLockName is the suffix passed to lockfile.Acquire. The
	// resulting file is os.TempDir()/gitmap-selfinstall.lock and prevents
	// two `gitmap self-install` processes from racing each other.
	SelfInstallLockName = "selfinstall"

	// ErrSelfInstallLockHeld is shown when another live gitmap self-install
	// process holds the lock. Includes its PID so the user can investigate.
	ErrSelfInstallLockHeld = "Error: another gitmap self-install is already running (pid=%d).\n" +
		"       Wait for it to finish, or pass --force-lock if it crashed.\n"

	// ErrSelfInstallLock is the catch-all for lock filesystem failures
	// (permission denied, disk full, etc.).
	ErrSelfInstallLock = "Error: could not acquire self-install lock: %v\n"
)
