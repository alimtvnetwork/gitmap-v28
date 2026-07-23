// Package cmd implements the CLI commands for gitmap-updater.
package cmd

// Version is the updater version, set via ldflags at build time.
var Version = "0.1.0"

// GitHub repository coordinates.
const (
	RepoOwner = "alimtvnetwork"
	RepoName  = "gitmap-v28"
	RepoSlug  = RepoOwner + "/" + RepoName
)

// API and URL templates.
const (
	GitHubAPILatest   = "https://api.github.com/repos/" + RepoSlug + "/releases/latest"
	ReleaseInstallURL = "https://github.com/" + RepoSlug + "/releases/download/%s/install.ps1"
)

// Binary names.
const (
	GitMapBin   = "gitmap"
	UpdaterCopy = "gitmap-updater-tmp-%d.exe"
)

// PowerShell execution.
const (
	PSBin        = "powershell"
	PSExecPolicy = "-ExecutionPolicy"
	PSBypass     = "Bypass"
	PSNoProfile  = "-NoProfile"
	PSNoLogo     = "-NoLogo"
	PSFile       = "-File"
)

// UI messages.
const (
	MsgChecking       = "\n  ■ Checking for updates...\n"
	MsgCurrentVer     = "  Current version: %s\n"
	MsgLatestVer      = "  Latest version:  %s\n"
	MsgUpToDate       = "  ✓ Already up to date (%s)\n\n"
	MsgUpdateAvail    = "  %s → %s\n"
	MsgDownloading    = "  ■ Downloading installer for %s...\n"
	MsgRunningInstall = "  ■ Running installer...\n"
	MsgDone           = "\n  ✓ Update complete.\n"
	MsgVerifyFail     = "  ✗ Version verification failed: expected %s, got %s\n"
	MsgHandoff        = "  → Handing off to worker...\n"
)

// Error messages.
const (
	ErrFetchRelease = "  ✗ Failed to check latest release: %v\n"
	ErrGetVersion   = "  ✗ Failed to get installed version: %v\n"
	ErrDownload     = "  ✗ Failed to download installer: %v\n"
	ErrRunInstaller = "  ✗ Installer failed: %v\n"
	ErrCreateCopy   = "  ✗ Failed to create handoff copy: %v\n"
	ErrLaunchWorker = "  ✗ Failed to launch worker: %v\n"
)
