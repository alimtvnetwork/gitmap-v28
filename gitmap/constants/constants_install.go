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
	HelpInstall   = "  install (in) <tool> Install a developer tool by name"
	HelpUninstall = "  uninstall (un) <tool> Remove a previously installed tool"
)

// Supported tool names — Core.
const (
	ToolVSCode         = "vscode"
	ToolNodeJS         = "node"
	ToolYarn           = "yarn"
	ToolBun            = "bun"
	ToolPnpm           = "pnpm"
	ToolPython         = "python"
	ToolGo             = "go"
	ToolGit            = "git"
	ToolGitLFS         = "git-lfs"
	ToolGHCLI          = "gh"
	ToolGitHubDesktop  = "github-desktop"
	ToolCPP            = "cpp"
	ToolPHP            = "php"
	ToolPowerShell     = "powershell"
	ToolChocolatey     = "chocolatey"
	ToolWinget         = "winget"
	ToolNpp            = "npp"
	ToolNppSettings    = "npp-settings"
	ToolNppInstall     = "install-npp"
	ToolVSCodeSync     = "vscode-settings"
	ToolOBSSync        = "obs-settings"
	ToolWTSync         = "wt-settings"
	ToolScripts        = "scripts"
	ToolDbeaver        = "dbeaver"
	ToolStickyNotes    = "sticky-notes"
	ToolLiteDB         = "litedb"
	ToolVSCodeCtx      = "vscode-ctx"
	ToolPwshCtx        = "pwsh-ctx"
	ToolOBS            = "obs"
	ToolAllDevTools    = "all"
	ToolGitmapOneliner = "gitmap-oneliner"
)

// Supported tool names — Databases.
const (
	ToolMySQL         = "mysql"
	ToolMariaDB       = "mariadb"
	ToolPostgreSQL    = "postgresql"
	ToolSQLite        = "sqlite"
	ToolMongoDB       = "mongodb"
	ToolCouchDB       = "couchdb"
	ToolRedis         = "redis"
	ToolCassandra     = "cassandra"
	ToolNeo4j         = "neo4j"
	ToolElasticsearch = "elasticsearch"
	ToolDuckDB        = "duckdb"
)

// Package manager names.
const (
	PkgMgrChocolatey = "choco"
	PkgMgrWinget     = "winget"
	PkgMgrApt        = "apt"
	PkgMgrBrew       = "brew"
	PkgMgrSnap       = "snap"
	PkgMgrDnf        = "dnf"
	PkgMgrPacman     = "pacman"
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

// Chocolatey package IDs.
const (
	ChocoPkgVSCode        = "vscode"
	ChocoPkgNodeJS        = "nodejs"
	ChocoPkgYarn          = "yarn"
	ChocoPkgBun           = "bun"
	ChocoPkgPnpm          = "pnpm"
	ChocoPkgPython        = "python"
	ChocoPkgGo            = "golang"
	ChocoPkgGit           = "git"
	ChocoPkgGitLFS        = "git-lfs"
	ChocoPkgGHCLI         = "gh"
	ChocoPkgGitHubDesktop = "github-desktop"
	ChocoPkgCPP           = "mingw"
	ChocoPkgPHP           = "php"
	ChocoPkgMySQL         = "mysql"
	ChocoPkgMariaDB       = "mariadb"
	ChocoPkgPostgreSQL    = "postgresql"
	ChocoPkgSQLite        = "sqlite"
	ChocoPkgMongoDB       = "mongodb"
	ChocoPkgCouchDB       = "couchdb"
	ChocoPkgRedis         = "redis-64"
	ChocoPkgNeo4j         = "neo4j-community"
	ChocoPkgElasticsearch = "elasticsearch"
	ChocoPkgDuckDB        = "duckdb"
	ChocoPkgNpp           = "notepadplusplus"
	ChocoPkgDbeaver       = "dbeaver"
	ChocoPkgOBS           = "obs-studio"
	ChocoPkgPowerShell    = "powershell-core"
	ChocoPkgStickyNotes   = "microsoft-windows-terminal" // sticky notes is a Windows Store app
)

// Winget package IDs.
const (
	WingetPkgVSCode        = "Microsoft.VisualStudioCode"
	WingetPkgPowerShell    = "Microsoft.PowerShell"
	WingetPkgDbeaver       = "dbeaver.DBeaverCommunity"
	WingetPkgOBS           = "OBSProject.OBSStudio"
	WingetPkgStickyNotes   = "9NBLGGH4QGHW" // Microsoft Sticky Notes Store ID
	WingetPkgGitHubDesktop = "GitHub.GitHubDesktop"
)

// Apt package IDs.
const (
	AptPkgNodeJS        = "nodejs"
	AptPkgPython        = "python3"
	AptPkgGo            = "golang"
	AptPkgGit           = "git"
	AptPkgGitLFS        = "git-lfs"
	AptPkgCPP           = "g++"
	AptPkgPHP           = "php"
	AptPkgMySQL         = "mysql-server"
	AptPkgMariaDB       = "mariadb-server"
	AptPkgPostgreSQL    = "postgresql"
	AptPkgSQLite        = "sqlite3"
	AptPkgMongoDB       = "mongod"
	AptPkgCouchDB       = "couchdb"
	AptPkgRedis         = "redis-server"
	AptPkgCassandra     = "cassandra"
	AptPkgElasticsearch = "elasticsearch"
)

// Brew package IDs.
const (
	BrewPkgNodeJS        = "node"
	BrewPkgPython        = "python"
	BrewPkgGo            = "go"
	BrewPkgGit           = "git"
	BrewPkgGitLFS        = "git-lfs"
	BrewPkgGHCLI         = "gh"
	BrewPkgCPP           = "gcc"
	BrewPkgPHP           = "php"
	BrewPkgMySQL         = "mysql"
	BrewPkgMariaDB       = "mariadb"
	BrewPkgPostgreSQL    = "postgresql"
	BrewPkgSQLite        = "sqlite"
	BrewPkgMongoDB       = "mongodb-community"
	BrewPkgCouchDB       = "couchdb"
	BrewPkgRedis         = "redis"
	BrewPkgNeo4j         = "neo4j"
	BrewPkgElasticsearch = "elasticsearch"
	BrewPkgDuckDB        = "duckdb"
	BrewPkgDbeaver       = "dbeaver-community"
	BrewPkgOBS           = "obs"
)

// Snap package IDs.
const (
	SnapPkgCouchDB = "couchdb"
	SnapPkgRedis   = "redis"
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
	ErrInstallToolRequired    = "Tool name is required. Use --list to see available tools."
	ErrInstallUnknownTool     = "Unknown tool: %s. Use --list to see available tools.\n"
	ErrInstallNoPkgMgr        = "No package manager found. Install Chocolatey or Winget first."
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
	MsgScriptsCloning = "  Cloning gitmap repo for scripts...\n    %s\n"
	MsgScriptsSkip    = "  ⚠ Skipped (not found): %s\n"
	MsgScriptsCopied  = "  ✓ Copied: %s\n"
	MsgScriptsDone    = "\n  ✅ %d scripts installed to %s\n"
	ErrScriptsMkdir   = "  ✗ Could not create target directory %s: %v\n"
	ErrScriptsTemp    = "  ✗ Could not create temp directory: %v\n"
	ErrScriptsClone   = "  ✗ Clone failed: %v\n"
	ErrScriptsCopy    = "  ✗ Failed to copy %s: %v\n"
)

// Tool categories.
const (
	ToolCategoryCore     = "Core Tools"
	ToolCategoryDatabase = "Databases"
)

// Tool display names for --list output.
var InstallToolDescriptions = map[string]string{
	ToolVSCode:         "Visual Studio Code editor",
	ToolNodeJS:         "Node.js JavaScript runtime",
	ToolYarn:           "Yarn package manager",
	ToolBun:            "Bun JavaScript runtime",
	ToolPnpm:           "pnpm package manager",
	ToolPython:         "Python programming language",
	ToolGo:             "Go programming language",
	ToolGit:            "Git version control",
	ToolGitLFS:         "Git Large File Storage",
	ToolGHCLI:          "GitHub CLI",
	ToolGitHubDesktop:  "GitHub Desktop application",
	ToolCPP:            "C++ compiler (MinGW/g++)",
	ToolPHP:            "PHP programming language",
	ToolPowerShell:     "PowerShell shell",
	ToolChocolatey:     "Chocolatey package manager",
	ToolWinget:         "Winget package manager",
	ToolMySQL:          "MySQL relational database",
	ToolMariaDB:        "MariaDB (MySQL-compatible fork)",
	ToolPostgreSQL:     "PostgreSQL relational database",
	ToolSQLite:         "SQLite embedded database",
	ToolMongoDB:        "MongoDB document database",
	ToolCouchDB:        "CouchDB document database (REST API)",
	ToolRedis:          "Redis in-memory key-value store",
	ToolCassandra:      "Apache Cassandra wide-column NoSQL",
	ToolNeo4j:          "Neo4j graph database",
	ToolElasticsearch:  "Elasticsearch search and analytics",
	ToolDuckDB:         "DuckDB analytical columnar database",
	ToolDbeaver:        "DBeaver database management tool",
	ToolStickyNotes:    "Microsoft Sticky Notes",
	ToolLiteDB:         "LiteDB embedded NoSQL database for .NET",
	ToolOBS:            "OBS Studio screen recorder and streamer",
	ToolVSCodeCtx:      "Add VS Code to Windows right-click context menu",
	ToolPwshCtx:        "Add PowerShell to Windows right-click context menu",
	ToolCtx:            ToolCtxDesc,
	ToolAllDevTools:    "Install all core developer tools at once",
	ToolNpp:            "NPP + Settings -- Notepad++ with settings",
	ToolNppSettings:    "NPP Settings -- Notepad++ settings sync only",
	ToolNppInstall:     "Install NPP -- Notepad++ install only (no settings)",
	ToolVSCodeSync:     "VS Code Settings -- sync VS Code settings and extensions",
	ToolOBSSync:        "OBS Settings -- sync OBS Studio profiles and scenes",
	ToolWTSync:         "WT Settings -- sync Windows Terminal settings.json",
	ToolScripts:        "Clone gitmap scripts to local folder",
	ToolGitmapOneliner: "Print the Windows + macOS install-gitmap one-liners",
}

// InstallToolCategories groups tools by category for display.
var InstallToolCategories = map[string][]string{
	ToolCategoryCore: {
		ToolVSCode, ToolNodeJS, ToolYarn, ToolBun, ToolPnpm,
		ToolPython, ToolGo, ToolGit, ToolGitLFS, ToolGHCLI,
		ToolGitHubDesktop, ToolCPP, ToolPHP, ToolPowerShell,
		ToolChocolatey, ToolWinget, ToolDbeaver, ToolOBS,
		ToolStickyNotes, ToolVSCodeCtx, ToolPwshCtx, ToolCtx,
		ToolNpp, ToolNppSettings, ToolNppInstall,
		ToolVSCodeSync, ToolOBSSync, ToolWTSync,
		ToolScripts, ToolGitmapOneliner, ToolAllDevTools,
	},
	ToolCategoryDatabase: {
		ToolMySQL, ToolMariaDB, ToolPostgreSQL, ToolSQLite,
		ToolMongoDB, ToolCouchDB, ToolRedis, ToolCassandra,
		ToolNeo4j, ToolElasticsearch, ToolDuckDB, ToolLiteDB,
	},
}
