package constants

const (
	WindowsOS       = "windows"
	WindowsShell    = "cmd.exe"
	WindowsShellArg = "/C"
	UnixShell       = "/bin/sh"
	UnixShellArg    = "-c"

	WingetInstallArg = "install"
	ChocoInstallArg  = "install"
	BrewInstallArg   = "install"
	AptInstallArg    = "install"

	WingetQuietArg = "--quiet"
	ChocoYesArg    = "-y"
	BrewQuietArg   = "-q"
	AptYesArg      = "-y"

	PkgMgrVersionArg = "--version"

	ExitCodeSuccess = 0
	ExitCodeError   = 1
	ExitCodeUsage   = 2

	ErrNoPackageManager = "no supported package manager found"
	FormatCmdSpace      = "%s %s %s %s"

	ClusterGitPullCmd   = "gitmap pull --all"
	ClusterGitPushCmd   = "gitmap push --all"
	ClusterGitCommitCmd = "gitmap commit --all"
	ClusterGitStatusCmd = "gitmap status --all"
)
