// Package constants defines all shared constant values for gitmap.
// No magic strings — all literals used for comparison, defaults,
// formats, and file extensions live here.
//
// Split into focused files:
//
//	constants.go          — core defaults, modes, formats, permissions
//	constants_git.go      — git command and argument strings
//	constants_cli.go      — CLI command names, help text, flag descriptions
//	constants_terminal.go — ANSI colors, terminal sections, table headers
//	constants_messages.go — user-facing messages and error strings
//	constants_release.go  — release workflow messages and setup sections
package constants

// Version.
const Version = "6.80.1"

// RepoPath is set at build time via -ldflags.
var RepoPath = ""

// Clone modes.
const (
	ModeHTTPS = "https"
	ModeSSH   = "ssh"
)

// Output formats.
const (
	OutputTerminal = "terminal"
	OutputCSV      = "csv"
	OutputJSON     = "json"
)

// URL prefixes.
const (
	PrefixHTTP      = "http://"
	PrefixHTTPS     = "https://"
	PrefixSSH       = "git@"
	PrefixSSHScheme = "ssh://"
)

// File extensions.
const (
	ExtCSV  = ".csv"
	ExtJSON = ".json"
	ExtTXT  = ".txt"
	ExtGit  = ".git"
)

// Root directory for all repo-local gitmap data.
const GitMapDir = ".gitmap"

// Subdirectory names within .gitmap/.
const (
	ReleaseDirName  = "release"
	OutputDirName   = "output"
	DeployedDirName = "deployed"
)

// Legacy directory names (pre-.gitmap migration).
const (
	LegacyOutputDir   = "gitmap-output"
	LegacyReleaseDir  = ".release"
	LegacyDeployedDir = ".deployed"
)

// Default file names.
const (
	DefaultCSVFile              = "gitmap.csv"
	DefaultJSONFile             = "gitmap.json"
	DefaultTextFile             = "gitmap.txt"
	DefaultVerboseLogDir        = GitMapDir + "/output"
	DefaultStructureFile        = "folder-structure.md"
	DefaultCloneScript          = "clone.ps1"
	DefaultDirectCloneScript    = "direct-clone.ps1"
	DefaultDirectCloneSSHScript = "direct-clone-ssh.ps1"
	DefaultDesktopScript        = "register-desktop.ps1"
	DefaultScanCacheFile        = "last-scan.json"
	DefaultConfigPath           = "./data/config.json"
	DefaultSetupConfigPath      = "./data/git-setup.json"
	DefaultBuildOutput          = "./bin"
	DefaultOutputDir            = ".gitmap/output"
	DefaultOutputFolder         = "output"
	DefaultBranch               = "main"
	DefaultDir                  = "."
	DefaultVersionFile          = "version.json"
)

// DefaultReleaseDir is a var so tests can override it.
var DefaultReleaseDir = GitMapDir + "/" + ReleaseDirName

const (
	DefaultLatestFile = "latest.json"
)

// JSON formatting.
const JSONIndent = "  "

// Date display formatting.
const (
	DateDisplayLayout = "02-Jan-2006 03:04 PM"
	DateUTCSuffix     = " (UTC)"
)

// Sort orders.
const (
	SortByDate = "date"
	SortByName = "name"
)

// Bump levels.
const (
	BumpMajor = "major"
	BumpMinor = "minor"
	BumpPatch = "patch"
)

// Clone and Desktop scripts are now generated from Go templates
// embedded in formatter/templates/. See clone.ps1.tmpl and desktop.ps1.tmpl.

// File and directory permissions.
const DirPermission = 0o755
const FilePermission = 0o644

// Safe-pull defaults.
const (
	SafePullRetryAttempts    = 4
	SafePullRetryDelayMS     = 600
	WindowsPathWarnThreshold = 240
)

// Verbose log file.
const VerboseLogFileFmt = "gitmap-verbose-%s.log"
