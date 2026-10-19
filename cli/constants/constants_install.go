package constants

// gitmap:cmd top-level
// Install CLI commands.
const (
	CmdInstall        = "install"
	CmdInstallAlias   = "in"
	CmdUninstall      = "uninstall"
	CmdUninstallAlias = "un"
)

// Install help text.
const (
	HelpInstall   = "  install (in) <tool>    Install a developer tool by name"
	HelpUninstall = "  uninstall (un) <tool>  Remove a previously installed tool"
)

// Supported tool names — Core.
const (
	ToolVSCode             = "vscode"
	ToolNodeJS             = "node"
	ToolYarn               = "yarn"
	ToolBun                = "bun"
	ToolPnpm               = "pnpm"
	ToolPython             = "python"
	ToolGo                 = "go"
	ToolGit                = "git"
	ToolGitLFS             = "git-lfs"
	ToolGHCLI              = "gh"
	ToolGitHubDesktop      = "github-desktop"
	ToolCPP                = "cpp"
	ToolPHP                = "php"
	ToolPowerShell         = "powershell"
	ToolChocolatey         = "chocolatey"
	ToolWinget             = "winget"
	ToolNpp                = "npp"
	ToolNppSettings        = "npp-settings"
	ToolNppInstall         = "install-npp"
	ToolVSCodeSync         = "vscode-settings"
	ToolOBSSync            = "obs-settings"
	ToolWTSync             = "wt-settings"
	ToolScripts            = "scripts"
	ToolDbeaver            = "dbeaver"
	ToolStickyNotes        = "sticky-notes"
	ToolLiteDB             = "litedb"
	ToolVSCodeCtx          = "vscode-ctx"
	ToolPwshCtx            = "pwsh-ctx"
	ToolOBS                = "obs"
	ToolAllDevTools        = "all"
	ToolGitmapOneliner     = "gitmap-oneliner"
	ToolScriptsFixer       = "scripts-fixer"
	ToolCodingGuidelines   = "coding-guidelines"
	ToolMacroAhk           = "macro-ahk"
	ToolAgManager          = "ag-manager"
	ToolAgCtx              = "ag-ctx"
	ToolAntigravity        = "antigravity"
	ToolAntigravityIde     = "antigravity-ide"
	ToolAntigravityDesktop = "antigravity-desktop"
	ToolAgy                = "agy"
	ToolAntigravityCli     = "antigravity-cli"
	ToolChrome             = "chrome"
	ToolGoogleChrome       = "google-chrome"
	ToolBuildEssential     = "build-essential"
	ToolUbuntuCommon       = "ubuntu-common"
	ToolWordPress          = "wordpress"
	ToolGitmap             = "gitmap"
	ToolWpCli              = "wp-cli"
	ToolComposer           = "composer"
	ToolGitCompact         = "git-compact"
)

// Install flag names.
const (
	FlagInstallManager = "manager"
	FlagInstallVersion = "version"
	FlagInstallVerbose = "verbose"
	FlagInstallDryRun  = "dry-run"
	FlagInstallCheck   = "check"
	FlagInstallList    = "list"
	FlagInstallStatus  = "status"
	FlagInstallUpgrade = "upgrade"
	FlagInstallYes     = "yes"
	FlagInstallExplain = "explain"
)

// Install flag descriptions.
const (
	FlagDescInstallManager = "Force package manager (choco, winget, apt, brew, snap)"
	FlagDescInstallVersion = "Install a specific version"
	FlagDescInstallVerbose = "Show full installer output"
	FlagDescInstallDryRun  = "Show install command without executing"
	FlagDescInstallCheck   = "Only check if tool is installed"
	FlagDescInstallList    = "List all supported tools"
	FlagDescInstallStatus  = "Show installed tools from database"
	FlagDescInstallUpgrade = "Upgrade an already-installed tool"
	FlagDescInstallExplain = "ctx only: each generated menu entry prints '> gitmap <args>' to the terminal/notification before executing, so users see the exact resolved invocation"
	FlagDescInstallYes     = "Auto-confirm install without prompting"
)

// Uninstall flag names.
const (
	FlagUninstallDryRun = "dry-run"
	FlagUninstallForce  = "force"
	FlagUninstallPurge  = "purge"
)

// Uninstall flag descriptions.
const (
	FlagDescUninstallDryRun = "Show uninstall command without executing"
	FlagDescUninstallForce  = "Skip confirmation prompt"
	FlagDescUninstallPurge  = "Remove config files too"
)

// Install terminal messages.
const (
	MsgInstallChecking     = "\n  Checking if %s is installed...\n"
	MsgInstallFound        = "  ✓ %s is already installed (version: %s)\n"
	MsgInstallNotFound     = "  ✗ %s is not installed.\n"
	MsgInstallInstalling   = "\n  Installing %s...\n"
	MsgInstallSuccess      = "  ✓ %s installed successfully.\n"
	MsgInstallDryCmd       = "  [dry-run] Would run: %s\n"
	MsgInstallVerifying    = "\n  Verifying %s installation...\n"
	MsgInstallListHeader   = "Supported tools:\n\n"
	MsgInstallListRow      = "  %-20s %s\n"
	MsgInstallRecorded     = "  ✓ Recorded %s v%s in database.\n"
	MsgInstallStatusHdr    = "Installed tools:\n\n"
	MsgInstallStatusRow    = "  %-20s %-12s %-8s %s\n"
	MsgInstallListCategory = "\n%s\n%s\n"
	MsgInstallListGrouped  = "  %s %-22s %-10s %s\n"
	MsgInstallListLegend   = "\nLegend: ● installed   ○ not installed   ? unknown\n"
	StatusInstalled        = "●"
	StatusNotInstalled     = "○"
	StatusUnknown          = "?"
	MsgInstallExeVerify    = "  Verifying %s binary at: %s\n"
	MsgInstallExeFound     = "  ✓ Binary confirmed: %s\n"
	MsgInstallNppSettings  = "Syncing Notepad++ settings...\n"
	MsgInstallNppSkipBin   = "Skipping Notepad++ installation (settings-only mode)\n"
	MsgInstallNppSkipSet   = "Skipping Notepad++ settings (install-only mode)\n"
	MsgInstallNppExtract   = "Extracting Notepad++ settings to %s...\n"
	MsgInstallPrompt       = "\n  → Install %s %s using %s? (y/N): "
	MsgInstallPromptNoVer  = "\n  → Install %s (latest) using %s? (y/N): "
	MsgInstallAborted      = "\n  Installation canceled by user.\n"
	MsgInstallVersion      = "  → Version: %s\n"
	MsgInstallVersionLabel = "  → Version: latest\n"
	MsgInstallManager      = "  → Package manager: %s\n"
)

// Install error messages.
const (
	ErrInstallToolRequired    = "Tool name is required. Use --list to see available tools.\n"
	ErrInstallUnknownTool     = "Unknown tool: %s. Use --list to see available tools.\n"
	ErrInstallNoPkgMgr        = "No package manager found. Install Chocolatey or Winget first.\n"
	ErrInstallFailed          = "\n  ✗ Installation failed for %s.\n"
	ErrInstallFailedReason    = "  → Reason: %v\n"
	ErrInstallFailedVersion   = "  → Attempted version: %s\n"
	ErrInstallFailedManager   = "  → Package manager: %s\n"
	ErrInstallFailedCmd       = "  → Command: %s\n"
	ErrInstallFailedLog       = "  → Error log: %s\n"
	ErrInstallFailedHint      = "  → Share the log file with an AI or support to diagnose the issue.\n"
	ErrInstallVerifyFailed    = "\n  ✗ Post-install verification failed for %s.\n"
	ErrInstallAdminRequired   = "%s requires administrator privileges to install.\n"
	ErrInstallNetworkRequired = "Network connection required for installation."
	ErrInstallExeNotFound     = "  Error: post-install binary not found at %s (operation: verify, reason: file does not exist)\n"
	ErrMuseDownloadFailed     = "E1085"
	ErrMissingDownloadTool    = "E1086"
	ErrMuseProbeFailed        = "E1087"
	ErrUnsupportedArch        = "E1088"
)

// Install log directory.
const (
	InstallLogDir = ".gitmap/logs"
)

// Apt-specific messages.
const (
	MsgInstallAptUpdate       = "\n  Updating package index (apt-get update)...\n"
	MsgInstallAptUpdateDone   = "  ✓ Package index updated.\n"
	ErrInstallAptUpdateFailed = "  ⚠ apt-get update failed (continuing anyway): %v\n"
)

// NPP error messages — Code Red: all file errors include exact path and reason.
const (
	ErrNppZipNotFound      = "Error: settings zip not found at %s: %v (operation: extract, reason: file does not exist)\n"
	ErrNppSourceDir        = "Error: settings source directory not found at %s: %v (operation: read, reason: directory does not exist)\n"
	ErrNppDirCreate        = "Error: failed to create directory %s: %v (operation: mkdir, reason: path is inaccessible)\n"
	ErrNppExtractEntry     = "Error: failed to open zip entry '%s' for extraction to %s: %v (operation: extract)\n"
	ErrNppFileCreate       = "Error: failed to create file at %s: %v (operation: write, reason: path is inaccessible)\n"
	ErrNppFileCopy         = "Error: failed to copy zip entry '%s' to %s: %v (operation: extract)\n"
	ErrNppFileRead         = "Error: failed to read settings file at %s: %v (operation: read)\n"
	ErrNppFileWrite        = "Error: failed to write settings file to %s: %v (operation: write)\n"
	ErrNppWindowsOnly      = "Error: Notepad++ settings sync is only supported on Windows (current OS: %s)\n"
	ErrNppNoAppData        = "Error: APPDATA environment variable not set (operation: resolve, reason: environment variable not set)\n"
	MsgNppSettingsSynced   = "Settings synced to %s\n"
	MsgNppSettingsFallback = "Settings synced to %s (fallback — zip was missing)\n"
)

// Uninstall messages.
const (
	MsgUninstallRemoving = "Removing %s...\n"
	MsgUninstallSuccess  = "%s uninstalled successfully.\n"
	MsgUninstallDryCmd   = "[dry-run] Would run: %s\n"
	MsgUninstallConfirm  = "Uninstall %s? (y/N): "
	ErrUninstallFailed   = "Uninstall failed for %s: %v\n"
	ErrUninstallNotFound = "%s is not tracked in the database. Use --force to try anyway.\n"
	ErrUninstallDBRemove = "Warning: could not remove %s from database: %v\n"
)

// Scripts install messages.
const (
	MsgScriptsTarget  = "  → Scripts target: %s\n"
	MsgScriptsCloning = "  Cloning gitmap repo for fspath...\n    %s\n"
	MsgScriptsSkip    = "  ⚠ Skipped (not found): %s\n"
	MsgScriptsCopied  = "  ✓ Copied: %s\n"
	MsgScriptsDone    = "\n  ✅ %d scripts installed to %s\n"
	ErrScriptsMkdir   = "  ✗ Could not create target directory %s: %v\n"
	ErrScriptsTemp    = "  ✗ Could not create temp directory: %v\n"
	ErrScriptsClone   = "  ✗ Clone failed: %v\n"
	ErrScriptsCopy    = "  ✗ Failed to copy %s: %v\n"
)
