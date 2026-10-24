package constants

// Subcommand names & routing aliases for cmdconfig.
const (
	CmdExportConfig   = "export-config"
	CmdImportConfig   = "import-config"
	CmdWhatConfigs    = "what-configs"
	CmdWhatConfig     = "what-config"
	CmdWhatConfigsAlt = "whatconfigs"
	CmdWC             = "wc"
	CmdWhichFormat    = "which-format"
	CmdWhich          = "which"
	CmdFormat         = "format"
	CmdInspect        = "inspect"
	CmdCat            = "cat"
	CmdVersion        = "version"
	CmdEnv            = "env"
	CmdInstaller      = "installer"
	CmdPipeline       = "pipeline"
	CmdSSH            = "ssh"
	CmdStorage        = "storage"
)

// CLI flags, short forms, and prefixes.
const (
	FlagDir            = "dir"
	FlagDirShort       = "d"
	FlagYes            = "yes"
	FlagYesShort       = "y"
	FlagHelp           = "help"
	FlagHelpShort      = "h"
	FlagArgDir         = "--dir"
	FlagArgDirShort    = "-d"
	FlagArgDirPrefix   = "--dir="
	FlagArgYes         = "--yes"
	FlagArgYesShort    = "-y"
	FlagArgHelp        = "--help"
	FlagArgHelpShort   = "-h"
	FlagListExtensions = "--list-extensions"
	FlagPrefixDash     = "-"
)

// Operating system identifiers and environment tokens.
const (
	OSWindows         = "windows"
	OSDarwin          = "darwin"
	OSLinux           = "linux"
	EnvAppData        = "APPDATA"
	PathMacLibrary    = "Library"
	PathMacAppSupport = "Application Support"
	PathLinuxConfig   = ".config"
)

// Tool directory subpaths and identifiers.
const (
	PathVSCodeUser     = "Code/User"
	PathQBittorrent    = "qBittorrent"
	PathUTorrent       = "uTorrent"
	PathUTorrentLower  = "utorrent"
	ToolVSCode         = "vscode"
	ToolQTorrent       = "qtorrent"
	ToolUTorrent       = "utorrent"
	ToolAll            = "all"
	ToolWildcard       = "*"
	ToolCode           = "code"
	ToolVCode          = "vcode"
	ToolQBittorrent    = "qbittorrent"
	ToolQBit           = "qbit"
	ToolUTTorrent      = "uttorrent"
	ToolUTorrentHyphen = "u-torrent"
)

// Canonical file names, basenames, and extensions.
const (
	FileSettingsJSON    = "settings.json"
	FileKeybindingsJSON = "keybindings.json"
	FileQBittorrentIni  = "qBittorrent.ini"
	FileQBittorrentConf = "qBittorrent.conf"
	FileSettingsDat     = "settings.dat"
	FileExtensionsTxt   = "extensions.txt"
	FilePackageJSON     = "package.json"
	FileVersionJSON     = "version.json"
	FileTSConfigJSON    = "tsconfig.json"
	FileCargoTOML       = "cargo.toml"
	FilePyprojectTOML   = "pyproject.toml"
	FileInstallationDB  = "installation.db"
	FileGitmapDB        = "gitmap.db"
	FileDotEnv          = ".env"
	FileDockerfile      = "dockerfile"
	FileMakefile        = "makefile"
	FileEditorConfig    = ".editorconfig"
	ExtJSON             = ".json"
	ExtYAML             = ".yaml"
	ExtYML              = ".yml"
	ExtDB               = ".db"
	ExtSQLite           = ".sqlite"
	ExtSQLite3          = ".sqlite3"
	ExtTOML             = ".toml"
	ExtINI              = ".ini"
	ExtCFG              = ".cfg"
	ExtCONF             = ".conf"
)

// Payload encoding tokens.
const (
	EncodingUTF8   = "utf-8"
	EncodingBase64 = "base64"
)

// Starter templates for default tool export fallbacks.
const (
	DefaultVSCodeSettings    = "{\n  \"editor.tabSize\": 2,\n  \"editor.formatOnSave\": true\n}\n"
	DefaultVSCodeKeybindings = "[]\n"
	DefaultQBittorrentConfig = "[LegalNotice]\nAccepted=true\n\n[Preferences]\nConnection\\PortRangeMin=6881\nDownloads\\SavePath=Downloads\n"
	DefaultUTorrentSettings  = "d4:name8:uTorrente"
)

// Table column headers and display titles.
const (
	ColFile                 = "FILE"
	ColFormat               = "FORMAT"
	ColPurpose              = "PURPOSE / SUBSYSTEM"
	ColRecommendedCmd       = "RECOMMENDED CLI COMMAND"
	HeaderWhatConfigs       = "GITMAP CONFIGURATION FILE CLASSIFIER (WHAT-CONFIGS)"
	HeaderWhichFormat       = "GITMAP JSON FORMAT INSPECTOR (WHICH-FORMAT)"
	FormatJSONTypedEnvelope = "JSON (Typed Envelope)"
	FormatJSON              = "JSON"
	FormatYAML              = "YAML"
	FormatTOML              = "TOML"
	FormatENV               = "ENV"
	FormatSQLiteDB          = "SQLite DB"
	FormatBinaryDB          = "Binary DB"
)

// Inspection categories for what-configs classifier.
const (
	CategoryUnreadable       = "Unreadable File"
	CategoryEnvVars          = "Environment Variables & Secrets"
	CategoryGenericConfig    = "Generic Configuration"
	CategoryNodeManifest     = "Node.js Package Manifest & Dependencies"
	CategoryVersionPin       = "GitMap Canonical Version Pin Site"
	CategoryTSConfig         = "TypeScript Compiler Options"
	CategoryGenericJSON      = "Generic JSON Document"
	CategoryGitHubWorkflow   = "GitHub Actions CI/CD Pipeline Workflow"
	CategoryDockerCompose    = "Docker Compose Container Infrastructure"
	CategoryYAMLManifest     = "YAML Configuration Manifest"
	CategoryInstallationDB   = "GitMap System & Tools Installation Registry"
	CategoryGitmapDB         = "GitMap Core System & SSH Nodes Split Database"
	CategoryPipelineDB       = "GitMap Local CI/CD Pipeline Database"
	CategorySQLiteDB         = "SQLite Relational Database"
	CategoryCargoTOML        = "Rust Cargo Package & Workspace Manifest"
	CategoryPyprojectTOML    = "Python PEP 518/621 Build & Packaging Manifest"
	CategoryTOMLFile         = "TOML Configuration File"
)

// Inspection system impacts for what-configs classifier.
const (
	ImpactEnvVars          = "Injects environment credentials and service runtime variables."
	ImpactGenericConfig    = "Application configuration settings."
	ImpactNodeManifest     = "Defines JavaScript/TypeScript project dependencies and fspath."
	ImpactVersionPin       = "Defines project release semantic version and build metadata."
	ImpactTSConfig         = "Governs TypeScript build rules, target ECMAScript version, and path aliases."
	ImpactGenericJSON      = "Generic JSON payload (inspect with 'gitmap which-format')."
	ImpactGitHubWorkflow   = "Orchestrates remote GitHub Actions runners, linters, tests, and releases."
	ImpactDockerCompose    = "Defines multi-container application services, networks, and storage volumes."
	ImpactYAMLManifest     = "Application declarative configuration."
	ImpactInstallationDB   = "Tracks installed packages, tool versions, and system telemetry."
	ImpactGitmapDB         = "Stores cluster fleet nodes, SSH host definitions, credentials, and settings."
	ImpactPipelineDB       = "Maintains local CI/CD pipeline task executions, outputs, and status logs."
	ImpactSQLiteDB         = "SQLite storage database (schema tables, rows, indexes)."
	ImpactCargoTOML        = "Defines Rust crate dependencies, compilation targets, and features."
	ImpactPyprojectTOML    = "Defines Python build dependencies, package configuration, and linter settings."
	ImpactTOMLFile         = "Structured TOML settings file."
)

// SQLite detection parameters and defaults.
const (
	SQLiteFormat3        = "SQLite format 3"
	SQLiteHeaderMinLen   = 16
	DefaultSchemaVersion = 1
	DefaultTargetDir     = "."
	DefaultExportFolder  = "."
	DashPlaceholder      = "-"
)

// Error codes for configuration operations.
const (
	ErrCodeConfigExport  = "E7201"
	ErrCodeConfigImport  = "E7202"
	ErrCodeInvalidTool   = "E7203"
	ErrCodeFileNotFound  = "E7204"
)
